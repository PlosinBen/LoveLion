package services

import (
	"context"
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
	db        *gorm.DB
	txnRepo   *repositories.TransactionRepo
	itemRepo  *repositories.TransactionItemRepo
	splitRepo *repositories.TransactionSplitRepo
}

func NewTransactionService(db *gorm.DB, txnRepo *repositories.TransactionRepo, itemRepo *repositories.TransactionItemRepo, splitRepo *repositories.TransactionSplitRepo) *TransactionService {
	return &TransactionService{
		db:        db,
		txnRepo:   txnRepo,
		itemRepo:  itemRepo,
		splitRepo: splitRepo,
	}
}

type TransactionItemInput struct {
	Name      string
	UnitPrice decimal.Decimal
	Quantity  decimal.Decimal
	Discount  decimal.Decimal
}

type TransactionSplitInput struct {
	Name    string
	Amount  decimal.Decimal
	IsPayer bool
}

type CreateTransactionInput struct {
	Payer         string
	Date          *time.Time
	Currency      string
	TotalAmount   decimal.Decimal
	ExchangeRate  decimal.Decimal
	BillingAmount decimal.Decimal
	HandlingFee   decimal.Decimal
	Category      string
	Title         string
	PaymentMethod string
	Note          string
	Items         []TransactionItemInput
	Splits        []TransactionSplitInput
}

type UpdateTransactionInput struct {
	Payer         string
	Date          *time.Time
	Currency      string
	TotalAmount   *decimal.Decimal
	ExchangeRate  *decimal.Decimal
	BillingAmount *decimal.Decimal
	HandlingFee   *decimal.Decimal
	Category      string
	Title         string
	PaymentMethod string
	Note          string
	Items         []TransactionItemInput
	Splits        []TransactionSplitInput
}

func buildItems(txnID string, inputs []TransactionItemInput) ([]models.TransactionItem, decimal.Decimal) {
	totalAmount := decimal.Zero
	var items []models.TransactionItem

	for _, inp := range inputs {
		quantity := inp.Quantity
		if quantity.IsZero() {
			quantity = decimal.NewFromInt(1)
		}

		amount := inp.UnitPrice.Sub(inp.Discount).Mul(quantity)

		items = append(items, models.TransactionItem{
			ID:            uuid.New(),
			TransactionID: txnID,
			Name:          inp.Name,
			UnitPrice:     inp.UnitPrice,
			Quantity:      quantity,
			Discount:      inp.Discount,
			Amount:        amount,
		})
		totalAmount = totalAmount.Add(amount)
	}

	return items, totalAmount
}

func buildSplits(txnID string, inputs []TransactionSplitInput) []models.TransactionSplit {
	var splits []models.TransactionSplit
	for _, inp := range inputs {
		splits = append(splits, models.TransactionSplit{
			ID:            uuid.New(),
			TransactionID: txnID,
			Name:          inp.Name,
			Amount:        inp.Amount,
			IsPayer:       inp.IsPayer,
		})
	}
	return splits
}

func (s *TransactionService) List(ctx context.Context, spaceID uuid.UUID) ([]models.Transaction, error) {
	transactions, err := s.txnRepo.FindBySpace(ctx, spaceID)
	if err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to fetch transactions")
	}
	return transactions, nil
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

func (s *TransactionService) Create(ctx context.Context, spaceID uuid.UUID, input CreateTransactionInput) (*models.Transaction, error) {
	txnID, err := utils.NewShortID(s.db, "transactions", "id")
	if err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to generate ID")
	}

	txn := &models.Transaction{
		ID:            txnID,
		SpaceID:       spaceID,
		Payer:         input.Payer,
		Currency:      input.Currency,
		ExchangeRate:  input.ExchangeRate,
		BillingAmount: input.BillingAmount,
		HandlingFee:   input.HandlingFee,
		Category:      input.Category,
		Title:         input.Title,
		PaymentMethod: input.PaymentMethod,
		Note:          input.Note,
	}

	if input.Date != nil {
		txn.Date = *input.Date
	} else {
		txn.Date = time.Now()
	}

	if txn.Currency == "" {
		txn.Currency = "TWD"
	}

	if txn.ExchangeRate.IsZero() {
		txn.ExchangeRate = decimal.NewFromInt(1)
	}

	if len(input.Items) > 0 {
		items, totalAmount := buildItems(txnID, input.Items)
		txn.Items = items
		txn.TotalAmount = totalAmount
	} else {
		txn.TotalAmount = input.TotalAmount
	}

	txn.Splits = buildSplits(txnID, input.Splits)

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txnRepo := s.txnRepo.WithTx(tx)
		return txnRepo.Create(ctx, txn)
	}); err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to create transaction")
	}

	return txn, nil
}

func (s *TransactionService) Update(ctx context.Context, txnID string, spaceID uuid.UUID, input UpdateTransactionInput) (*models.Transaction, error) {
	// Verify existence
	if _, err := s.txnRepo.FindByID(ctx, txnID, spaceID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.Wrap(errorx.ErrNotFound, "Transaction not found")
		}
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to fetch transaction")
	}

	// Build update params
	params := repositories.TransactionUpdateParams{
		Date:          input.Date,
		TotalAmount:   input.TotalAmount,
		ExchangeRate:  input.ExchangeRate,
		BillingAmount: input.BillingAmount,
		HandlingFee:   input.HandlingFee,
		Note:          &input.Note,
	}
	if input.Payer != "" {
		params.Payer = &input.Payer
	}
	if input.Currency != "" {
		params.Currency = &input.Currency
	}
	if input.Category != "" {
		params.Category = &input.Category
	}
	if input.Title != "" {
		params.Title = &input.Title
	}
	if input.PaymentMethod != "" {
		params.PaymentMethod = &input.PaymentMethod
	}

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txnRepo := s.txnRepo.WithTx(tx)
		itemRepo := s.itemRepo.WithTx(tx)
		splitRepo := s.splitRepo.WithTx(tx)

		// Replace items if provided
		if input.Items != nil {
			if err := itemRepo.DeleteByTransaction(ctx, txnID); err != nil {
				return err
			}
			if len(input.Items) > 0 {
				items, totalAmount := buildItems(txnID, input.Items)
				if err := itemRepo.BatchCreate(ctx, items); err != nil {
					return err
				}
				params.TotalAmount = &totalAmount
			}
		}

		// Replace splits if provided
		if input.Splits != nil {
			if err := splitRepo.DeleteByTransaction(ctx, txnID); err != nil {
				return err
			}
			splits := buildSplits(txnID, input.Splits)
			if err := splitRepo.BatchCreate(ctx, splits); err != nil {
				return err
			}
		}

		return txnRepo.Update(ctx, txnID, params)
	}); err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to update transaction")
	}

	// Reload with associations
	return s.txnRepo.FindByID(ctx, txnID, spaceID)
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
