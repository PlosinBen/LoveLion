package repositories

import (
	"context"

	"lovelion/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionExpenseItemRepo struct {
	db *gorm.DB
}

func NewTransactionExpenseItemRepo(db *gorm.DB) *TransactionExpenseItemRepo {
	return &TransactionExpenseItemRepo{db: db}
}

func (r *TransactionExpenseItemRepo) WithTx(tx *gorm.DB) *TransactionExpenseItemRepo {
	return &TransactionExpenseItemRepo{db: tx}
}

func (r *TransactionExpenseItemRepo) BatchCreate(ctx context.Context, items []models.TransactionExpenseItem) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&items).Error
}

func (r *TransactionExpenseItemRepo) DeleteByExpense(ctx context.Context, expenseID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("expense_id = ?", expenseID).Delete(&models.TransactionExpenseItem{}).Error
}
