package services

import (
	"time"

	"lovelion/internal/models"
	"lovelion/internal/utils"
	"lovelion/internal/utils/errorx"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TransactionService struct {
	db *gorm.DB
}

func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{db: db}
}

type TransactionItemInput struct {
	Name      string
	UnitPrice decimal.Decimal
	Quantity  decimal.Decimal
	Discount  decimal.Decimal
}

type TransactionSplitInput struct {
	MemberID *uuid.UUID
	Name     string
	Amount   decimal.Decimal
	IsPayer  bool
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
			MemberID:      inp.MemberID,
			Name:          inp.Name,
			Amount:        inp.Amount,
			IsPayer:       inp.IsPayer,
		})
	}
	return splits
}

func (s *TransactionService) preload(query *gorm.DB) *gorm.DB {
	return query.Preload("Items").Preload("Splits").Preload("Images", "entity_type = ?", "transaction")
}

func (s *TransactionService) List(spaceID uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := s.preload(s.db).Where("ledger_id = ?", spaceID).
		Order("date DESC").
		Find(&transactions).Error

	if err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to fetch transactions")
	}

	return transactions, nil
}

func (s *TransactionService) GetByID(txnID string, spaceID uuid.UUID) (*models.Transaction, error) {
	var txn models.Transaction
	err := s.preload(s.db).Where("id = ? AND ledger_id = ?", txnID, spaceID).First(&txn).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.Wrap(errorx.ErrNotFound, "Transaction not found")
		}
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to fetch transaction")
	}

	return &txn, nil
}

func (s *TransactionService) Create(spaceID uuid.UUID, input CreateTransactionInput) (*models.Transaction, error) {
	txnID, err := utils.NewShortID(s.db, "transactions", "id")
	if err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to generate ID")
	}

	txn := &models.Transaction{
		ID:            txnID,
		LedgerID:      spaceID,
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

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(txn).Error
	}); err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to create transaction")
	}

	return txn, nil
}

func (s *TransactionService) Update(txnID string, spaceID uuid.UUID, input UpdateTransactionInput) (*models.Transaction, error) {
	var txn models.Transaction
	if err := s.db.Where("id = ? AND ledger_id = ?", txnID, spaceID).First(&txn).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.Wrap(errorx.ErrNotFound, "Transaction not found")
		}
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to fetch transaction")
	}

	// Apply field updates
	if input.Payer != "" {
		txn.Payer = input.Payer
	}
	if input.Date != nil {
		txn.Date = *input.Date
	}
	if input.Currency != "" {
		txn.Currency = input.Currency
	}
	if input.TotalAmount != nil {
		txn.TotalAmount = *input.TotalAmount
	}
	if input.ExchangeRate != nil {
		txn.ExchangeRate = *input.ExchangeRate
	}
	if input.BillingAmount != nil {
		txn.BillingAmount = *input.BillingAmount
	}
	if input.HandlingFee != nil {
		txn.HandlingFee = *input.HandlingFee
	}
	if input.Category != "" {
		txn.Category = input.Category
	}
	if input.Title != "" {
		txn.Title = input.Title
	}
	if input.PaymentMethod != "" {
		txn.PaymentMethod = input.PaymentMethod
	}
	txn.Note = input.Note

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if input.Items != nil {
			if err := tx.Where("transaction_id = ?", txnID).Delete(&models.TransactionItem{}).Error; err != nil {
				return err
			}

			if len(input.Items) > 0 {
				items, totalAmount := buildItems(txnID, input.Items)
				txn.Items = items
				txn.TotalAmount = totalAmount
			} else {
				txn.Items = nil
			}
		}

		if input.Splits != nil {
			if err := tx.Where("transaction_id = ?", txnID).Delete(&models.TransactionSplit{}).Error; err != nil {
				return err
			}

			txn.Splits = buildSplits(txnID, input.Splits)
		}

		return tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(&txn).Error
	}); err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to update transaction")
	}

	// Reload with associations
	s.preload(s.db).First(&txn, "id = ?", txnID)

	return &txn, nil
}

func (s *TransactionService) Delete(txnID string, spaceID uuid.UUID) error {
	result := s.db.Where("id = ? AND ledger_id = ?", txnID, spaceID).Delete(&models.Transaction{})
	if result.Error != nil {
		return errorx.Wrap(errorx.ErrInternal, "Failed to delete transaction")
	}
	if result.RowsAffected == 0 {
		return errorx.Wrap(errorx.ErrNotFound, "Transaction not found")
	}
	return nil
}
