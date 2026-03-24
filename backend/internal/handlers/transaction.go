package handlers

import (
	"net/http"

	"lovelion/internal/models"
	"lovelion/internal/services"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	svc *services.TransactionService
}

func NewTransactionHandler(svc *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: svc}
}

// List transactions for a space (all types)
func (h *TransactionHandler) List(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Space)

	transactions, err := h.svc.List(c.Request.Context(), space.ID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// Get a single transaction (any type)
func (h *TransactionHandler) Get(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Space)
	txnID := c.Param("txn_id")

	txn, err := h.svc.GetByID(c.Request.Context(), txnID, space.ID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, txn)
}

// Delete a transaction (any type)
func (h *TransactionHandler) Delete(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Space)
	txnID := c.Param("txn_id")

	if err := h.svc.Delete(c.Request.Context(), txnID, space.ID); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted"})
}
