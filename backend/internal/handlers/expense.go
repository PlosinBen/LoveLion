package handlers

import (
	"net/http"
	"time"

	"lovelion/internal/models"
	"lovelion/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type ExpenseHandler struct {
	svc *services.TransactionService
}

func NewExpenseHandler(svc *services.TransactionService) *ExpenseHandler {
	return &ExpenseHandler{svc: svc}
}

type ExpenseItemRequest struct {
	Name      string          `json:"name" binding:"required"`
	UnitPrice decimal.Decimal `json:"unit_price"`
	Quantity  decimal.Decimal `json:"quantity"`
	Discount  decimal.Decimal `json:"discount"`
}

type DebtRequest struct {
	PayerName  string          `json:"payer_name" binding:"required"`
	PayeeName  string          `json:"payee_name" binding:"required"`
	Amount     decimal.Decimal `json:"amount"`
	IsSpotPaid bool            `json:"is_spot_paid"`
}

type ExpenseDetailRequest struct {
	Category      string               `json:"category"`
	ExchangeRate  decimal.Decimal      `json:"exchange_rate"`
	BillingAmount decimal.Decimal      `json:"billing_amount"`
	HandlingFee   decimal.Decimal      `json:"handling_fee"`
	PaymentMethod string               `json:"payment_method"`
	Items         []ExpenseItemRequest `json:"items"`
}

type CreateExpenseRequest struct {
	Date     *time.Time           `json:"date"`
	Currency string               `json:"currency"`
	Title    string               `json:"title"`
	Note     string               `json:"note"`
	Expense  ExpenseDetailRequest `json:"expense"`
	Debts    []DebtRequest        `json:"debts"`
}

type UpdateExpenseRequest struct {
	Date     *time.Time           `json:"date"`
	Currency string               `json:"currency"`
	Title    string               `json:"title"`
	Note     string               `json:"note"`
	Expense  ExpenseDetailRequest `json:"expense"`
	Debts    []DebtRequest        `json:"debts"`
}

func toExpenseItemInputs(reqs []ExpenseItemRequest) []services.ExpenseItemInput {
	if reqs == nil {
		return nil
	}
	inputs := make([]services.ExpenseItemInput, len(reqs))
	for i, r := range reqs {
		inputs[i] = services.ExpenseItemInput{
			Name:      r.Name,
			UnitPrice: r.UnitPrice,
			Quantity:  r.Quantity,
			Discount:  r.Discount,
		}
	}
	return inputs
}

func toDebtInputs(reqs []DebtRequest) []services.DebtInput {
	if reqs == nil {
		return nil
	}
	inputs := make([]services.DebtInput, len(reqs))
	for i, r := range reqs {
		inputs[i] = services.DebtInput{
			PayerName:  r.PayerName,
			PayeeName:  r.PayeeName,
			Amount:     r.Amount,
			IsSpotPaid: r.IsSpotPaid,
		}
	}
	return inputs
}

func (h *ExpenseHandler) Create(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Space)

	var req CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txn, err := h.svc.CreateExpense(c.Request.Context(), space.ID, services.CreateExpenseInput{
		Date:     req.Date,
		Currency: req.Currency,
		Title:    req.Title,
		Note:     req.Note,
		Expense: services.ExpenseInput{
			Category:      req.Expense.Category,
			ExchangeRate:  req.Expense.ExchangeRate,
			BillingAmount: req.Expense.BillingAmount,
			HandlingFee:   req.Expense.HandlingFee,
			PaymentMethod: req.Expense.PaymentMethod,
			Items:         toExpenseItemInputs(req.Expense.Items),
		},
		Debts: toDebtInputs(req.Debts),
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, txn)
}

func (h *ExpenseHandler) Update(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Space)
	txnID := c.Param("txn_id")

	var req UpdateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txn, err := h.svc.UpdateExpense(c.Request.Context(), txnID, space.ID, services.UpdateExpenseInput{
		Date:     req.Date,
		Currency: req.Currency,
		Title:    req.Title,
		Note:     req.Note,
		Expense: services.ExpenseInput{
			Category:      req.Expense.Category,
			ExchangeRate:  req.Expense.ExchangeRate,
			BillingAmount: req.Expense.BillingAmount,
			HandlingFee:   req.Expense.HandlingFee,
			PaymentMethod: req.Expense.PaymentMethod,
			Items:         toExpenseItemInputs(req.Expense.Items),
		},
		Debts: toDebtInputs(req.Debts),
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, txn)
}
