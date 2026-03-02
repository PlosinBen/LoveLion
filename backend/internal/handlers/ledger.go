package handlers

import (
	"encoding/json"
	"net/http"

	"lovelion/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type LedgerHandler struct {
	db *gorm.DB
}

func NewLedgerHandler(db *gorm.DB) *LedgerHandler {
	return &LedgerHandler{db: db}
}

type CreateLedgerRequest struct {
	Name           string   `json:"name" binding:"required,min=1,max=100"`
	Type           string   `json:"type"`
	BaseCurrency   string   `json:"base_currency"`
	Currencies     []string `json:"currencies"`
	Members        []string `json:"members"`
	Categories     []string `json:"categories"`
	PaymentMethods []string `json:"payment_methods"`
}

type UpdateLedgerRequest struct {
	Name           string   `json:"name" binding:"omitempty,min=1,max=100"`
	BaseCurrency   string   `json:"base_currency"`
	Currencies     []string `json:"currencies"`
	Members        []string `json:"members"`
	Categories     []string `json:"categories"`
	PaymentMethods []string `json:"payment_methods"`
}

func toJSON(v interface{}) datatypes.JSON {
	if v == nil {
		return datatypes.JSON("[]")
	}
	bytes, _ := json.Marshal(v)
	return datatypes.JSON(bytes)
}

// List user's ledgers (including shared ones)
func (h *LedgerHandler) List(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var ledgers []models.Ledger
	// Find ledgers where user is a member (either owner or member)
	err := h.db.
		Joins("JOIN ledger_members ON ledger_members.ledger_id = ledgers.id").
		Where("ledger_members.user_id = ?", userID).
		Preload("User"). // Preload owner info
		Order("ledgers.created_at DESC").
		Find(&ledgers).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ledgers"})
		return
	}

	c.JSON(http.StatusOK, ledgers)
}

// Create a new ledger
func (h *LedgerHandler) Create(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var req CreateLedgerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ledger := &models.Ledger{
		ID:           uuid.New(),
		UserID:       userID,
		Name:         req.Name,
		Type:         req.Type,
		BaseCurrency: req.BaseCurrency,
	}

	if ledger.BaseCurrency == "" {
		ledger.BaseCurrency = "TWD"
	}

	if ledger.Type == "" {
		ledger.Type = "personal"
	}

	if req.Currencies != nil {
		ledger.Currencies = toJSON(req.Currencies)
	} else {
		ledger.Currencies = toJSON([]string{"TWD"})
	}

	if req.Members != nil {
		ledger.MemberNames = toJSON(req.Members)
	} else {
		ledger.MemberNames = toJSON([]string{})
	}

	if req.Categories != nil {
		ledger.Categories = toJSON(req.Categories)
	} else {
		ledger.Categories = toJSON([]string{})
	}

	if req.PaymentMethods != nil {
		ledger.PaymentMethods = toJSON(req.PaymentMethods)
	} else {
		ledger.PaymentMethods = toJSON([]string{})
	}

	// Use a transaction to create both the ledger and the owner membership
	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(ledger).Error; err != nil {
			return err
		}

		// Add owner to members table
		member := &models.LedgerMember{
			ID:       uuid.New(),
			LedgerID: ledger.ID,
			UserID:   userID,
			Role:     "owner",
		}

		if err := tx.Create(member).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ledger"})
		return
	}

	c.JSON(http.StatusCreated, ledger)
}

// Get a single ledger
func (h *LedgerHandler) Get(c *gin.Context) {
	ledger, _ := c.Get("ledger")
	c.JSON(http.StatusOK, ledger)
}

// Update a ledger
func (h *LedgerHandler) Update(c *gin.Context) {
	ledgerVal, _ := c.Get("ledger")
	ledger := ledgerVal.(*models.Ledger)

	var req UpdateLedgerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		ledger.Name = req.Name
	}
	if req.BaseCurrency != "" {
		ledger.BaseCurrency = req.BaseCurrency
	}
	if req.Currencies != nil {
		ledger.Currencies = toJSON(req.Currencies)
	}
	if req.Members != nil {
		ledger.MemberNames = toJSON(req.Members)
	}
	if req.Categories != nil {
		ledger.Categories = toJSON(req.Categories)
	}
	if req.PaymentMethods != nil {
		ledger.PaymentMethods = toJSON(req.PaymentMethods)
	}

	if err := h.db.Save(&ledger).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ledger"})
		return
	}

	c.JSON(http.StatusOK, ledger)
}

// Delete a ledger
func (h *LedgerHandler) Delete(c *gin.Context) {
	ledgerVal, _ := c.Get("ledger")
	ledger := ledgerVal.(*models.Ledger)

	if err := h.db.Delete(ledger).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ledger"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ledger deleted"})
}
