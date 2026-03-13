package handlers

import (
	"net/http"
	"time"

	"lovelion/internal/models"
	"lovelion/internal/repositories"
	"lovelion/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SpaceSharingHandler struct {
	inviteService *services.InviteService
	memberRepo    *repositories.MemberRepo
}

func NewSpaceSharingHandler(inviteService *services.InviteService, memberRepo *repositories.MemberRepo) *SpaceSharingHandler {
	return &SpaceSharingHandler{
		inviteService: inviteService,
		memberRepo:    memberRepo,
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

	invite, err := h.inviteService.Create(c.Request.Context(), space.ID, userID, services.CreateInviteParams{
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

	info, err := h.inviteService.GetInfo(c.Request.Context(), token)
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

	if err := h.inviteService.Join(c.Request.Context(), token, userID); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully joined the space"})
}

// List all members of a space
func (h *SpaceSharingHandler) ListMembers(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Ledger)

	members, err := h.memberRepo.FindBySpace(c.Request.Context(), space.ID)
	if err != nil {
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

	rows, err := h.memberRepo.UpdateAlias(c.Request.Context(), space.ID, targetUserID, req.Alias)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update alias"})
		return
	}

	if rows == 0 {
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

	rows, err := h.memberRepo.Delete(c.Request.Context(), space.ID, targetUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove member"})
		return
	}

	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found in this space"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed successfully"})
}

// List all active invites for a space
func (h *SpaceSharingHandler) ListInvites(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Ledger)

	invites, err := h.inviteService.ListActive(c.Request.Context(), space.ID)
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

	if err := h.inviteService.Revoke(c.Request.Context(), inviteID, space.ID); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invite revoked successfully"})
}
