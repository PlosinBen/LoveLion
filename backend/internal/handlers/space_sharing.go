package handlers

import (
	"net/http"
	"time"

	"lovelion/internal/models"
	"lovelion/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SpaceSharingHandler struct {
	db            *gorm.DB
	inviteService *services.InviteService
}

func NewSpaceSharingHandler(db *gorm.DB) *SpaceSharingHandler {
	return &SpaceSharingHandler{
		db:            db,
		inviteService: services.NewInviteService(db),
	}
}

type CreateInviteRequest struct {
	IsOneTime bool       `json:"is_one_time"`
	MaxUses   int        `json:"max_uses"`
	ExpiresAt *time.Time `json:"expires_at"`
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

	invite, err := h.inviteService.Create(space.ID, userID, services.CreateInviteParams{
		IsOneTime: req.IsOneTime,
		MaxUses:   req.MaxUses,
		ExpiresAt: req.ExpiresAt,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, invite)
}

// Get invite info (publicly accessible to show preview)
func (h *SpaceSharingHandler) GetInviteInfo(c *gin.Context) {
	token := c.Param("token")

	info, err := h.inviteService.GetInfo(token)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, info)
}

// Join a space via invite link
func (h *SpaceSharingHandler) JoinSpace(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	token := c.Param("token")

	if err := h.inviteService.Join(token, userID); err != nil {
		respondError(c, err)
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

// Update a member's alias
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

// Remove a member
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

// List all active invites for a space
func (h *SpaceSharingHandler) ListInvites(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Ledger)

	invites, err := h.inviteService.ListActive(space.ID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, invites)
}

// Revoke an invite link
func (h *SpaceSharingHandler) RevokeInvite(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Ledger)
	inviteID, err := uuid.Parse(c.Param("invite_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Invite ID"})
		return
	}

	if err := h.inviteService.Revoke(inviteID, space.ID); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invite revoked successfully"})
}
