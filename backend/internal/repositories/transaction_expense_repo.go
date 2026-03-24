package repositories

import (
	"context"

	"lovelion/internal/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TransactionExpenseRepo struct {
	db *gorm.DB
}

func NewTransactionExpenseRepo(db *gorm.DB) *TransactionExpenseRepo {
	return &TransactionExpenseRepo{db: db}
}

func (r *TransactionExpenseRepo) WithTx(tx *gorm.DB) *TransactionExpenseRepo {
	return &TransactionExpenseRepo{db: tx}
}

type ExpenseUpdateParams struct {
	Category      *string
	ExchangeRate  *decimal.Decimal
	BillingAmount *decimal.Decimal
	HandlingFee   *decimal.Decimal
	PaymentMethod *string
}

func (r *TransactionExpenseRepo) Create(ctx context.Context, expense *models.TransactionExpense) error {
	return r.db.WithContext(ctx).Create(expense).Error
}

func (r *TransactionExpenseRepo) Update(ctx context.Context, transactionID string, params ExpenseUpdateParams) error {
	updates := map[string]interface{}{}

	if params.Category != nil {
		updates["category"] = *params.Category
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
	if params.PaymentMethod != nil {
		updates["payment_method"] = *params.PaymentMethod
	}

	if len(updates) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Model(&models.TransactionExpense{}).
		Where("transaction_id = ?", transactionID).Updates(updates).Error
}

func (r *TransactionExpenseRepo) DeleteByTransaction(ctx context.Context, txnID string) error {
	return r.db.WithContext(ctx).Where("transaction_id = ?", txnID).Delete(&models.TransactionExpense{}).Error
}
