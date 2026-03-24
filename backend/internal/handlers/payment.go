package handlers

import (
	"net/http"
	"time"

	"lovelion/internal/models"
	"lovelion/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type PaymentHandler struct {
	svc *services.TransactionService
}

func NewPaymentHandler(svc *services.TransactionService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

type CreatePaymentRequest struct {
	Date        *time.Time      `json:"date"`
	Title       string          `json:"title"`
	Note        string          `json:"note"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	PayerName   string          `json:"payer_name" binding:"required"`
	PayeeName   string          `json:"payee_name" binding:"required"`
}

type UpdatePaymentRequest struct {
	Date        *time.Time       `json:"date"`
	Title       string           `json:"title"`
	Note        string           `json:"note"`
	TotalAmount *decimal.Decimal `json:"total_amount"`
	PayerName   string           `json:"payer_name" binding:"required"`
	PayeeName   string           `json:"payee_name" binding:"required"`
}

func (h *PaymentHandler) Create(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Space)

	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txn, err := h.svc.CreatePayment(c.Request.Context(), space.ID, space.BaseCurrency, services.CreatePaymentInput{
		Date:        req.Date,
		Title:       req.Title,
		Note:        req.Note,
		TotalAmount: req.TotalAmount,
		PayerName:   req.PayerName,
		PayeeName:   req.PayeeName,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, txn)
}

func (h *PaymentHandler) Update(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Space)
	txnID := c.Param("txn_id")

	var req UpdatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txn, err := h.svc.UpdatePayment(c.Request.Context(), txnID, space.ID, services.UpdatePaymentInput{
		Date:        req.Date,
		Title:       req.Title,
		Note:        req.Note,
		TotalAmount: req.TotalAmount,
		PayerName:   req.PayerName,
		PayeeName:   req.PayeeName,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, txn)
}
