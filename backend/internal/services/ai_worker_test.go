package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"lovelion/internal/models"
	"lovelion/internal/testutil"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// --- Fakes ---

type fakeExtractor struct {
	mu    sync.Mutex
	calls int
	// If set, these functions drive behaviour on each call in order.
	results []*ReceiptData
	errs    []error
	// sleep is applied before returning; lets tests drive cancellation.
	sleep time.Duration
}

func (f *fakeExtractor) Extract(ctx context.Context, image []byte, mimeType string, _ ExtractHints) (*ReceiptData, error) {
	f.mu.Lock()
	idx := f.calls
	f.calls++
	f.mu.Unlock()

	if f.sleep > 0 {
		select {
		case <-time.After(f.sleep):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	if idx < len(f.errs) && f.errs[idx] != nil {
		return nil, f.errs[idx]
	}
	if idx < len(f.results) {
		return f.results[idx], nil
	}
	return &ReceiptData{}, nil
}

type fakeStorage struct {
	data        []byte
	contentType string
	err         error
	downloaded  []string
}

func (f *fakeStorage) DownloadByURL(ctx context.Context, fullURL string) ([]byte, string, error) {
	f.downloaded = append(f.downloaded, fullURL)
	if f.err != nil {
		return nil, "", f.err
	}
	return f.data, f.contentType, nil
}

// --- Helpers ---

func createTestSpace(t *testing.T, db *gorm.DB, ownerID uuid.UUID) *models.Space {
	t.Helper()
	space := &models.Space{
		ID:     uuid.New(),
		Name:   "Test Space",
		UserID: ownerID,
	}
	require.NoError(t, db.Create(space).Error)
	return space
}

// createPendingExpense creates a transaction+expense row with ai_status=pending
// and an attached image record. Returns the transaction ID.
func createPendingExpense(t *testing.T, db *gorm.DB, spaceID uuid.UUID, imageURL string) string {
	t.Helper()
	txnID := "txn_" + uuid.NewString()[:8]
	pending := aiStatusPending
	txn := &models.Transaction{
		ID:          txnID,
		SpaceID:     spaceID,
		Type:        "expense",
		Title:       "AI receipt",
		Date:        time.Now(),
		Currency:    "TWD",
		TotalAmount: decimal.Zero,
		AIStatus:    &pending,
	}
	require.NoError(t, db.Create(txn).Error)

	expense := &models.TransactionExpense{
		ID:            uuid.New(),
		TransactionID: txnID,
		Category:      "其他",
		ExchangeRate:  decimal.NewFromInt(1),
	}
	require.NoError(t, db.Create(expense).Error)

	img := &models.Image{
		ID:         uuid.New(),
		EntityID:   txnID,
		EntityType: "transaction",
		FilePath:   imageURL,
		SortOrder:  0,
	}
	require.NoError(t, db.Create(img).Error)

	return txnID
}

func loadTxn(t *testing.T, db *gorm.DB, txnID string) *models.Transaction {
	t.Helper()
	var txn models.Transaction
	require.NoError(t, db.Preload("Expense.Items").First(&txn, "id = ?", txnID).Error)
	return &txn
}

func newTestWorker(db *gorm.DB, ext ReceiptExtractor, store ImageDownloader) *AIWorker {
	return NewAIWorker(db, ext, store, AIWorkerConfig{
		PollInterval: 50 * time.Millisecond,
		BatchSize:    5,
	})
}

// --- Tests ---

func TestAIWorker_ProcessOne_Success(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)
	space := createTestSpace(t, db, user.ID)

	date := time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC)
	ext := &fakeExtractor{
		results: []*ReceiptData{{
			Date: &date,
			Items: []ReceiptItem{
				{Name: "Coffee", UnitPrice: decimal.NewFromInt(120), Quantity: decimal.NewFromInt(1)},
				{Name: "Bagel", UnitPrice: decimal.NewFromInt(80), Quantity: decimal.NewFromInt(2)},
			},
		}},
	}
	store := &fakeStorage{data: []byte("img"), contentType: "image/jpeg"}

	txnID := createPendingExpense(t, db, space.ID, "https://cdn/test/transaction/a.jpg")

	worker := newTestWorker(db, ext, store)
	worker.processOne(context.Background(), txnID, "AI receipt")

	txn := loadTxn(t, db, txnID)
	require.NotNil(t, txn.AIStatus)
	assert.Equal(t, aiStatusCompleted, *txn.AIStatus)
	assert.Equal(t, "", txn.AIError)
	assert.Equal(t, "2025-03-15", txn.Date.UTC().Format("2006-01-02"))

	require.NotNil(t, txn.Expense)
	require.Len(t, txn.Expense.Items, 2)
	assert.Equal(t, "Coffee", txn.Expense.Items[0].Name)
	// total = 120*1 + 80*2 = 280
	assert.True(t, decimal.NewFromInt(280).Equal(txn.TotalAmount), "total_amount=%s", txn.TotalAmount)

	assert.Equal(t, 1, ext.calls)
	assert.Equal(t, []string{"https://cdn/test/transaction/a.jpg"}, store.downloaded)
}

func TestAIWorker_ProcessOne_LLMError_MarksFailed(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)
	space := createTestSpace(t, db, user.ID)

	ext := &fakeExtractor{errs: []error{errors.New("gemini http 500: boom")}}
	store := &fakeStorage{data: []byte("img"), contentType: "image/jpeg"}

	txnID := createPendingExpense(t, db, space.ID, "https://cdn/x.jpg")

	worker := newTestWorker(db, ext, store)
	worker.processOne(context.Background(), txnID, "AI receipt")

	txn := loadTxn(t, db, txnID)
	require.NotNil(t, txn.AIStatus)
	assert.Equal(t, aiStatusFailed, *txn.AIStatus)
	assert.Contains(t, txn.AIError, "辨識服務錯誤")
}

func TestAIWorker_ProcessOne_NoImage_MarksFailed(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)
	space := createTestSpace(t, db, user.ID)

	// Create a pending transaction without any image.
	txnID := "txn_" + uuid.NewString()[:8]
	pending := aiStatusPending
	txn := &models.Transaction{
		ID: txnID, SpaceID: space.ID, Type: "expense", Title: "no img",
		Date: time.Now(), Currency: "TWD", TotalAmount: decimal.Zero, AIStatus: &pending,
	}
	require.NoError(t, db.Create(txn).Error)
	require.NoError(t, db.Create(&models.TransactionExpense{
		ID: uuid.New(), TransactionID: txnID, ExchangeRate: decimal.NewFromInt(1),
	}).Error)

	ext := &fakeExtractor{}
	store := &fakeStorage{}
	worker := newTestWorker(db, ext, store)
	// No text extractor configured → no-image row fails with a service-off message.
	worker.processOne(context.Background(), txnID, txn.Title)

	got := loadTxn(t, db, txnID)
	require.NotNil(t, got.AIStatus)
	assert.Equal(t, aiStatusFailed, *got.AIStatus)
	assert.Contains(t, got.AIError, "辨識服務未啟用")
	assert.Equal(t, 0, ext.calls, "image extractor should not be called when no image")
}

// fakeTextExtractor records ExtractText calls and returns a canned ReceiptData.
type fakeTextExtractor struct {
	mu      sync.Mutex
	calls   int
	lastIn  string
	result  *ReceiptData
	err     error
}

func (f *fakeTextExtractor) ExtractText(ctx context.Context, text string, _ ExtractHints) (*ReceiptData, error) {
	f.mu.Lock()
	f.calls++
	f.lastIn = text
	f.mu.Unlock()
	if f.err != nil {
		return nil, f.err
	}
	return f.result, nil
}

func TestAIWorker_ProcessOne_TextExtract_Success(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)
	space := createTestSpace(t, db, user.ID)

	// Pending row with no image, title holds the raw user input.
	txnID := "txn_" + uuid.NewString()[:8]
	pending := aiStatusPending
	txn := &models.Transaction{
		ID: txnID, SpaceID: space.ID, Type: "expense", Title: "停車費 100",
		Date: time.Now(), Currency: "TWD", TotalAmount: decimal.Zero, AIStatus: &pending,
	}
	require.NoError(t, db.Create(txn).Error)
	require.NoError(t, db.Create(&models.TransactionExpense{
		ID: uuid.New(), TransactionID: txnID, ExchangeRate: decimal.NewFromInt(1),
	}).Error)

	text := &fakeTextExtractor{result: &ReceiptData{
		Items: []ReceiptItem{
			{Name: "停車費", UnitPrice: decimal.NewFromInt(100), Quantity: decimal.NewFromInt(1)},
		},
	}}
	worker := newTestWorker(db, &fakeExtractor{}, &fakeStorage{}).WithTextExtractor(text)
	worker.processOne(context.Background(), txnID, txn.Title)

	got := loadTxn(t, db, txnID)
	require.NotNil(t, got.AIStatus)
	assert.Equal(t, aiStatusCompleted, *got.AIStatus)
	assert.Equal(t, "停車費 100", text.lastIn, "text extractor receives the raw title")
	// Title is overwritten with the cleaned item name.
	assert.Equal(t, "停車費", got.Title)
	require.Len(t, got.Expense.Items, 1)
	assert.True(t, decimal.NewFromInt(100).Equal(got.TotalAmount))
}

func TestAIWorker_ProcessOne_ClaimRaceLosesSilently(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)
	space := createTestSpace(t, db, user.ID)

	txnID := createPendingExpense(t, db, space.ID, "https://cdn/x.jpg")

	// Simulate another actor flipping status before worker claims.
	require.NoError(t, db.Model(&models.Transaction{}).
		Where("id = ?", txnID).
		Update("ai_status", gorm.Expr("NULL")).Error)

	ext := &fakeExtractor{}
	store := &fakeStorage{data: []byte("img"), contentType: "image/jpeg"}
	worker := newTestWorker(db, ext, store)
	worker.processOne(context.Background(), txnID, "AI receipt")

	var txn models.Transaction
	require.NoError(t, db.First(&txn, "id = ?", txnID).Error)
	assert.Nil(t, txn.AIStatus)
	assert.Equal(t, 0, ext.calls)
}

func TestAIWorker_ProcessOne_CancelMidCall_NoWriteBack(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)
	space := createTestSpace(t, db, user.ID)

	// Simulate cancel by flipping the row to NULL AFTER claim but BEFORE write-back:
	// the fake extractor sleeps, and in the meantime we manually clear ai_status.
	txnID := createPendingExpense(t, db, space.ID, "https://cdn/x.jpg")

	date := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	ext := &fakeExtractor{
		results: []*ReceiptData{{
			Date:  &date,
			Items: []ReceiptItem{{Name: "X", UnitPrice: decimal.NewFromInt(10), Quantity: decimal.NewFromInt(1)}},
		}},
		sleep: 100 * time.Millisecond,
	}
	store := &fakeStorage{data: []byte("img"), contentType: "image/jpeg"}
	worker := newTestWorker(db, ext, store)

	done := make(chan struct{})
	go func() {
		worker.processOne(context.Background(), txnID, "AI receipt")
		close(done)
	}()

	// While worker is in the fake LLM sleep, clear ai_status (simulating cancel).
	time.Sleep(30 * time.Millisecond)
	require.NoError(t, db.Model(&models.Transaction{}).
		Where("id = ?", txnID).
		Update("ai_status", gorm.Expr("NULL")).Error)

	<-done

	txn := loadTxn(t, db, txnID)
	assert.Nil(t, txn.AIStatus, "cancel should win, status stays NULL")
	// items should not have been created by the worker
	assert.Len(t, txn.Expense.Items, 0)
}

func TestAIWorker_RecoverStuck(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)
	space := createTestSpace(t, db, user.ID)

	// Create a row stuck in processing.
	txnID := "txn_" + uuid.NewString()[:8]
	processing := aiStatusProcessing
	require.NoError(t, db.Create(&models.Transaction{
		ID: txnID, SpaceID: space.ID, Type: "expense", Title: "stuck",
		Date: time.Now(), Currency: "TWD", TotalAmount: decimal.Zero, AIStatus: &processing,
	}).Error)

	worker := newTestWorker(db, &fakeExtractor{}, &fakeStorage{})
	require.NoError(t, worker.recoverStuck(context.Background()))

	var got models.Transaction
	require.NoError(t, db.First(&got, "id = ?", txnID).Error)
	require.NotNil(t, got.AIStatus)
	assert.Equal(t, aiStatusPending, *got.AIStatus)
}

func TestAIWorker_Run_ProcessesPendingThenShutsDown(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)
	space := createTestSpace(t, db, user.ID)

	ext := &fakeExtractor{
		results: []*ReceiptData{{
			Items: []ReceiptItem{{Name: "One", UnitPrice: decimal.NewFromInt(50), Quantity: decimal.NewFromInt(1)}},
		}},
	}
	store := &fakeStorage{data: []byte("img"), contentType: "image/jpeg"}
	worker := newTestWorker(db, ext, store)

	txnID := createPendingExpense(t, db, space.ID, "https://cdn/x.jpg")

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		worker.Run(ctx)
		close(done)
	}()

	// Wait for processing to complete.
	require.Eventually(t, func() bool {
		var txn models.Transaction
		if err := db.First(&txn, "id = ?", txnID).Error; err != nil {
			return false
		}
		return txn.AIStatus != nil && *txn.AIStatus == aiStatusCompleted
	}, 2*time.Second, 25*time.Millisecond, "expected completion")

	cancel()
	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("worker did not shut down on context cancel")
	}
}

func TestAIWorker_FriendlyErrorMessages(t *testing.T) {
	cases := []struct {
		in       error
		contains string
	}{
		{errors.New("context deadline exceeded"), "逾時"},
		{errors.New("gemini http 429: rate limited"), "忙碌"},
		{errors.New("gemini http 500: boom"), "服務錯誤"},
		{errors.New("parse receipt json: unexpected"), "格式錯誤"},
		{errors.New("something unknown"), "辨識失敗"},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%v", c.in), func(t *testing.T) {
			assert.Contains(t, friendlyExtractError(c.in), c.contains)
		})
	}
}
