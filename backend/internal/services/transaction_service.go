package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"strings"
	"time"

	"lovelion/internal/models"
	"lovelion/internal/repositories"
	"lovelion/internal/utils"
	"lovelion/internal/utils/errorx"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ImageStorage is the subset of a blob store that the transaction service
// uses when creating transactions with attached images.
type ImageStorage interface {
	Upload(ctx context.Context, key string, body io.Reader, contentType string) (string, error)
	Delete(ctx context.Context, key string) error
}

type TransactionService struct {
	db          *gorm.DB
	txnRepo     *repositories.TransactionRepo
	expenseRepo *repositories.TransactionExpenseRepo
	itemRepo    *repositories.TransactionExpenseItemRepo
	debtRepo    *repositories.TransactionDebtRepo
	storage     ImageStorage // optional — nil means image-bearing flows are rejected
}

func NewTransactionService(
	db *gorm.DB,
	txnRepo *repositories.TransactionRepo,
	expenseRepo *repositories.TransactionExpenseRepo,
	itemRepo *repositories.TransactionExpenseItemRepo,
	debtRepo *repositories.TransactionDebtRepo,
	storage ImageStorage,
) *TransactionService {
	return &TransactionService{
		db:          db,
		txnRepo:     txnRepo,
		expenseRepo: expenseRepo,
		itemRepo:    itemRepo,
		debtRepo:    debtRepo,
		storage:     storage,
	}
}

// --- Input types ---

type ExpenseItemInput struct {
	Name      string
	UnitPrice decimal.Decimal
	Quantity  decimal.Decimal
	Discount  decimal.Decimal
}

type ExpenseInput struct {
	Category      string
	ExchangeRate  decimal.Decimal
	BillingAmount decimal.Decimal
	HandlingFee   decimal.Decimal
	PaymentMethod string
	Items         []ExpenseItemInput
}

type DebtInput struct {
	PayerName  string
	PayeeName  string
	Amount     decimal.Decimal
	IsSpotPaid bool
}

type ImageUpload struct {
	FileName    string // original upload name, used for extension
	Body        []byte
	ContentType string
	BlurHash    string
}

type CreateExpenseInput struct {
	Date        *time.Time
	Currency    string
	TotalAmount decimal.Decimal
	Title       string
	Note        string
	Expense     ExpenseInput
	Debts       []DebtInput
	Images      []ImageUpload // optional — uploaded to R2 in the same tx
	AIExtract   bool          // when true, ai_status is set to pending for worker pickup
}

type UpdateExpenseInput struct {
	Date        *time.Time
	Currency    string
	TotalAmount *decimal.Decimal
	Title       string
	Note        string
	Expense     ExpenseInput
	Debts       []DebtInput
	// AIExtract toggles the AI re-run flow when the current row is in `failed`.
	// See UpdateExpense for the full transition table.
	AIExtract bool
}

type CreatePaymentInput struct {
	Date        *time.Time
	Title       string
	Note        string
	TotalAmount decimal.Decimal
	PayerName   string
	PayeeName   string
}

type UpdatePaymentInput struct {
	Date        *time.Time
	Title       string
	Note        string
	TotalAmount *decimal.Decimal
	PayerName   string
	PayeeName   string
}

// --- Helpers ---

func buildExpenseItems(expenseID uuid.UUID, inputs []ExpenseItemInput) ([]models.TransactionExpenseItem, decimal.Decimal) {
	totalAmount := decimal.Zero
	var items []models.TransactionExpenseItem

	for _, inp := range inputs {
		quantity := inp.Quantity
		if quantity.IsZero() {
			quantity = decimal.NewFromInt(1)
		}

		amount := inp.UnitPrice.Sub(inp.Discount).Mul(quantity)

		items = append(items, models.TransactionExpenseItem{
			ID:        uuid.New(),
			ExpenseID: expenseID,
			Name:      inp.Name,
			UnitPrice: inp.UnitPrice,
			Quantity:  quantity,
			Discount:  inp.Discount,
			Amount:    amount,
		})
		totalAmount = totalAmount.Add(amount)
	}

	return items, totalAmount
}

func calcSettledAmount(debt DebtInput, totalAmount decimal.Decimal, expense ExpenseInput, currency string) decimal.Decimal {
	if debt.IsSpotPaid {
		return decimal.Zero
	}

	billingAmount := expense.BillingAmount

	// Base currency: settled = amount
	isBaseCurrency := currency == "TWD" || billingAmount.IsZero()
	if isBaseCurrency {
		return debt.Amount
	}

	// Foreign currency with billing: ceiling(amount / totalAmount * billingAmount)
	if totalAmount.IsPositive() && billingAmount.IsPositive() {
		ratio := debt.Amount.Div(totalAmount).Mul(billingAmount)
		// Ceiling: round up to integer
		f, _ := ratio.Float64()
		return decimal.NewFromFloat(math.Ceil(f))
	}

	return decimal.Zero
}

func buildDebts(txnID string, inputs []DebtInput, totalAmount decimal.Decimal, expense *ExpenseInput, currency string) []models.TransactionDebt {
	var debts []models.TransactionDebt
	for _, inp := range inputs {
		settled := inp.Amount // default for base currency / payment
		if expense != nil {
			settled = calcSettledAmount(inp, totalAmount, *expense, currency)
		}

		debts = append(debts, models.TransactionDebt{
			ID:            uuid.New(),
			TransactionID: txnID,
			PayerName:     inp.PayerName,
			PayeeName:     inp.PayeeName,
			Amount:        inp.Amount,
			SettledAmount: settled,
			IsSpotPaid:    inp.IsSpotPaid,
		})
	}
	return debts
}

// --- Read operations (shared) ---

func (s *TransactionService) List(ctx context.Context, spaceID uuid.UUID) ([]models.Transaction, error) {
	transactions, err := s.txnRepo.FindBySpace(ctx, spaceID)
	if err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to fetch transactions")
	}
	return transactions, nil
}

func (s *TransactionService) ListPaginated(ctx context.Context, spaceID uuid.UUID, limit, offset int, filter *repositories.TransactionFilter) ([]models.Transaction, int64, error) {
	transactions, total, err := s.txnRepo.FindBySpacePaginated(ctx, spaceID, limit, offset, filter)
	if err != nil {
		return nil, 0, errorx.Wrap(errorx.ErrInternal, "Failed to fetch transactions")
	}
	return transactions, total, nil
}

func (s *TransactionService) GetByID(ctx context.Context, txnID string, spaceID uuid.UUID) (*models.Transaction, error) {
	txn, err := s.txnRepo.FindByID(ctx, txnID, spaceID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.Wrap(errorx.ErrNotFound, "Transaction not found")
		}
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to fetch transaction")
	}
	return txn, nil
}

func (s *TransactionService) Delete(ctx context.Context, txnID string, spaceID uuid.UUID) error {
	rows, err := s.txnRepo.Delete(ctx, txnID, spaceID)
	if err != nil {
		return errorx.Wrap(errorx.ErrInternal, "Failed to delete transaction")
	}
	if rows == 0 {
		return errorx.Wrap(errorx.ErrNotFound, "Transaction not found")
	}
	return nil
}

// --- Expense operations ---

func (s *TransactionService) CreateExpense(ctx context.Context, spaceID uuid.UUID, input CreateExpenseInput) (*models.Transaction, error) {
	txnID, err := utils.NewShortID(s.db, "transactions", "id")
	if err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to generate ID")
	}

	expenseID := uuid.New()

	// Build items and calculate total; fall back to user-provided total when no items
	items, totalAmount := buildExpenseItems(expenseID, input.Expense.Items)
	if totalAmount.IsZero() && input.TotalAmount.IsPositive() {
		totalAmount = input.TotalAmount
	}

	currency := input.Currency
	if currency == "" {
		currency = "TWD"
	}

	exchangeRate := input.Expense.ExchangeRate
	if exchangeRate.IsZero() {
		exchangeRate = decimal.NewFromInt(1)
	}

	txn := &models.Transaction{
		ID:          txnID,
		SpaceID:     spaceID,
		Type:        "expense",
		Title:       input.Title,
		Currency:    currency,
		TotalAmount: totalAmount,
		Note:        input.Note,
	}

	if input.Date != nil {
		txn.Date = *input.Date
	} else {
		txn.Date = time.Now()
	}

	expense := &models.TransactionExpense{
		ID:            expenseID,
		TransactionID: txnID,
		Category:      input.Expense.Category,
		ExchangeRate:  exchangeRate,
		BillingAmount: input.Expense.BillingAmount,
		HandlingFee:   input.Expense.HandlingFee,
		PaymentMethod: input.Expense.PaymentMethod,
	}

	debts := buildDebts(txnID, input.Debts, totalAmount, &input.Expense, currency)

	// Validate image-bearing input before doing any DB work.
	if len(input.Images) > 0 && s.storage == nil {
		return nil, errorx.Wrap(errorx.ErrInternal, "Image storage not configured")
	}
	if input.AIExtract && len(input.Images) == 0 && strings.TrimSpace(input.Title) == "" {
		return nil, errorx.Wrap(errorx.ErrBadRequest, "AI extraction requires an image or text")
	}

	// Keys of objects written to R2 so we can clean them up if the DB tx rolls back.
	var uploadedKeys []string

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txnRepo := s.txnRepo.WithTx(tx)
		expenseRepo := s.expenseRepo.WithTx(tx)
		itemRepo := s.itemRepo.WithTx(tx)
		debtRepo := s.debtRepo.WithTx(tx)

		if err := txnRepo.Create(ctx, txn); err != nil {
			return err
		}
		if err := expenseRepo.Create(ctx, expense); err != nil {
			return err
		}
		if len(items) > 0 {
			if err := itemRepo.BatchCreate(ctx, items); err != nil {
				return err
			}
		}
		if len(debts) > 0 {
			if err := debtRepo.BatchCreate(ctx, debts); err != nil {
				return err
			}
		}

		// Upload images to R2 and insert image records under the same entity_id.
		for i, img := range input.Images {
			_, key, err := s.uploadImageForTransaction(ctx, tx, txnID, i, img)
			if key != "" {
				// Track the key for rollback even if the DB insert failed.
				uploadedKeys = append(uploadedKeys, key)
			}
			if err != nil {
				return err
			}
		}

		// Set ai_status=pending so the worker will pick the row up.
		if input.AIExtract {
			pending := aiStatusPending
			if err := tx.Model(&models.Transaction{}).
				Where("id = ?", txnID).
				Update("ai_status", pending).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		// Rollback R2 objects using background ctx so cleanup still runs if
		// the request was cancelled.
		cleanupCtx := context.Background()
		for _, key := range uploadedKeys {
			_ = s.storage.Delete(cleanupCtx, key)
		}
		// Preserve validation / not-found errors from the tx closure.
		var appErr *errorx.AppError
		if errors.As(err, &appErr) {
			return nil, appErr
		}
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to create expense")
	}

	return s.txnRepo.FindByID(ctx, txnID, spaceID)
}

// uploadImageForTransaction uploads a single image to R2 and inserts an image
// record bound to the transaction. On any error the caller is responsible for
// rolling back both the DB tx and the R2 object (the key is returned so the
// caller can add it to a cleanup list before returning).
func (s *TransactionService) uploadImageForTransaction(
	ctx context.Context,
	tx *gorm.DB,
	txnID string,
	sortOrder int,
	img ImageUpload,
) (*models.Image, string, error) {
	ext := strings.ToLower(filepath.Ext(img.FileName))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return nil, "", errorx.Wrap(errorx.ErrBadRequest, "Only jpg, jpeg, and png are allowed")
	}

	fileID := uuid.New()
	key := fmt.Sprintf("transaction/%s%s", fileID.String(), ext)

	url, err := s.storage.Upload(ctx, key, bytes.NewReader(img.Body), img.ContentType)
	if err != nil {
		return nil, "", errorx.Wrap(errorx.ErrInternal, "Failed to upload image")
	}

	record := &models.Image{
		ID:         fileID,
		EntityID:   txnID,
		EntityType: "transaction",
		FilePath:   url,
		BlurHash:   img.BlurHash,
		SortOrder:  sortOrder,
	}
	if err := tx.Create(record).Error; err != nil {
		// Return the key so caller can delete the uploaded object.
		return nil, key, errorx.Wrap(errorx.ErrInternal, "Failed to save image record")
	}
	return record, key, nil
}

func (s *TransactionService) UpdateExpense(ctx context.Context, txnID string, spaceID uuid.UUID, input UpdateExpenseInput) (*models.Transaction, error) {
	existing, err := s.txnRepo.FindByID(ctx, txnID, spaceID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.Wrap(errorx.ErrNotFound, "Transaction not found")
		}
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to fetch transaction")
	}

	if existing.Type != "expense" {
		return nil, errorx.Wrap(errorx.ErrBadRequest, "Cannot update non-expense transaction as expense")
	}

	if existing.Expense == nil {
		return nil, errorx.Wrap(errorx.ErrInternal, "Expense data not found")
	}

	// Enforce the ai_status transition rules (see design doc §PUT 自動轉換).
	// The caller never sets ai_status directly — we derive the next value from
	// (currentAIStatus, input.AIExtract).
	currentAIStatus := ""
	if existing.AIStatus != nil {
		currentAIStatus = *existing.AIStatus
	}
	if currentAIStatus == aiStatusPending || currentAIStatus == aiStatusProcessing {
		return nil, errorx.Wrap(errorx.ErrConflict, "Transaction is being processed by AI, cannot update")
	}
	// When re-running AI we want the worker to repopulate items from
	// scratch, so clear any items the caller sent.
	if input.AIExtract {
		input.Expense.Items = nil
	}

	expenseID := existing.Expense.ID

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txnRepo := s.txnRepo.WithTx(tx)
		expenseRepo := s.expenseRepo.WithTx(tx)
		itemRepo := s.itemRepo.WithTx(tx)
		debtRepo := s.debtRepo.WithTx(tx)

		// Update transaction base fields
		params := repositories.TransactionUpdateParams{
			Date: input.Date,
			Note: &input.Note,
		}
		if input.Currency != "" {
			params.Currency = &input.Currency
		}
		if input.Title != "" {
			params.Title = &input.Title
		}

		// Replace items
		if err := itemRepo.DeleteByExpense(ctx, expenseID); err != nil {
			return err
		}
		items, totalAmount := buildExpenseItems(expenseID, input.Expense.Items)
		if len(items) > 0 {
			if err := itemRepo.BatchCreate(ctx, items); err != nil {
				return err
			}
		}
		if totalAmount.IsZero() && input.TotalAmount != nil && input.TotalAmount.IsPositive() {
			totalAmount = *input.TotalAmount
		}
		params.TotalAmount = &totalAmount

		// Update expense fields
		exchangeRate := input.Expense.ExchangeRate
		if exchangeRate.IsZero() {
			exchangeRate = decimal.NewFromInt(1)
		}
		expenseParams := repositories.ExpenseUpdateParams{
			Category:      &input.Expense.Category,
			ExchangeRate:  &exchangeRate,
			BillingAmount: &input.Expense.BillingAmount,
			HandlingFee:   &input.Expense.HandlingFee,
			PaymentMethod: &input.Expense.PaymentMethod,
		}
		if err := expenseRepo.Update(ctx, txnID, expenseParams); err != nil {
			return err
		}

		// Replace debts with recalculated settled_amount
		currency := input.Currency
		if currency == "" {
			currency = existing.Currency
		}
		if err := debtRepo.DeleteByTransaction(ctx, txnID); err != nil {
			return err
		}
		debts := buildDebts(txnID, input.Debts, totalAmount, &input.Expense, currency)
		if len(debts) > 0 {
			if err := debtRepo.BatchCreate(ctx, debts); err != nil {
				return err
			}
		}

		if err := txnRepo.Update(ctx, txnID, params); err != nil {
			return err
		}

		// ai_status transitions:
		//   ai_extract=true  → pending (worker picks up)
		//   failed + ai_extract=false → NULL (user editing manually)
		if input.AIExtract {
			if err := tx.Model(&models.Transaction{}).
				Where("id = ?", txnID).
				Updates(map[string]interface{}{
					"ai_status": aiStatusPending,
					"ai_error":  gorm.Expr("NULL"),
				}).Error; err != nil {
				return err
			}
		} else if currentAIStatus == aiStatusFailed {
			if err := tx.Model(&models.Transaction{}).
				Where("id = ?", txnID).
				Updates(map[string]interface{}{
					"ai_status": gorm.Expr("NULL"),
					"ai_error":  gorm.Expr("NULL"),
				}).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		var appErr *errorx.AppError
		if errors.As(err, &appErr) {
			return nil, appErr
		}
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to update expense")
	}

	return s.txnRepo.FindByID(ctx, txnID, spaceID)
}

// CancelAIExtract aborts an in-flight AI extraction by resetting ai_status to
// NULL. Only `pending` / `processing` rows are eligible — other states return
// Conflict so callers can distinguish "nothing to cancel" from success.
//
// The conditional WHERE means a concurrently-running worker's write-back
// (which is also guarded by ai_status='processing') becomes a no-op, so the
// cancel always wins the race.
func (s *TransactionService) CancelAIExtract(ctx context.Context, txnID string, spaceID uuid.UUID) error {
	result := s.db.WithContext(ctx).
		Model(&models.Transaction{}).
		Where("id = ? AND space_id = ? AND ai_status IN ?", txnID, spaceID, []string{aiStatusPending, aiStatusProcessing}).
		Updates(map[string]interface{}{
			"ai_status": gorm.Expr("NULL"),
			"ai_error":  gorm.Expr("NULL"),
		})
	if result.Error != nil {
		return errorx.Wrap(errorx.ErrInternal, "Failed to cancel AI extraction")
	}
	if result.RowsAffected == 0 {
		return errorx.Wrap(errorx.ErrConflict, "Transaction is not being processed by AI")
	}
	return nil
}

// --- Payment operations ---

func (s *TransactionService) CreatePayment(ctx context.Context, spaceID uuid.UUID, baseCurrency string, input CreatePaymentInput) (*models.Transaction, error) {
	if input.PayerName == "" || input.PayeeName == "" {
		return nil, errorx.Wrap(errorx.ErrBadRequest, "Payer and payee are required")
	}
	if input.PayerName == input.PayeeName {
		return nil, errorx.Wrap(errorx.ErrBadRequest, "Payer and payee must be different")
	}
	if !input.TotalAmount.IsPositive() {
		return nil, errorx.Wrap(errorx.ErrBadRequest, "Amount must be positive")
	}

	txnID, err := utils.NewShortID(s.db, "transactions", "id")
	if err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to generate ID")
	}

	if baseCurrency == "" {
		baseCurrency = "TWD"
	}

	txn := &models.Transaction{
		ID:          txnID,
		SpaceID:     spaceID,
		Type:        "payment",
		Title:       input.Title,
		Currency:    baseCurrency,
		TotalAmount: input.TotalAmount,
		Note:        input.Note,
	}

	if input.Date != nil {
		txn.Date = *input.Date
	} else {
		txn.Date = time.Now()
	}

	debt := models.TransactionDebt{
		ID:            uuid.New(),
		TransactionID: txnID,
		PayerName:     input.PayerName,
		PayeeName:     input.PayeeName,
		Amount:        input.TotalAmount,
		SettledAmount: input.TotalAmount,
		IsSpotPaid:    false,
	}

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txnRepo := s.txnRepo.WithTx(tx)
		debtRepo := s.debtRepo.WithTx(tx)

		if err := txnRepo.Create(ctx, txn); err != nil {
			return err
		}
		return debtRepo.BatchCreate(ctx, []models.TransactionDebt{debt})
	}); err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to create payment")
	}

	return s.txnRepo.FindByID(ctx, txnID, spaceID)
}

func (s *TransactionService) UpdatePayment(ctx context.Context, txnID string, spaceID uuid.UUID, input UpdatePaymentInput) (*models.Transaction, error) {
	existing, err := s.txnRepo.FindByID(ctx, txnID, spaceID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.Wrap(errorx.ErrNotFound, "Transaction not found")
		}
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to fetch transaction")
	}

	if existing.Type != "payment" {
		return nil, errorx.Wrap(errorx.ErrBadRequest, "Cannot update non-payment transaction as payment")
	}

	if input.PayerName == "" || input.PayeeName == "" {
		return nil, errorx.Wrap(errorx.ErrBadRequest, "Payer and payee are required")
	}
	if input.PayerName == input.PayeeName {
		return nil, errorx.Wrap(errorx.ErrBadRequest, "Payer and payee must be different")
	}

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txnRepo := s.txnRepo.WithTx(tx)
		debtRepo := s.debtRepo.WithTx(tx)

		params := repositories.TransactionUpdateParams{
			Date:        input.Date,
			TotalAmount: input.TotalAmount,
			Note:        &input.Note,
		}
		if input.Title != "" {
			params.Title = &input.Title
		}

		// Replace debt
		if err := debtRepo.DeleteByTransaction(ctx, txnID); err != nil {
			return err
		}

		amount := existing.TotalAmount
		if input.TotalAmount != nil {
			amount = *input.TotalAmount
		}

		debt := models.TransactionDebt{
			ID:            uuid.New(),
			TransactionID: txnID,
			PayerName:     input.PayerName,
			PayeeName:     input.PayeeName,
			Amount:        amount,
			SettledAmount: amount,
			IsSpotPaid:    false,
		}
		if err := debtRepo.BatchCreate(ctx, []models.TransactionDebt{debt}); err != nil {
			return err
		}

		return txnRepo.Update(ctx, txnID, params)
	}); err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to update payment")
	}

	return s.txnRepo.FindByID(ctx, txnID, spaceID)
}
