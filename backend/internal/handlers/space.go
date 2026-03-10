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
	Members        []string   `json:"members"` // legacy member names
	Categories     []string   `json:"categories"`
	PaymentMethods []string   `json:"payment_methods"`
	StartDate      *time.Time `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
	IsPinned       bool       `json:"is_pinned"`
}

type UpdateSpaceRequest struct {
	Name           string     `json:"name" binding:"omitempty,min=1,max=100"`
	Description    *string    `json:"description"`
	BaseCurrency   string     `json:"base_currency"`
	Currencies     []string   `json:"currencies"`
	Members        []string   `json:"members"`
	Categories     []string   `json:"categories"`
	PaymentMethods []string   `json:"payment_methods"`
	StartDate      *time.Time `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
	IsPinned       *bool      `json:"is_pinned"`
}

func toJSON(v interface{}) datatypes.JSON {
	if v == nil {
		return datatypes.JSON("[]")
	}
	bytes, _ := json.Marshal(v)
	return datatypes.JSON(bytes)
}

// List user's spaces
func (h *SpaceHandler) List(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	spaceType := c.Query("type")

	query := h.db.
		Joins("JOIN ledger_members ON ledger_members.ledger_id = ledgers.id").
		Where("ledger_members.user_id = ?", userID).
		Preload("User").
		Preload("Images", "entity_type = ?", "space")

	if spaceType != "" {
		query = query.Where("ledgers.type = ?", spaceType)
	}

	var spaces []models.Ledger
	err := query.Order("ledgers.is_pinned DESC, ledgers.created_at DESC").
		Find(&spaces).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch spaces"})
		return
	}

	for i := range spaces {
		spaces[i].PopulateCoverImage()
	}

	c.JSON(http.StatusOK, spaces)
}

// Create a new space
func (h *SpaceHandler) Create(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var req CreateSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	space := &models.Ledger{
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

	space.Currencies = toJSON(req.Currencies)
	if req.Currencies == nil {
		space.Currencies = toJSON([]string{space.BaseCurrency})
	}

	space.MemberNames = toJSON(req.Members)
	space.Categories = toJSON(req.Categories)
	space.PaymentMethods = toJSON(req.PaymentMethods)

	// Use a transaction to create both the space and the owner membership
	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(space).Error; err != nil {
			return err
		}

		// Add owner to members table
		member := &models.LedgerMember{
			ID:       uuid.New(),
			LedgerID: space.ID,
			UserID:   userID,
			Role:     "owner",
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
	space := spaceVal.(*models.Ledger)

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
	if req.Currencies != nil {
		space.Currencies = toJSON(req.Currencies)
	}
	if req.Members != nil {
		space.MemberNames = toJSON(req.Members)
	}
	if req.Categories != nil {
		space.Categories = toJSON(req.Categories)
	}
	if req.PaymentMethods != nil {
		space.PaymentMethods = toJSON(req.PaymentMethods)
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
	space := spaceVal.(*models.Ledger)

	if err := h.db.Delete(space).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete space"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Space deleted"})
}
