package handlers

import (
	"net/http"
	"strconv"
	"time"

	"lovelion/internal/models"
	"lovelion/internal/repositories"
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
// Supports optional pagination: ?limit=N&offset=M
// Without limit, returns all transactions (backward compatible).
// With limit, returns paginated results + X-Total-Count header.
func (h *TransactionHandler) List(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Space)

	limitStr := c.Query("limit")
	if limitStr == "" {
		transactions, err := h.svc.List(c.Request.Context(), space.ID)
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, transactions)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
			return
		}
	}

	filter := &repositories.TransactionFilter{
		Search:   c.Query("search"),
		Category: c.Query("category"),
		Type:     c.Query("type"),
	}
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		if t, err := time.Parse("2006-01-02", dateFrom); err == nil {
			filter.DateFrom = &t
		}
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		if t, err := time.Parse("2006-01-02", dateTo); err == nil {
			end := t.Add(24*time.Hour - time.Nanosecond)
			filter.DateTo = &end
		}
	}

	transactions, total, err := h.svc.ListPaginated(c.Request.Context(), space.ID, limit, offset, filter)
	if err != nil {
		respondError(c, err)
		return
	}

	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
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
