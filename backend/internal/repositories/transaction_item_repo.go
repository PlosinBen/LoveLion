package repositories

import (
	"context"

	"lovelion/internal/models"

	"gorm.io/gorm"
)

type TransactionItemRepo struct {
	db *gorm.DB
}

func NewTransactionItemRepo(db *gorm.DB) *TransactionItemRepo {
	return &TransactionItemRepo{db: db}
}

func (r *TransactionItemRepo) WithTx(tx *gorm.DB) *TransactionItemRepo {
	return &TransactionItemRepo{db: tx}
}

func (r *TransactionItemRepo) BatchCreate(ctx context.Context, items []models.TransactionItem) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&items).Error
}

func (r *TransactionItemRepo) DeleteByTransaction(ctx context.Context, txnID string) error {
	return r.db.WithContext(ctx).Where("transaction_id = ?", txnID).Delete(&models.TransactionItem{}).Error
}
