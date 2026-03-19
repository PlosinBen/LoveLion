package repositories

import (
	"context"
	"time"

	"lovelion/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InviteRepo struct {
	db *gorm.DB
}

func NewInviteRepo(db *gorm.DB) *InviteRepo {
	return &InviteRepo{db: db}
}

func (r *InviteRepo) WithTx(tx *gorm.DB) *InviteRepo {
	return &InviteRepo{db: tx}
}

func (r *InviteRepo) Create(ctx context.Context, invite *models.SpaceInvite) error {
	return r.db.WithContext(ctx).Create(invite).Error
}

func (r *InviteRepo) FindByToken(ctx context.Context, token string) (*models.SpaceInvite, error) {
	var invite models.SpaceInvite
	err := r.db.WithContext(ctx).
		Where("token = ?", token).
		Preload("Space").
		Preload("Creator").
		First(&invite).Error
	if err != nil {
		return nil, err
	}
	return &invite, nil
}

func (r *InviteRepo) FindByTokenForUpdate(ctx context.Context, token string) (*models.SpaceInvite, error) {
	var invite models.SpaceInvite
	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("token = ?", token).
		First(&invite).Error
	if err != nil {
		return nil, err
	}
	return &invite, nil
}

func (r *InviteRepo) FindActiveBySpace(ctx context.Context, spaceID uuid.UUID) ([]models.SpaceInvite, error) {
	var invites []models.SpaceInvite
	err := r.db.WithContext(ctx).
		Where("space_id = ? AND (expires_at IS NULL OR expires_at > ?) AND (max_uses = 0 OR use_count < max_uses)", spaceID, time.Now()).
		Order("created_at DESC").
		Find(&invites).Error
	return invites, err
}

func (r *InviteRepo) IncrementUseCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.SpaceInvite{}).
		Where("id = ?", id).
		UpdateColumn("use_count", gorm.Expr("use_count + 1")).
		Error
}

func (r *InviteRepo) Delete(ctx context.Context, id uuid.UUID, spaceID uuid.UUID) (int64, error) {
	result := r.db.WithContext(ctx).Where("id = ? AND space_id = ?", id, spaceID).Delete(&models.SpaceInvite{})
	return result.RowsAffected, result.Error
}
