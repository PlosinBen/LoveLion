package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"lovelion/internal/models"
	"lovelion/internal/repositories"
	"lovelion/internal/utils/errorx"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InviteService struct {
	db         *gorm.DB
	inviteRepo *repositories.InviteRepo
	memberRepo *repositories.MemberRepo
}

func NewInviteService(db *gorm.DB, inviteRepo *repositories.InviteRepo, memberRepo *repositories.MemberRepo) *InviteService {
	return &InviteService{
		db:         db,
		inviteRepo: inviteRepo,
		memberRepo: memberRepo,
	}
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

func validateInvite(invite *models.SpaceInvite) error {
	if invite.ExpiresAt != nil && invite.ExpiresAt.Before(time.Now()) {
		return errorx.Wrap(errorx.ErrExpired, "Invite link has expired")
	}
	if invite.MaxUses > 0 && invite.UseCount >= invite.MaxUses {
		return errorx.Wrap(errorx.ErrExhausted, "Invite link has reached its maximum usage")
	}
	return nil
}

func (s *InviteService) Create(ctx context.Context, spaceID uuid.UUID, userID uuid.UUID, params CreateInviteParams) (*models.SpaceInvite, error) {
	invite := &models.SpaceInvite{
		ID:        uuid.New(),
		SpaceID:   spaceID,
		Token:     generateToken(),
		IsOneTime: params.IsOneTime,
		MaxUses:   params.MaxUses,
		ExpiresAt: params.ExpiresAt,
		CreatedBy: userID,
	}

	if invite.IsOneTime && invite.MaxUses <= 0 {
		invite.MaxUses = 1
	}

	if err := s.inviteRepo.Create(ctx, invite); err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to create invite")
	}

	return invite, nil
}

func (s *InviteService) GetInfo(ctx context.Context, token string) (*InviteInfo, error) {
	invite, err := s.inviteRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, errorx.Wrap(errorx.ErrNotFound, "Invite link invalid or expired")
	}

	if err := validateInvite(invite); err != nil {
		return nil, err
	}

	return &InviteInfo{
		SpaceName:   invite.Space.Name,
		CreatorName: invite.Creator.DisplayName,
		IsOneTime:   invite.IsOneTime,
	}, nil
}

func (s *InviteService) Join(ctx context.Context, token string, userID uuid.UUID) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		inviteRepo := s.inviteRepo.WithTx(tx)
		memberRepo := s.memberRepo.WithTx(tx)

		invite, err := inviteRepo.FindByTokenForUpdate(ctx, token)
		if err != nil {
			return errorx.Wrap(errorx.ErrNotFound, "Invite link invalid")
		}

		if err := validateInvite(invite); err != nil {
			return err
		}

		// Check if already a member
		if _, err := memberRepo.FindBySpaceAndUser(ctx, invite.SpaceID, userID); err == nil {
			return nil // Already a member, no-op
		}

		member := &models.SpaceMember{
			ID:      uuid.New(),
			SpaceID: invite.SpaceID,
			UserID:  userID,
			Role:    "member",
		}

		if err := memberRepo.Create(ctx, member); err != nil {
			return errorx.Wrap(errorx.ErrInternal, "Failed to add member")
		}

		if err := inviteRepo.IncrementUseCount(ctx, invite.ID); err != nil {
			return errorx.Wrap(errorx.ErrInternal, "Failed to update invite usage")
		}

		return nil
	})
}

func (s *InviteService) ListActive(ctx context.Context, spaceID uuid.UUID) ([]models.SpaceInvite, error) {
	invites, err := s.inviteRepo.FindActiveBySpace(ctx, spaceID)
	if err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, "Failed to fetch invites")
	}
	return invites, nil
}

func (s *InviteService) Revoke(ctx context.Context, inviteID uuid.UUID, spaceID uuid.UUID) error {
	rows, err := s.inviteRepo.Delete(ctx, inviteID, spaceID)
	if err != nil {
		return errorx.Wrap(errorx.ErrInternal, "Failed to revoke invite")
	}
	if rows == 0 {
		return errorx.Wrap(errorx.ErrNotFound, "Invite not found")
	}
	return nil
}
