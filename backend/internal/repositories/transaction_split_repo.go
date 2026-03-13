package repositories

import (
	"context"

	"lovelion/internal/models"

	"gorm.io/gorm"
)

type TransactionSplitRepo struct {
	db *gorm.DB
}

func NewTransactionSplitRepo(db *gorm.DB) *TransactionSplitRepo {
	return &TransactionSplitRepo{db: db}
}

func (r *TransactionSplitRepo) WithTx(tx *gorm.DB) *TransactionSplitRepo {
	return &TransactionSplitRepo{db: tx}
}

func (r *TransactionSplitRepo) BatchCreate(ctx context.Context, splits []models.TransactionSplit) error {
	if len(splits) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&splits).Error
}

func (r *TransactionSplitRepo) DeleteByTransaction(ctx context.Context, txnID string) error {
	return r.db.WithContext(ctx).Where("transaction_id = ?", txnID).Delete(&models.TransactionSplit{}).Error
}
