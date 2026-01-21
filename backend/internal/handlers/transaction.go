package handlers

import (
	"net/http"
	"time"

	"lovelion/internal/models"
	"lovelion/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TransactionHandler struct {
	db *gorm.DB
}

func NewTransactionHandler(db *gorm.DB) *TransactionHandler {
	return &TransactionHandler{db: db}
}

type TransactionItemRequest struct {
	Name      string          `json:"name" binding:"required"`
	UnitPrice decimal.Decimal `json:"unit_price"`
	Quantity  decimal.Decimal `json:"quantity"`
	Discount  decimal.Decimal `json:"discount"`
}

type TransactionSplitRequest struct {
	MemberID uuid.UUID       `json:"member_id" binding:"required"`
	Amount   decimal.Decimal `json:"amount"`
	IsPayer  bool            `json:"is_payer"`
}

type CreateTransactionRequest struct {
	Payer         string                    `json:"payer"`
	Date          *time.Time                `json:"date"`
	Currency      string                    `json:"currency"`
	ExchangeRate  decimal.Decimal           `json:"exchange_rate"`
	BillingAmount decimal.Decimal           `json:"billing_amount"`
	HandlingFee   decimal.Decimal           `json:"handling_fee"`
	Category      string                    `json:"category"`
	PaymentMethod string                    `json:"payment_method"`
	Note          string                    `json:"note"`
	Items         []TransactionItemRequest  `json:"items"`
	Splits        []TransactionSplitRequest `json:"splits"`
}

type UpdateTransactionRequest struct {
	Payer         string                    `json:"payer"`
	Date          *time.Time                `json:"date"`
	Currency      string                    `json:"currency"`
	ExchangeRate  *decimal.Decimal          `json:"exchange_rate"`
	BillingAmount *decimal.Decimal          `json:"billing_amount"`
	HandlingFee   *decimal.Decimal          `json:"handling_fee"`
	Category      string                    `json:"category"`
	PaymentMethod string                    `json:"payment_method"`
	Note          string                    `json:"note"`
	Items         []TransactionItemRequest  `json:"items"`
	Splits        []TransactionSplitRequest `json:"splits"`
}

// Helper to verify ledger ownership
// Helper to verify ledger ownership or access
func (h *TransactionHandler) verifyLedgerOwnership(ledgerID uuid.UUID, userID uuid.UUID) (*models.Ledger, error) {
	var ledger models.Ledger
	if err := h.db.Where("id = ?", ledgerID).First(&ledger).Error; err != nil {
		return nil, err
	}

	// Personal ledger: must match UserID
	if ledger.Type == "personal" {
		if ledger.UserID != userID {
			return nil, gorm.ErrRecordNotFound
		}
		return &ledger, nil
	}

	// Trip ledger: check if user is in the trip
	if ledger.Type == "trip" {
		var trip models.Trip
		if err := h.db.Where("ledger_id = ?", ledgerID).Preload("Members").First(&trip).Error; err != nil {
			return nil, err
		}

		// Check if creator
		if trip.CreatedBy == userID {
			return &ledger, nil
		}

		// Check members
		for _, m := range trip.Members {
			if m.UserID != nil && *m.UserID == userID {
				return &ledger, nil
			}
		}

		return nil, gorm.ErrRecordNotFound
	}

	return nil, gorm.ErrRecordNotFound
}

// List transactions for a ledger
func (h *TransactionHandler) List(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	ledgerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ledger ID"})
		return
	}

	// Verify ownership
	if _, err := h.verifyLedgerOwnership(ledgerID, userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ledger not found"})
		return
	}

	var transactions []models.Transaction
	if err := h.db.Where("ledger_id = ?", ledgerID).
		Preload("Items").
		Preload("Splits").
		Order("date DESC").
		Order("date DESC").
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// Create a new transaction with items
func (h *TransactionHandler) Create(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	ledgerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ledger ID"})
		return
	}

	// Verify ownership
	if _, err := h.verifyLedgerOwnership(ledgerID, userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ledger not found"})
		return
	}

	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate short ID for transaction
	txnID, err := utils.NewShortID(h.db, "transactions", "id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate ID"})
		return
	}

	txn := &models.Transaction{
		ID:            txnID,
		LedgerID:      ledgerID,
		Payer:         req.Payer,
		Currency:      req.Currency,
		ExchangeRate:  req.ExchangeRate,
		BillingAmount: req.BillingAmount,
		HandlingFee:   req.HandlingFee,
		Category:      req.Category,
		PaymentMethod: req.PaymentMethod,
		Note:          req.Note,
	}

	if req.Date != nil {
		txn.Date = *req.Date
	} else {
		txn.Date = time.Now()
	}

	if txn.Currency == "" {
		txn.Currency = "TWD"
	}

	if txn.ExchangeRate.IsZero() {
		txn.ExchangeRate = decimal.NewFromInt(1)
	}

	// Calculate items and total
	totalAmount := decimal.Zero
	var items []models.TransactionItem
	for _, itemReq := range req.Items {
		quantity := itemReq.Quantity
		if quantity.IsZero() {
			quantity = decimal.NewFromInt(1)
		}

		amount := itemReq.UnitPrice.Mul(quantity).Sub(itemReq.Discount)

		item := models.TransactionItem{
			ID:            uuid.New(),
			TransactionID: txnID,
			Name:          itemReq.Name,
			UnitPrice:     itemReq.UnitPrice,
			Quantity:      quantity,
			Discount:      itemReq.Discount,
			Amount:        amount,
		}
		items = append(items, item)
		totalAmount = totalAmount.Add(amount)
	}

	txn.TotalAmount = totalAmount
	txn.Items = items

	// Calculate splits
	var splits []models.TransactionSplit
	for _, splitReq := range req.Splits {
		split := models.TransactionSplit{
			ID:            uuid.New(),
			TransactionID: txnID,
			MemberID:      splitReq.MemberID,
			Amount:        splitReq.Amount,
			IsPayer:       splitReq.IsPayer,
		}
		splits = append(splits, split)
	}
	txn.Splits = splits

	// Use transaction to save both
	if err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(txn).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	c.JSON(http.StatusCreated, txn)
}

// Get a single transaction
func (h *TransactionHandler) Get(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	ledgerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ledger ID"})
		return
	}
	txnID := c.Param("txn_id")

	// Verify ownership
	if _, err := h.verifyLedgerOwnership(ledgerID, userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ledger not found"})
		return
	}

	var txn models.Transaction
	if err := h.db.Where("id = ? AND ledger_id = ?", txnID, ledgerID).
		Preload("Items").
		Preload("Splits").
		First(&txn).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction"})
		return
	}

	c.JSON(http.StatusOK, txn)
}

// Update a transaction
func (h *TransactionHandler) Update(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	ledgerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ledger ID"})
		return
	}
	txnID := c.Param("txn_id")

	// Verify ownership
	if _, err := h.verifyLedgerOwnership(ledgerID, userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ledger not found"})
		return
	}

	var txn models.Transaction
	if err := h.db.Where("id = ? AND ledger_id = ?", txnID, ledgerID).First(&txn).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction"})
		return
	}

	var req UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	if req.Payer != "" {
		txn.Payer = req.Payer
	}
	if req.Date != nil {
		txn.Date = *req.Date
	}
	if req.Currency != "" {
		txn.Currency = req.Currency
	}
	if req.ExchangeRate != nil {
		txn.ExchangeRate = *req.ExchangeRate
	}
	if req.BillingAmount != nil {
		txn.BillingAmount = *req.BillingAmount
	}
	if req.HandlingFee != nil {
		txn.HandlingFee = *req.HandlingFee
	}
	if req.Category != "" {
		txn.Category = req.Category
	}
	if req.PaymentMethod != "" {
		txn.PaymentMethod = req.PaymentMethod
	}
	txn.Note = req.Note

	// If items are provided, replace them
	if req.Items != nil {
		// Delete existing items
		h.db.Where("transaction_id = ?", txnID).Delete(&models.TransactionItem{})

		// Create new items
		totalAmount := decimal.Zero
		var items []models.TransactionItem
		for _, itemReq := range req.Items {
			quantity := itemReq.Quantity
			if quantity.IsZero() {
				quantity = decimal.NewFromInt(1)
			}

			amount := itemReq.UnitPrice.Mul(quantity).Sub(itemReq.Discount)

			item := models.TransactionItem{
				ID:            uuid.New(),
				TransactionID: txnID,
				Name:          itemReq.Name,
				UnitPrice:     itemReq.UnitPrice,
				Quantity:      quantity,
				Discount:      itemReq.Discount,
				Amount:        amount,
			}
			items = append(items, item)
			totalAmount = totalAmount.Add(amount)
		}

		txn.TotalAmount = totalAmount
		txn.Items = items
	}

	// If splits are provided, replace them
	if req.Splits != nil {
		// Delete existing splits
		h.db.Where("transaction_id = ?", txnID).Delete(&models.TransactionSplit{})

		var splits []models.TransactionSplit
		for _, splitReq := range req.Splits {
			split := models.TransactionSplit{
				ID:            uuid.New(),
				TransactionID: txnID,
				MemberID:      splitReq.MemberID,
				Amount:        splitReq.Amount,
				IsPayer:       splitReq.IsPayer,
			}
			splits = append(splits, split)
		}
		txn.Splits = splits
	}

	if err := h.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&txn).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction"})
		return
	}

	// Reload with items and splits
	h.db.Preload("Items").Preload("Splits").First(&txn, "id = ?", txnID)

	c.JSON(http.StatusOK, txn)
}

// Delete a transaction
func (h *TransactionHandler) Delete(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	ledgerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ledger ID"})
		return
	}
	txnID := c.Param("txn_id")

	// Verify ownership
	if _, err := h.verifyLedgerOwnership(ledgerID, userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ledger not found"})
		return
	}

	result := h.db.Where("id = ? AND ledger_id = ?", txnID, ledgerID).Delete(&models.Transaction{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted"})
}
