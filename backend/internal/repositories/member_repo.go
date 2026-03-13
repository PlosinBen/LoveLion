package repositories

import (
	"context"

	"lovelion/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MemberRepo struct {
	db *gorm.DB
}

func NewMemberRepo(db *gorm.DB) *MemberRepo {
	return &MemberRepo{db: db}
}

func (r *MemberRepo) WithTx(tx *gorm.DB) *MemberRepo {
	return &MemberRepo{db: tx}
}

func (r *MemberRepo) Create(ctx context.Context, member *models.LedgerMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *MemberRepo) FindBySpaceAndUser(ctx context.Context, spaceID uuid.UUID, userID uuid.UUID) (*models.LedgerMember, error) {
	var member models.LedgerMember
	err := r.db.WithContext(ctx).
		Where("ledger_id = ? AND user_id = ?", spaceID, userID).
		First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *MemberRepo) FindBySpace(ctx context.Context, spaceID uuid.UUID) ([]models.LedgerMember, error) {
	var members []models.LedgerMember
	err := r.db.WithContext(ctx).
		Where("ledger_id = ?", spaceID).
		Preload("User").
		Find(&members).Error
	return members, err
}

func (r *MemberRepo) UpdateAlias(ctx context.Context, spaceID uuid.UUID, userID uuid.UUID, alias string) (int64, error) {
	result := r.db.WithContext(ctx).
		Model(&models.LedgerMember{}).
		Where("ledger_id = ? AND user_id = ?", spaceID, userID).
		Update("alias", alias)
	return result.RowsAffected, result.Error
}

func (r *MemberRepo) Delete(ctx context.Context, spaceID uuid.UUID, userID uuid.UUID) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("ledger_id = ? AND user_id = ?", spaceID, userID).
		Delete(&models.LedgerMember{})
	return result.RowsAffected, result.Error
}
