package services

import (
	"context"
	"math"
	"time"

	"lovelion/internal/models"
	"lovelion/internal/repositories"
	"lovelion/internal/utils"
	"lovelion/internal/utils/errorx"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TransactionService struct {
	db          *gorm.DB
	txnRepo     *repositories.TransactionRepo
	expenseRepo *repositories.TransactionExpenseRepo
	itemRepo    *repositories.TransactionExpenseItemRepo
	debtRepo    *repositories.TransactionDebtRepo
}

func NewTransactionService(
	db *gorm.DB,
	txnRepo *repositories.TransactionRepo,
	expenseRepo *repositories.TransactionExpenseRepo,
	itemRepo *repositories.TransactionExpenseItemRepo,
	debtRepo *repositories.TransactionDebtRepo,
) *TransactionService {
	return &TransactionService{
		db:          db,
		txnRepo:     txnRepo,
		expenseRepo: expenseRepo,
		itemRepo:    itemRepo,
		debtRepo:    debtRepo,
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

type CreateExpenseInput struct {
	Date     *time.Time
	Currency string
	Title    string
	Note     string
	Expense  ExpenseInput
	Debts    []DebtInput
}

type UpdateExpenseInput struct {
	Date        *time.Time
	Currency    string
	TotalAmount *decimal.Decimal
	Title       string
	Note        string
	Expense     ExpenseInput
	Debts       []DebtInput
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

func (s *TransactionService) ListPaginated(ctx context.Context, spaceID uuid.UUID, limit, offset int) ([]models.Transaction, int64, error) {
	transactions, total, err := s.txnRepo.FindBySpacePaginated(ctx, spaceID, limit, offset)
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

	// Build items and calculate total
	items, totalAmount := buildExpenseItems(expenseID, input.Expense.Items)

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
		return nil
	}); err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to create expense")
	}

	return s.txnRepo.FindByID(ctx, txnID, spaceID)
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

		return txnRepo.Update(ctx, txnID, params)
	}); err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to update expense")
	}

	return s.txnRepo.FindByID(ctx, txnID, spaceID)
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
