package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"lovelion/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SpaceSharingHandler struct {
	db *gorm.DB
}

func NewSpaceSharingHandler(db *gorm.DB) *SpaceSharingHandler {
	return &SpaceSharingHandler{db: db}
}

type CreateInviteRequest struct {
	IsOneTime bool       `json:"is_one_time"`
	MaxUses   int        `json:"max_uses"` // 0 for unlimited
	ExpiresAt *time.Time `json:"expires_at"`
}

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// Create an invite link
func (h *SpaceSharingHandler) CreateInvite(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Ledger)

	var req CreateInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invite := &models.LedgerInvite{
		ID:        uuid.New(),
		LedgerID:  space.ID,
		Token:     generateToken(),
		IsOneTime: req.IsOneTime,
		MaxUses:   req.MaxUses,
		ExpiresAt: req.ExpiresAt,
		CreatedBy: userID,
	}

	// If one-time, max_uses should be 1 if not specified
	if invite.IsOneTime && invite.MaxUses <= 0 {
		invite.MaxUses = 1
	}

	if err := h.db.Create(invite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invite"})
		return
	}

	c.JSON(http.StatusCreated, invite)
}

// Get invite info (publicly accessible to show preview)
func (h *SpaceSharingHandler) GetInviteInfo(c *gin.Context) {
	token := c.Param("token")

	var invite models.LedgerInvite
	if err := h.db.Where("token = ?", token).Preload("Ledger").Preload("Creator").First(&invite).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invite link invalid or expired"})
		return
	}

	// Check expiration
	if invite.ExpiresAt != nil && invite.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusGone, gin.H{"error": "Invite link has expired"})
		return
	}

	// Check uses
	if invite.MaxUses > 0 && invite.UseCount >= invite.MaxUses {
		c.JSON(http.StatusGone, gin.H{"error": "Invite link has reached its maximum usage"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"space_name":   invite.Ledger.Name,
		"creator_name": invite.Creator.DisplayName,
		"is_one_time":  invite.IsOneTime,
	})
}

// Join a space via invite link
func (h *SpaceSharingHandler) JoinSpace(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	token := c.Param("token")

	err := h.db.Transaction(func(tx *gorm.DB) error {
		var invite models.LedgerInvite
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("token = ?", token).First(&invite).Error; err != nil {
			return gorm.ErrRecordNotFound
		}

		// Check expiration
		if invite.ExpiresAt != nil && invite.ExpiresAt.Before(time.Now()) {
			return gorm.ErrInvalidData // Custom error would be better
		}

		// Check uses
		if invite.MaxUses > 0 && invite.UseCount >= invite.MaxUses {
			return gorm.ErrInvalidData
		}

		// Check if already a member
		var existingMember models.LedgerMember
		if err := tx.Where("ledger_id = ? AND user_id = ?", invite.LedgerID, userID).First(&existingMember).Error; err == nil {
			// Already a member, success (no-op)
			return nil
		}

		// Add as member
		member := &models.LedgerMember{
			ID:       uuid.New(),
			LedgerID: invite.LedgerID,
			UserID:   userID,
			Role:     "member",
		}

		if err := tx.Create(member).Error; err != nil {
			return err
		}

		// Update use count
		invite.UseCount++
		if err := tx.Save(&invite).Error; err != nil {
			return err
		}

		// Pre-reserve for "NotifyMemberJoined" logic later
		// h.NotifyMemberJoined(invite.LedgerID, userID)

		return nil
	})

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invite link invalid"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to join space: invite link may be expired or full"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully joined the space"})
}

// List all members of a space
func (h *SpaceSharingHandler) ListMembers(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Ledger)

	var members []models.LedgerMember
	if err := h.db.Where("ledger_id = ?", space.ID).Preload("User").Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch members"})
		return
	}

	c.JSON(http.StatusOK, members)
}

type UpdateMemberRequest struct {
	Alias string `json:"alias"`
}

// Update a member's alias (Owner only via SpaceOwnerOnly middleware)
func (h *SpaceSharingHandler) UpdateMemberAlias(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Ledger)
	targetUserID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	var req UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := h.db.Model(&models.LedgerMember{}).
		Where("ledger_id = ? AND user_id = ?", space.ID, targetUserID).
		Update("alias", req.Alias)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update alias"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found in this space"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Alias updated successfully"})
}

// Remove a member (Kick by owner OR leave by self)
func (h *SpaceSharingHandler) RemoveMember(c *gin.Context) {
	requestorID := c.MustGet("userID").(uuid.UUID)
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Ledger)
	memberVal, _ := c.Get("member")
	requestorMember := memberVal.(*models.LedgerMember)
	
	targetUserID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	// Permission logic:
	// 1. If requestor is Owner, they can remove anyone EXCEPT themselves (to delete space, use DeleteSpace)
	// 2. If requestor is NOT Owner, they can ONLY remove themselves (leaving)
	
	isOwner := requestorMember.Role == "owner"
	isSelf := requestorID == targetUserID

	if !isOwner && !isSelf {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to remove this member"})
		return
	}

	if isOwner && isSelf {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Owner cannot leave their own space. Please delete the space instead."})
		return
	}

	result := h.db.Where("ledger_id = ? AND user_id = ?", space.ID, targetUserID).Delete(&models.LedgerMember{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove member"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found in this space"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed successfully"})
}

// List all active invites for a space (Owner only via SpaceOwnerOnly)
func (h *SpaceSharingHandler) ListInvites(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Ledger)

	var invites []models.LedgerInvite
	// Only show active ones: not expired AND (max_uses=0 OR use_count < max_uses)
	err := h.db.Where("ledger_id = ? AND (expires_at IS NULL OR expires_at > ?) AND (max_uses = 0 OR use_count < max_uses)", space.ID, time.Now()).
		Order("created_at DESC").
		Find(&invites).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch invites"})
		return
	}

	c.JSON(http.StatusOK, invites)
}

// Revoke an invite link (Owner only via SpaceOwnerOnly)
func (h *SpaceSharingHandler) RevokeInvite(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Ledger)
	inviteID, err := uuid.Parse(c.Param("invite_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Invite ID"})
		return
	}

	result := h.db.Where("id = ? AND ledger_id = ?", inviteID, space.ID).Delete(&models.LedgerInvite{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke invite"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invite not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invite revoked successfully"})
}
