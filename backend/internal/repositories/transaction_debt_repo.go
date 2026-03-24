package repositories

import (
	"context"

	"lovelion/internal/models"

	"gorm.io/gorm"
)

type TransactionDebtRepo struct {
	db *gorm.DB
}

func NewTransactionDebtRepo(db *gorm.DB) *TransactionDebtRepo {
	return &TransactionDebtRepo{db: db}
}

func (r *TransactionDebtRepo) WithTx(tx *gorm.DB) *TransactionDebtRepo {
	return &TransactionDebtRepo{db: tx}
}

func (r *TransactionDebtRepo) BatchCreate(ctx context.Context, debts []models.TransactionDebt) error {
	if len(debts) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&debts).Error
}

func (r *TransactionDebtRepo) DeleteByTransaction(ctx context.Context, txnID string) error {
	return r.db.WithContext(ctx).Where("transaction_id = ?", txnID).Delete(&models.TransactionDebt{}).Error
}
