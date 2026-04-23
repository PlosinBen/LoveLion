package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"lovelion/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SpaceHandler struct {
	db *gorm.DB
}

func NewSpaceHandler(db *gorm.DB) *SpaceHandler {
	return &SpaceHandler{db: db}
}

type CreateSpaceRequest struct {
	Name           string     `json:"name" binding:"required,min=1,max=100"`
	Description    string     `json:"description"`
	Type           string     `json:"type"` // personal, trip, group
	BaseCurrency   string     `json:"base_currency"`
	Currencies     []string   `json:"currencies"`
	SplitMembers   []string   `json:"split_members"`
	Categories     []string   `json:"categories"`
	PaymentMethods []string   `json:"payment_methods"`
	StartDate      *time.Time `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
	IsPinned       bool       `json:"is_pinned"`
}

type UpdateSpaceRequest struct {
	Name           string      `json:"name" binding:"omitempty,min=1,max=100"`
	Description    *string     `json:"description"`
	BaseCurrency   string      `json:"base_currency"`
	Currencies     *[]string   `json:"currencies"`
	SplitMembers   *[]string   `json:"split_members"`
	Categories     *[]string   `json:"categories"`
	PaymentMethods *[]string   `json:"payment_methods"`
	StartDate      *time.Time  `json:"start_date"`
	EndDate        *time.Time  `json:"end_date"`
	IsPinned       *bool       `json:"is_pinned"`
}

func toJSON(v interface{}) (datatypes.JSON, error) {
	if v == nil {
		return datatypes.JSON("[]"), nil
	}
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(bytes), nil
}

// List user's spaces
func (h *SpaceHandler) List(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	spaceType := c.Query("type")

	query := h.db.
		Joins("JOIN space_members ON space_members.space_id = spaces.id").
		Where("space_members.user_id = ?", userID).
		Preload("User").
		Preload("Images", "entity_type = ?", "space").
		Preload("Members")

	if spaceType != "" {
		query = query.Where("spaces.type = ?", spaceType)
	}

	var spaces []models.Space
	err := query.Order("spaces.is_pinned DESC, spaces.created_at DESC").
		Find(&spaces).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch spaces"})
		return
	}

	type spaceResponse struct {
		models.Space
		MyRole      string `json:"my_role"`
		MemberCount int    `json:"member_count"`
	}

	var results []spaceResponse
	for i := range spaces {
		spaces[i].PopulateCoverImage()

		role := ""
		memberCount := 0
		for _, m := range spaces[i].Members {
			memberCount++
			if m.UserID == userID {
				role = m.Role
			}
		}

		spaces[i].Members = nil // don't leak full member list
		results = append(results, spaceResponse{
			Space:       spaces[i],
			MyRole:      role,
			MemberCount: memberCount,
		})
	}

	c.JSON(http.StatusOK, results)
}

// Create a new space
func (h *SpaceHandler) Create(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var req CreateSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	space := &models.Space{
		ID:           uuid.New(),
		UserID:       userID,
		Name:         req.Name,
		Description:  req.Description,
		Type:         req.Type,
		BaseCurrency: req.BaseCurrency,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		IsPinned:     req.IsPinned,
	}

	if space.BaseCurrency == "" {
		space.BaseCurrency = "TWD"
	}

	if space.Type == "" {
		space.Type = "personal"
	}

	currencies := req.Currencies
	if currencies == nil {
		currencies = []string{space.BaseCurrency}
	}

	jsonFields := []struct {
		target *datatypes.JSON
		value  interface{}
	}{
		{&space.Currencies, currencies},
		{&space.SplitMembers, req.SplitMembers},
		{&space.Categories, req.Categories},
		{&space.PaymentMethods, req.PaymentMethods},
	}
	for _, f := range jsonFields {
		data, err := toJSON(f.value)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode space data"})
			return
		}
		*f.target = data
	}

	// Use a transaction to create both the space and the owner membership
	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(space).Error; err != nil {
			return err
		}

		// Add owner to members table
		member := &models.SpaceMember{
			ID:      uuid.New(),
			SpaceID: space.ID,
			UserID:  userID,
			Role:    "owner",
		}

		if err := tx.Create(member).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create space"})
		return
	}

	c.JSON(http.StatusCreated, space)
}

// Get a single space
func (h *SpaceHandler) Get(c *gin.Context) {
	space, _ := c.Get("space")
	c.JSON(http.StatusOK, space)
}

// Update a space
func (h *SpaceHandler) Update(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Space)

	var req UpdateSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		space.Name = req.Name
	}
	if req.Description != nil {
		space.Description = *req.Description
	}
	if req.BaseCurrency != "" {
		space.BaseCurrency = req.BaseCurrency
	}
	jsonUpdates := []struct {
		target *datatypes.JSON
		value  *[]string
	}{
		{&space.Currencies, req.Currencies},
		{&space.SplitMembers, req.SplitMembers},
		{&space.Categories, req.Categories},
		{&space.PaymentMethods, req.PaymentMethods},
	}
	for _, f := range jsonUpdates {
		if f.value == nil {
			continue
		}
		data, err := toJSON(*f.value)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode space data"})
			return
		}
		*f.target = data
	}
	if req.StartDate != nil {
		space.StartDate = req.StartDate
	}
	if req.EndDate != nil {
		space.EndDate = req.EndDate
	}
	if req.IsPinned != nil {
		space.IsPinned = *req.IsPinned
	}

	if err := h.db.Save(&space).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update space"})
		return
	}

	// Reload images for response
	h.db.Where("entity_id = ? AND entity_type = ?", space.ID, "space").Find(&space.Images)
	space.PopulateCoverImage()

	c.JSON(http.StatusOK, space)
}

// Delete a space
func (h *SpaceHandler) Delete(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Space)

	if err := h.db.Delete(space).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete space"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Space deleted"})
}

// Leave a space
func (h *SpaceHandler) Leave(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Space)

	// Check if user is a member
	var member models.SpaceMember
	if err := h.db.Where("space_id = ? AND user_id = ?", space.ID, userID).First(&member).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found in this space"})
		return
	}

	// Count total members
	var memberCount int64
	h.db.Model(&models.SpaceMember{}).Where("space_id = ?", space.ID).Count(&memberCount)

	// If owner and only member, suggest deleting instead
	if member.Role == "owner" && memberCount == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You are the only member and owner. Please delete the space instead of leaving."})
		return
	}

	// If owner but there are others, force transfer ownership or block
	if member.Role == "owner" && memberCount > 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You are the owner. Please transfer ownership before leaving or delete the space."})
		return
	}

	// Proceed to remove membership
	if err := h.db.Delete(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to leave space"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully left the space"})
}
