package repositories

import (
	"context"
	"time"

	"lovelion/internal/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TransactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) WithTx(tx *gorm.DB) *TransactionRepo {
	return &TransactionRepo{db: tx}
}

type TransactionUpdateParams struct {
	Payer         *string
	Date          *time.Time
	Currency      *string
	TotalAmount   *decimal.Decimal
	ExchangeRate  *decimal.Decimal
	BillingAmount *decimal.Decimal
	HandlingFee   *decimal.Decimal
	Category      *string
	Title         *string
	PaymentMethod *string
	Note          *string
}

func (r *TransactionRepo) Create(ctx context.Context, txn *models.Transaction) error {
	return r.db.WithContext(ctx).Create(txn).Error
}

func (r *TransactionRepo) FindBySpace(ctx context.Context, spaceID uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.WithContext(ctx).
		Where("ledger_id = ?", spaceID).
		Preload("Items").
		Preload("Splits").
		Preload("Images", "entity_type = ?", "transaction").
		Order("date DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepo) FindByID(ctx context.Context, id string, spaceID uuid.UUID) (*models.Transaction, error) {
	var txn models.Transaction
	err := r.db.WithContext(ctx).
		Where("id = ? AND ledger_id = ?", id, spaceID).
		Preload("Items").
		Preload("Splits").
		Preload("Images", "entity_type = ?", "transaction").
		First(&txn).Error
	if err != nil {
		return nil, err
	}
	return &txn, nil
}

func (r *TransactionRepo) Update(ctx context.Context, id string, params TransactionUpdateParams) error {
	updates := map[string]interface{}{}

	if params.Payer != nil {
		updates["payer"] = *params.Payer
	}
	if params.Date != nil {
		updates["date"] = *params.Date
	}
	if params.Currency != nil {
		updates["currency"] = *params.Currency
	}
	if params.TotalAmount != nil {
		updates["total_amount"] = *params.TotalAmount
	}
	if params.ExchangeRate != nil {
		updates["exchange_rate"] = *params.ExchangeRate
	}
	if params.BillingAmount != nil {
		updates["billing_amount"] = *params.BillingAmount
	}
	if params.HandlingFee != nil {
		updates["handling_fee"] = *params.HandlingFee
	}
	if params.Category != nil {
		updates["category"] = *params.Category
	}
	if params.Title != nil {
		updates["title"] = *params.Title
	}
	if params.PaymentMethod != nil {
		updates["payment_method"] = *params.PaymentMethod
	}
	if params.Note != nil {
		updates["note"] = *params.Note
	}

	if len(updates) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Model(&models.Transaction{}).Where("id = ?", id).Updates(updates).Error
}

func (r *TransactionRepo) Delete(ctx context.Context, id string, spaceID uuid.UUID) (int64, error) {
	result := r.db.WithContext(ctx).Where("id = ? AND ledger_id = ?", id, spaceID).Delete(&models.Transaction{})
	return result.RowsAffected, result.Error
}
