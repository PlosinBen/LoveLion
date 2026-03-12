package services

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"lovelion/internal/models"
	"lovelion/internal/utils/errorx"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InviteService struct {
	db *gorm.DB
}

func NewInviteService(db *gorm.DB) *InviteService {
	return &InviteService{db: db}
}

type CreateInviteParams struct {
	IsOneTime bool
	MaxUses   int
	ExpiresAt *time.Time
}

type InviteInfo struct {
	SpaceName   string `json:"space_name"`
	CreatorName string `json:"creator_name"`
	IsOneTime   bool   `json:"is_one_time"`
}

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// validateInvite checks expiration and usage limits.
// Shared by GetInviteInfo and JoinSpace to eliminate duplication.
func validateInvite(invite *models.LedgerInvite) error {
	if invite.ExpiresAt != nil && invite.ExpiresAt.Before(time.Now()) {
		return errorx.Wrap(errorx.ErrExpired, "Invite link has expired")
	}
	if invite.MaxUses > 0 && invite.UseCount >= invite.MaxUses {
		return errorx.Wrap(errorx.ErrExhausted, "Invite link has reached its maximum usage")
	}
	return nil
}

func (s *InviteService) Create(spaceID uuid.UUID, userID uuid.UUID, params CreateInviteParams) (*models.LedgerInvite, error) {
	invite := &models.LedgerInvite{
		ID:        uuid.New(),
		LedgerID:  spaceID,
		Token:     generateToken(),
		IsOneTime: params.IsOneTime,
		MaxUses:   params.MaxUses,
		ExpiresAt: params.ExpiresAt,
		CreatedBy: userID,
	}

	if invite.IsOneTime && invite.MaxUses <= 0 {
		invite.MaxUses = 1
	}

	if err := s.db.Create(invite).Error; err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to create invite")
	}

	return invite, nil
}

func (s *InviteService) GetInfo(token string) (*InviteInfo, error) {
	var invite models.LedgerInvite
	if err := s.db.Where("token = ?", token).Preload("Ledger").Preload("Creator").First(&invite).Error; err != nil {
		return nil, errorx.Wrap(errorx.ErrNotFound, "Invite link invalid or expired")
	}

	if err := validateInvite(&invite); err != nil {
		return nil, err
	}

	return &InviteInfo{
		SpaceName:   invite.Ledger.Name,
		CreatorName: invite.Creator.DisplayName,
		IsOneTime:   invite.IsOneTime,
	}, nil
}

func (s *InviteService) Join(token string, userID uuid.UUID) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var invite models.LedgerInvite
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("token = ?", token).First(&invite).Error; err != nil {
			return errorx.Wrap(errorx.ErrNotFound, "Invite link invalid")
		}

		if err := validateInvite(&invite); err != nil {
			return err
		}

		// Check if already a member
		var existing models.LedgerMember
		if err := tx.Where("ledger_id = ? AND user_id = ?", invite.LedgerID, userID).First(&existing).Error; err == nil {
			return nil // Already a member, no-op
		}

		member := &models.LedgerMember{
			ID:       uuid.New(),
			LedgerID: invite.LedgerID,
			UserID:   userID,
			Role:     "member",
		}

		if err := tx.Create(member).Error; err != nil {
			return errorx.Wrap(errorx.ErrInternal, "Failed to add member")
		}

		invite.UseCount++
		if err := tx.Save(&invite).Error; err != nil {
			return errorx.Wrap(errorx.ErrInternal, "Failed to update invite usage")
		}

		return nil
	})
}

func (s *InviteService) ListActive(spaceID uuid.UUID) ([]models.LedgerInvite, error) {
	var invites []models.LedgerInvite
	err := s.db.Where(
		"ledger_id = ? AND (expires_at IS NULL OR expires_at > ?) AND (max_uses = 0 OR use_count < max_uses)",
		spaceID, time.Now(),
	).Order("created_at DESC").Find(&invites).Error

	if err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to fetch invites")
	}

	return invites, nil
}

func (s *InviteService) Revoke(inviteID uuid.UUID, spaceID uuid.UUID) error {
	result := s.db.Where("id = ? AND ledger_id = ?", inviteID, spaceID).Delete(&models.LedgerInvite{})
	if result.Error != nil {
		return errorx.Wrap(errorx.ErrInternal, "Failed to revoke invite")
	}
	if result.RowsAffected == 0 {
		return errorx.Wrap(errorx.ErrNotFound, "Invite not found")
	}
	return nil
}
