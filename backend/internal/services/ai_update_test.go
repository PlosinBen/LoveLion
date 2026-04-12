package services

import (
	"context"
	"testing"
	"time"

	"lovelion/internal/models"
	"lovelion/internal/repositories"
	"lovelion/internal/testutil"
	"lovelion/internal/utils/errorx"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// newTestTransactionService wires up a service with real repositories against
// the given test DB. storage is nil — ai-related update/cancel tests never
// upload images.
func newTestTransactionService(db *gorm.DB) *TransactionService {
	return NewTransactionService(
		db,
		repositories.NewTransactionRepo(db),
		repositories.NewTransactionExpenseRepo(db),
		repositories.NewTransactionExpenseItemRepo(db),
		repositories.NewTransactionDebtRepo(db),
		nil,
	)
}

// createExpenseWithAIStatus creates a fully-formed expense transaction with
// the given ai_status. Status empty means ai_status=NULL.
func createExpenseWithAIStatus(t *testing.T, db *gorm.DB, spaceID uuid.UUID, status string) string {
	t.Helper()
	txnID := "txn_" + uuid.NewString()[:8]
	txn := &models.Transaction{
		ID:          txnID,
		SpaceID:     spaceID,
		Type:        "expense",
		Title:       "Lunch",
		Date:        time.Now(),
		Currency:    "TWD",
		TotalAmount: decimal.NewFromInt(100),
	}
	if status != "" {
		s := status
		txn.AIStatus = &s
		if status == aiStatusFailed {
			txn.AIError = "辨識逾時，請稍後再試"
		}
	}
	require.NoError(t, db.Create(txn).Error)

	expenseID := uuid.New()
	require.NoError(t, db.Create(&models.TransactionExpense{
		ID:            expenseID,
		TransactionID: txnID,
		Category:      "Food",
		ExchangeRate:  decimal.NewFromInt(1),
	}).Error)

	require.NoError(t, db.Create(&models.TransactionExpenseItem{
		ID:        uuid.New(),
		ExpenseID: expenseID,
		Name:      "original item",
		UnitPrice: decimal.NewFromInt(100),
		Quantity:  decimal.NewFromInt(1),
		Amount:    decimal.NewFromInt(100),
	}).Error)

	return txnID
}

func baseUpdateInput() UpdateExpenseInput {
	return UpdateExpenseInput{
		Title:    "Updated",
		Currency: "TWD",
		Expense: ExpenseInput{
			Category:     "Food",
			ExchangeRate: decimal.NewFromInt(1),
			Items: []ExpenseItemInput{
				{Name: "Manual item", UnitPrice: decimal.NewFromInt(200), Quantity: decimal.NewFromInt(1)},
			},
		},
	}
}

func TestUpdateExpense_FailedPlusAIExtractFalse_ClearsAIStatus(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)
	space := createTestSpace(t, db, user.ID)
	txnID := createExpenseWithAIStatus(t, db, space.ID, aiStatusFailed)

	svc := newTestTransactionService(db)
	_, err := svc.UpdateExpense(context.Background(), txnID, space.ID, baseUpdateInput())
	require.NoError(t, err)

	txn := loadTxn(t, db, txnID)
	assert.Nil(t, txn.AIStatus, "ai_status should be cleared to NULL")
	assert.Equal(t, "", txn.AIError)
	// Manual items replaced the original: only the new one should remain.
	require.Len(t, txn.Expense.Items, 1)
	assert.Equal(t, "Manual item", txn.Expense.Items[0].Name)
}

func TestUpdateExpense_FailedPlusAIExtractTrue_ResetsToPendingAndClearsItems(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)
	space := createTestSpace(t, db, user.ID)
	txnID := createExpenseWithAIStatus(t, db, space.ID, aiStatusFailed)

	input := baseUpdateInput()
	input.AIExtract = true
	// Caller may still send items; the service must ignore them so the worker
	// can repopulate from the LLM response.
	input.Expense.Items = []ExpenseItemInput{
		{Name: "stale client-side item", UnitPrice: decimal.NewFromInt(99), Quantity: decimal.NewFromInt(1)},
	}

	svc := newTestTransactionService(db)
	_, err := svc.UpdateExpense(context.Background(), txnID, space.ID, input)
	require.NoError(t, err)

	txn := loadTxn(t, db, txnID)
	require.NotNil(t, txn.AIStatus)
	assert.Equal(t, aiStatusPending, *txn.AIStatus)
	assert.Equal(t, "", txn.AIError)
	assert.Empty(t, txn.Expense.Items, "items should be cleared for worker to repopulate")
}

func TestUpdateExpense_PendingRejected(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)
	space := createTestSpace(t, db, user.ID)
	txnID := createExpenseWithAIStatus(t, db, space.ID, aiStatusPending)

	svc := newTestTransactionService(db)
	_, err := svc.UpdateExpense(context.Background(), txnID, space.ID, baseUpdateInput())
	require.Error(t, err)
	assert.True(t, errorx.Is(err, errorx.ErrConflict), "expected ErrConflict, got %v", err)

	// The row should be untouched.
	txn := loadTxn(t, db, txnID)
	require.NotNil(t, txn.AIStatus)
	assert.Equal(t, aiStatusPending, *txn.AIStatus)
	require.Len(t, txn.Expense.Items, 1)
	assert.Equal(t, "original item", txn.Expense.Items[0].Name)
}

func TestUpdateExpense_ProcessingRejected(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)
	space := createTestSpace(t, db, user.ID)
	txnID := createExpenseWithAIStatus(t, db, space.ID, aiStatusProcessing)

	svc := newTestTransactionService(db)
	_, err := svc.UpdateExpense(context.Background(), txnID, space.ID, baseUpdateInput())
	require.Error(t, err)
	assert.True(t, errorx.Is(err, errorx.ErrConflict))
}

func TestUpdateExpense_NullStatusUntouched(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)
	space := createTestSpace(t, db, user.ID)
	txnID := createExpenseWithAIStatus(t, db, space.ID, "")

	svc := newTestTransactionService(db)
	_, err := svc.UpdateExpense(context.Background(), txnID, space.ID, baseUpdateInput())
	require.NoError(t, err)

	txn := loadTxn(t, db, txnID)
	assert.Nil(t, txn.AIStatus)
}

func TestUpdateExpense_CompletedUntouched(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)
	space := createTestSpace(t, db, user.ID)
	txnID := createExpenseWithAIStatus(t, db, space.ID, aiStatusCompleted)

	svc := newTestTransactionService(db)
	_, err := svc.UpdateExpense(context.Background(), txnID, space.ID, baseUpdateInput())
	require.NoError(t, err)

	txn := loadTxn(t, db, txnID)
	require.NotNil(t, txn.AIStatus)
	assert.Equal(t, aiStatusCompleted, *txn.AIStatus)
}

func TestCancelAIExtract_Pending_Success(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)
	space := createTestSpace(t, db, user.ID)
	txnID := createExpenseWithAIStatus(t, db, space.ID, aiStatusPending)

	svc := newTestTransactionService(db)
	require.NoError(t, svc.CancelAIExtract(context.Background(), txnID, space.ID))

	txn := loadTxn(t, db, txnID)
	assert.Nil(t, txn.AIStatus)
	assert.Equal(t, "", txn.AIError)
}

func TestCancelAIExtract_Processing_Success(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)
	space := createTestSpace(t, db, user.ID)
	txnID := createExpenseWithAIStatus(t, db, space.ID, aiStatusProcessing)

	svc := newTestTransactionService(db)
	require.NoError(t, svc.CancelAIExtract(context.Background(), txnID, space.ID))

	txn := loadTxn(t, db, txnID)
	assert.Nil(t, txn.AIStatus)
}

func TestCancelAIExtract_Failed_Conflict(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)
	space := createTestSpace(t, db, user.ID)
	txnID := createExpenseWithAIStatus(t, db, space.ID, aiStatusFailed)

	svc := newTestTransactionService(db)
	err := svc.CancelAIExtract(context.Background(), txnID, space.ID)
	require.Error(t, err)
	assert.True(t, errorx.Is(err, errorx.ErrConflict))

	// Row is untouched.
	txn := loadTxn(t, db, txnID)
	require.NotNil(t, txn.AIStatus)
	assert.Equal(t, aiStatusFailed, *txn.AIStatus)
	assert.NotEmpty(t, txn.AIError)
}

func TestCancelAIExtract_NullStatus_Conflict(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)
	space := createTestSpace(t, db, user.ID)
	txnID := createExpenseWithAIStatus(t, db, space.ID, "")

	svc := newTestTransactionService(db)
	err := svc.CancelAIExtract(context.Background(), txnID, space.ID)
	require.Error(t, err)
	assert.True(t, errorx.Is(err, errorx.ErrConflict))
}

func TestCancelAIExtract_WrongSpace_Conflict(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)
	space := createTestSpace(t, db, user.ID)
	other := createTestSpace(t, db, user.ID)
	txnID := createExpenseWithAIStatus(t, db, space.ID, aiStatusPending)

	svc := newTestTransactionService(db)
	err := svc.CancelAIExtract(context.Background(), txnID, other.ID)
	require.Error(t, err)
	assert.True(t, errorx.Is(err, errorx.ErrConflict))

	// Original row untouched.
	txn := loadTxn(t, db, txnID)
	require.NotNil(t, txn.AIStatus)
	assert.Equal(t, aiStatusPending, *txn.AIStatus)
}
