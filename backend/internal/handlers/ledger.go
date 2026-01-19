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
	Currencies     []string `json:"currencies"`
	Members        []string `json:"members"`
	Categories     []string `json:"categories"`
	PaymentMethods []string `json:"payment_methods"`
}

type UpdateLedgerRequest struct {
	Name           string   `json:"name" binding:"omitempty,min=1,max=100"`
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

// List user's ledgers
func (h *LedgerHandler) List(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var ledgers []models.Ledger
	if err := h.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&ledgers).Error; err != nil {
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
		ID:     uuid.New(),
		UserID: userID,
		Name:   req.Name,
		Type:   req.Type,
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
		ledger.Members = toJSON(req.Members)
	} else {
		ledger.Members = toJSON([]string{})
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

	if err := h.db.Create(ledger).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ledger"})
		return
	}

	c.JSON(http.StatusCreated, ledger)
}

// Get a single ledger
func (h *LedgerHandler) Get(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	ledgerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ledger ID"})
		return
	}

	var ledger models.Ledger
	if err := h.db.Where("id = ? AND user_id = ?", ledgerID, userID).First(&ledger).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ledger not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ledger"})
		return
	}

	c.JSON(http.StatusOK, ledger)
}

// Update a ledger
func (h *LedgerHandler) Update(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	ledgerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ledger ID"})
		return
	}

	var ledger models.Ledger
	if err := h.db.Where("id = ? AND user_id = ?", ledgerID, userID).First(&ledger).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ledger not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ledger"})
		return
	}

	var req UpdateLedgerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		ledger.Name = req.Name
	}
	if req.Currencies != nil {
		ledger.Currencies = toJSON(req.Currencies)
	}
	if req.Members != nil {
		ledger.Members = toJSON(req.Members)
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
	userID := c.MustGet("userID").(uuid.UUID)
	ledgerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ledger ID"})
		return
	}

	result := h.db.Where("id = ? AND user_id = ?", ledgerID, userID).Delete(&models.Ledger{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ledger"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ledger not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ledger deleted"})
}
