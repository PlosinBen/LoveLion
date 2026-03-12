package handlers

import (
	"net/http"
	"time"

	"lovelion/internal/models"
	"lovelion/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TransactionHandler struct {
	svc *services.TransactionService
}

func NewTransactionHandler(db *gorm.DB) *TransactionHandler {
	return &TransactionHandler{svc: services.NewTransactionService(db)}
}

type TransactionItemRequest struct {
	Name      string          `json:"name" binding:"required"`
	UnitPrice decimal.Decimal `json:"unit_price"`
	Quantity  decimal.Decimal `json:"quantity"`
	Discount  decimal.Decimal `json:"discount"`
}

type TransactionSplitRequest struct {
	MemberID *uuid.UUID      `json:"member_id"`
	Name     string          `json:"name" binding:"required"`
	Amount   decimal.Decimal `json:"amount"`
	IsPayer  bool            `json:"is_payer"`
}

type CreateTransactionRequest struct {
	Payer         string                    `json:"payer"`
	Date          *time.Time                `json:"date"`
	Currency      string                    `json:"currency"`
	TotalAmount   decimal.Decimal           `json:"total_amount"`
	ExchangeRate  decimal.Decimal           `json:"exchange_rate"`
	BillingAmount decimal.Decimal           `json:"billing_amount"`
	HandlingFee   decimal.Decimal           `json:"handling_fee"`
	Category      string                    `json:"category"`
	Title         string                    `json:"title"`
	PaymentMethod string                    `json:"payment_method"`
	Note          string                    `json:"note"`
	Items         []TransactionItemRequest  `json:"items"`
	Splits        []TransactionSplitRequest `json:"splits"`
}

type UpdateTransactionRequest struct {
	Payer         string                    `json:"payer"`
	Date          *time.Time                `json:"date"`
	Currency      string                    `json:"currency"`
	TotalAmount   *decimal.Decimal          `json:"total_amount"`
	ExchangeRate  *decimal.Decimal          `json:"exchange_rate"`
	BillingAmount *decimal.Decimal          `json:"billing_amount"`
	HandlingFee   *decimal.Decimal          `json:"handling_fee"`
	Category      string                    `json:"category"`
	Title         string                    `json:"title"`
	PaymentMethod string                    `json:"payment_method"`
	Note          string                    `json:"note"`
	Items         []TransactionItemRequest  `json:"items"`
	Splits        []TransactionSplitRequest `json:"splits"`
}

func toItemInputs(reqs []TransactionItemRequest) []services.TransactionItemInput {
	if reqs == nil {
		return nil
	}
	inputs := make([]services.TransactionItemInput, len(reqs))
	for i, r := range reqs {
		inputs[i] = services.TransactionItemInput{
			Name:      r.Name,
			UnitPrice: r.UnitPrice,
			Quantity:  r.Quantity,
			Discount:  r.Discount,
		}
	}
	return inputs
}

func toSplitInputs(reqs []TransactionSplitRequest) []services.TransactionSplitInput {
	if reqs == nil {
		return nil
	}
	inputs := make([]services.TransactionSplitInput, len(reqs))
	for i, r := range reqs {
		inputs[i] = services.TransactionSplitInput{
			MemberID: r.MemberID,
			Name:     r.Name,
			Amount:   r.Amount,
			IsPayer:  r.IsPayer,
		}
	}
	return inputs
}

// List transactions for a space
func (h *TransactionHandler) List(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Ledger)

	transactions, err := h.svc.List(space.ID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// Create a new transaction with items
func (h *TransactionHandler) Create(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Ledger)

	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txn, err := h.svc.Create(space.ID, services.CreateTransactionInput{
		Payer:         req.Payer,
		Date:          req.Date,
		Currency:      req.Currency,
		TotalAmount:   req.TotalAmount,
		ExchangeRate:  req.ExchangeRate,
		BillingAmount: req.BillingAmount,
		HandlingFee:   req.HandlingFee,
		Category:      req.Category,
		Title:         req.Title,
		PaymentMethod: req.PaymentMethod,
		Note:          req.Note,
		Items:         toItemInputs(req.Items),
		Splits:        toSplitInputs(req.Splits),
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, txn)
}

// Get a single transaction
func (h *TransactionHandler) Get(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Ledger)
	txnID := c.Param("txn_id")

	txn, err := h.svc.GetByID(txnID, space.ID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, txn)
}

// Update a transaction
func (h *TransactionHandler) Update(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Ledger)
	txnID := c.Param("txn_id")

	var req UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txn, err := h.svc.Update(txnID, space.ID, services.UpdateTransactionInput{
		Payer:         req.Payer,
		Date:          req.Date,
		Currency:      req.Currency,
		TotalAmount:   req.TotalAmount,
		ExchangeRate:  req.ExchangeRate,
		BillingAmount: req.BillingAmount,
		HandlingFee:   req.HandlingFee,
		Category:      req.Category,
		Title:         req.Title,
		PaymentMethod: req.PaymentMethod,
		Note:          req.Note,
		Items:         toItemInputs(req.Items),
		Splits:        toSplitInputs(req.Splits),
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, txn)
}

// Delete a transaction
func (h *TransactionHandler) Delete(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Ledger)
	txnID := c.Param("txn_id")

	if err := h.svc.Delete(txnID, space.ID); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted"})
}
