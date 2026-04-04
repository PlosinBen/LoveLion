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
	Date        *time.Time
	Currency    *string
	TotalAmount *decimal.Decimal
	Title       *string
	Note        *string
}

func (r *TransactionRepo) Create(ctx context.Context, txn *models.Transaction) error {
	return r.db.WithContext(ctx).Create(txn).Error
}

func (r *TransactionRepo) FindBySpace(ctx context.Context, spaceID uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.WithContext(ctx).
		Where("space_id = ?", spaceID).
		Preload("Expense").
		Preload("Expense.Items").
		Preload("Debts").
		Preload("Images", "entity_type = ?", "transaction").
		Order("date DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepo) FindBySpacePaginated(ctx context.Context, spaceID uuid.UUID, limit, offset int) ([]models.Transaction, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&models.Transaction{}).Where("space_id = ?", spaceID).Count(&total)

	var transactions []models.Transaction
	err := r.db.WithContext(ctx).
		Where("space_id = ?", spaceID).
		Preload("Expense").
		Preload("Expense.Items").
		Preload("Debts").
		Preload("Images", "entity_type = ?", "transaction").
		Order("date DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error
	return transactions, total, err
}

func (r *TransactionRepo) FindByID(ctx context.Context, id string, spaceID uuid.UUID) (*models.Transaction, error) {
	var txn models.Transaction
	err := r.db.WithContext(ctx).
		Where("id = ? AND space_id = ?", id, spaceID).
		Preload("Expense").
		Preload("Expense.Items").
		Preload("Debts").
		Preload("Images", "entity_type = ?", "transaction").
		First(&txn).Error
	if err != nil {
		return nil, err
	}
	return &txn, nil
}

func (r *TransactionRepo) Update(ctx context.Context, id string, params TransactionUpdateParams) error {
	updates := map[string]interface{}{}

	if params.Date != nil {
		updates["date"] = *params.Date
	}
	if params.Currency != nil {
		updates["currency"] = *params.Currency
	}
	if params.TotalAmount != nil {
		updates["total_amount"] = *params.TotalAmount
	}
	if params.Title != nil {
		updates["title"] = *params.Title
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
	result := r.db.WithContext(ctx).Where("id = ? AND space_id = ?", id, spaceID).Delete(&models.Transaction{})
	return result.RowsAffected, result.Error
}
