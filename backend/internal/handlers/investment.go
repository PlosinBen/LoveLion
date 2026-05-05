package handlers

import (
	"math"
	"net/http"
	"time"

	"lovelion/internal/models"
	"lovelion/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type InvestmentHandler struct {
	db *gorm.DB
}

func NewInvestmentHandler(db *gorm.DB) *InvestmentHandler {
	return &InvestmentHandler{db: db}
}

// --- Members ---

func (h *InvestmentHandler) ListMembers(c *gin.Context) {
	var members []models.InvMember
	query := h.db.Order("sort_order ASC, created_at ASC")

	if c.Query("active") == "true" {
		query = query.Where("active = ?", true)
	}

	if err := query.Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch members"})
		return
	}
	c.JSON(http.StatusOK, members)
}

type CreateMemberRequest struct {
	Name   string     `json:"name" binding:"required,min=1,max=50"`
	UserID *uuid.UUID `json:"user_id"`
}

func (h *InvestmentHandler) CreateMember(c *gin.Context) {
	var req CreateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member := models.InvMember{
		ID:     utils.MustNewShortID(h.db, "inv_members", "id"),
		Name:   req.Name,
		UserID: req.UserID,
		Active: true,
	}

	if err := h.db.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create member"})
		return
	}
	c.JSON(http.StatusCreated, member)
}

type UpdateInvMemberRequest struct {
	Name      *string `json:"name"`
	Active    *bool   `json:"active"`
	SortOrder *int    `json:"sort_order"`
}

func (h *InvestmentHandler) UpdateMember(c *gin.Context) {
	id := c.Param("id")
	var member models.InvMember
	if err := h.db.First(&member, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}

	var req UpdateInvMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Active != nil {
		updates["active"] = *req.Active
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}

	if len(updates) > 0 {
		if err := h.db.Model(&member).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update member"})
			return
		}
	}

	h.db.First(&member, "id = ?", id)
	c.JSON(http.StatusOK, member)
}

// --- Settlements ---

func (h *InvestmentHandler) ListSettlements(c *gin.Context) {
	var settlements []models.InvSettlement
	if err := h.db.Order("year_month DESC").Find(&settlements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch settlements"})
		return
	}
	c.JSON(http.StatusOK, settlements)
}

type CreateSettlementRequest struct {
	YearMonth string `json:"year_month" binding:"required,len=7"`
}

func (h *InvestmentHandler) CreateSettlement(c *gin.Context) {
	var req CreateSettlementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate format YYYY-MM
	if _, err := time.Parse("2006-01", req.YearMonth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year_month format, expected YYYY-MM"})
		return
	}

	settlement := models.InvSettlement{
		YearMonth: req.YearMonth,
		Status:    "draft",
	}

	if err := h.db.Create(&settlement).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Settlement already exists for this month"})
		return
	}
	c.JSON(http.StatusCreated, settlement)
}

type SettlementDetail struct {
	models.InvSettlement
	FuturesStatement *models.InvFuturesStatement `json:"futures_statement"`
	StockStatement   *models.InvStockStatement   `json:"stock_statement"`
	Allocations      []AllocationPreview         `json:"allocations"`
}

type AllocationPreview struct {
	MemberID   string `json:"member_id"`
	MemberName string `json:"member_name"`
	IsOwner    bool   `json:"is_owner"`
	Weight     int    `json:"weight"`
	Amount     int    `json:"amount"`
	Deposit    int    `json:"deposit"`
	Withdrawal int    `json:"withdrawal"`
	Balance    int    `json:"balance"`
}

func (h *InvestmentHandler) GetSettlement(c *gin.Context) {
	ym := c.Param("ym")

	var settlement models.InvSettlement
	if err := h.db.First(&settlement, "year_month = ?", ym).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Settlement not found"})
		return
	}

	// Load sub-statements
	var futures *models.InvFuturesStatement
	var fStmt models.InvFuturesStatement
	if err := h.db.First(&fStmt, "year_month = ?", ym).Error; err == nil {
		futures = &fStmt
	}

	var stocks *models.InvStockStatement
	var sStmt models.InvStockStatement
	if err := h.db.Preload("Holdings").First(&sStmt, "year_month = ?", ym).Error; err == nil {
		stocks = &sStmt
	}

	// Build allocation preview
	allocations := h.buildAllocationPreview(ym, futures, stocks)

	detail := SettlementDetail{
		InvSettlement:    settlement,
		FuturesStatement: futures,
		StockStatement:   stocks,
		Allocations:      allocations,
	}

	c.JSON(http.StatusOK, detail)
}

func (h *InvestmentHandler) buildAllocationPreview(ym string, futures *models.InvFuturesStatement, stocks *models.InvStockStatement) []AllocationPreview {
	// Calculate total profit/loss
	totalPL := 0
	if futures != nil {
		totalPL += futures.ProfitLoss
	}
	if stocks != nil {
		totalPL += stocks.ProfitLoss
	}

	// Get active members
	var members []models.InvMember
	h.db.Where("active = ?", true).Order("sort_order ASC").Find(&members)

	if len(members) == 0 {
		return nil
	}

	// Get previous month's allocations for balance
	prevYM := prevYearMonth(ym)
	prevBalances := map[string]int{}
	var prevAllocs []models.InvSettlementAllocation
	h.db.Where("year_month = ?", prevYM).Find(&prevAllocs)
	for _, a := range prevAllocs {
		prevBalances[a.MemberID] = a.Balance
	}

	// Get current month member transactions
	monthStart, monthEnd := monthRange(ym)
	var txns []models.InvMemberTransaction
	h.db.Where("date >= ? AND date <= ?", monthStart, monthEnd).Find(&txns)

	memberDeposits := map[string]int{}
	memberWithdrawals := map[string]int{}
	for _, t := range txns {
		switch t.Type {
		case "deposit":
			memberDeposits[t.MemberID] += t.Amount
		case "withdrawal":
			memberWithdrawals[t.MemberID] += t.Amount
		}
	}

	// Calculate weights
	type memberWeight struct {
		member     models.InvMember
		weight     int
		deposit    int
		withdrawal int
	}

	var mws []memberWeight
	totalWeight := 0
	for _, m := range members {
		prevBal := prevBalances[m.ID]
		withdrawal := memberWithdrawals[m.ID]
		w := int(math.Floor(float64(prevBal-withdrawal) / 5000))
		if w < 1 {
			w = 1
		}
		mws = append(mws, memberWeight{
			member:     m,
			weight:     w,
			deposit:    memberDeposits[m.ID],
			withdrawal: withdrawal,
		})
		totalWeight += w
	}

	// Calculate allocations
	plPerWeight := 0
	if totalWeight > 0 {
		plPerWeight = int(math.Floor(float64(totalPL) / float64(totalWeight)))
	}

	var previews []AllocationPreview
	allocated := 0
	var ownerIdx int

	for i, mw := range mws {
		amount := 0
		if mw.member.IsOwner {
			ownerIdx = i
		} else {
			amount = plPerWeight * mw.weight
			allocated += amount
		}

		prevBal := prevBalances[mw.member.ID]
		balance := prevBal + mw.deposit - mw.withdrawal + amount

		previews = append(previews, AllocationPreview{
			MemberID:   mw.member.ID,
			MemberName: mw.member.Name,
			IsOwner:    mw.member.IsOwner,
			Weight:     mw.weight,
			Amount:     amount,
			Deposit:    mw.deposit,
			Withdrawal: mw.withdrawal,
			Balance:    balance,
		})
	}

	// Owner absorbs remainder
	if len(previews) > 0 {
		ownerAmount := totalPL - allocated
		previews[ownerIdx].Amount = ownerAmount
		prevBal := prevBalances[mws[ownerIdx].member.ID]
		previews[ownerIdx].Balance = prevBal + previews[ownerIdx].Deposit - previews[ownerIdx].Withdrawal + ownerAmount
	}

	return previews
}

func (h *InvestmentHandler) CompleteSettlement(c *gin.Context) {
	ym := c.Param("ym")

	var settlement models.InvSettlement
	if err := h.db.First(&settlement, "year_month = ?", ym).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Settlement not found"})
		return
	}

	if settlement.Status != "draft" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Settlement is not in draft status"})
		return
	}

	// Verify all required statements exist
	var futuresCount int64
	h.db.Model(&models.InvFuturesStatement{}).Where("year_month = ?", ym).Count(&futuresCount)
	var stocksCount int64
	h.db.Model(&models.InvStockStatement{}).Where("year_month = ?", ym).Count(&stocksCount)

	if futuresCount == 0 || stocksCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All statement types must be filled before completing"})
		return
	}

	// Load statements
	var futures models.InvFuturesStatement
	h.db.First(&futures, "year_month = ?", ym)
	var stocks models.InvStockStatement
	h.db.First(&stocks, "year_month = ?", ym)

	// Calculate
	totalPL := futures.ProfitLoss + stocks.ProfitLoss
	previews := h.buildAllocationPreview(ym, &futures, &stocks)

	totalWeight := 0
	for _, p := range previews {
		totalWeight += p.Weight
	}
	plPerWeight := 0
	if totalWeight > 0 {
		plPerWeight = int(math.Floor(float64(totalPL) / float64(totalWeight)))
	}

	// Transaction: update settlement + write allocations + write member_transactions
	tx := h.db.Begin()

	tx.Model(&settlement).Updates(map[string]interface{}{
		"status":                 "completed",
		"total_profit_loss":      totalPL,
		"total_weight":           totalWeight,
		"profit_loss_per_weight": plPerWeight,
	})

	// Delete old allocations and profit_loss transactions for this month
	tx.Where("year_month = ?", ym).Delete(&models.InvSettlementAllocation{})

	monthEnd := lastDayOfMonth(ym)
	tx.Where("date = ? AND type = ?", monthEnd, "profit_loss").Delete(&models.InvMemberTransaction{})

	for _, p := range previews {
		alloc := models.InvSettlementAllocation{
			YearMonth:  ym,
			MemberID:   p.MemberID,
			Weight:     p.Weight,
			Amount:     p.Amount,
			Deposit:    p.Deposit,
			Withdrawal: p.Withdrawal,
			Balance:    p.Balance,
		}
		tx.Create(&alloc)

		// Write profit_loss transaction
		if p.Amount != 0 {
			plTxn := models.InvMemberTransaction{
				ID:       uuid.New(),
				MemberID: p.MemberID,
				Date:     monthEnd,
				Type:     "profit_loss",
				Amount:   p.Amount,
			}
			tx.Create(&plTxn)
		}
	}

	// Update member net_investment
	for _, p := range previews {
		if p.Deposit != 0 || p.Withdrawal != 0 {
			tx.Model(&models.InvMember{}).Where("id = ?", p.MemberID).
				Update("net_investment", gorm.Expr("net_investment + ? - ?", p.Deposit, p.Withdrawal))
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete settlement"})
		return
	}

	h.db.First(&settlement, "year_month = ?", ym)
	c.JSON(http.StatusOK, settlement)
}

func (h *InvestmentHandler) ReopenSettlement(c *gin.Context) {
	ym := c.Param("ym")

	var settlement models.InvSettlement
	if err := h.db.First(&settlement, "year_month = ?", ym).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Settlement not found"})
		return
	}

	if settlement.Status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Settlement is not completed"})
		return
	}

	tx := h.db.Begin()
	tx.Model(&settlement).Update("status", "draft")
	tx.Where("year_month = ?", ym).Delete(&models.InvSettlementAllocation{})

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reopen settlement"})
		return
	}

	h.db.First(&settlement, "year_month = ?", ym)
	c.JSON(http.StatusOK, settlement)
}

func (h *InvestmentHandler) DeleteSettlement(c *gin.Context) {
	ym := c.Param("ym")

	var settlement models.InvSettlement
	if err := h.db.First(&settlement, "year_month = ?", ym).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Settlement not found"})
		return
	}

	if settlement.Status != "draft" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only draft settlements can be deleted"})
		return
	}

	if err := h.db.Delete(&settlement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete settlement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settlement deleted"})
}

// --- Futures Statement ---

type UpsertFuturesRequest struct {
	EndingEquity       int `json:"ending_equity"`
	FloatingProfitLoss int `json:"floating_profit_loss"`
	RealizedProfitLoss int `json:"realized_profit_loss"`
	Deposit            int `json:"deposit"`
	Withdrawal         int `json:"withdrawal"`
}

func (h *InvestmentHandler) UpsertFutures(c *gin.Context) {
	ym := c.Param("ym")

	// Verify settlement exists and is draft
	var settlement models.InvSettlement
	if err := h.db.First(&settlement, "year_month = ?", ym).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Settlement not found"})
		return
	}
	if settlement.Status != "draft" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Settlement is not in draft status"})
		return
	}

	var req UpsertFuturesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate profit_loss
	// 實質權益 = 期末權益 - 浮動損益
	realEquity := req.EndingEquity - req.FloatingProfitLoss

	// Get previous month's real equity
	prevYM := prevYearMonth(ym)
	prevRealEquity := 0
	var prevFutures models.InvFuturesStatement
	if err := h.db.First(&prevFutures, "year_month = ?", prevYM).Error; err == nil {
		prevRealEquity = prevFutures.EndingEquity - prevFutures.FloatingProfitLoss
	}

	// 權益損益 = 實質權益 - 前期實質權益 - 入金 + 出金
	equityPL := realEquity - prevRealEquity - req.Deposit + req.Withdrawal

	// 期貨損益 = min(沖銷損益, 權益損益)
	profitLoss := equityPL
	if req.RealizedProfitLoss < equityPL {
		profitLoss = req.RealizedProfitLoss
	}

	stmt := models.InvFuturesStatement{
		YearMonth:          ym,
		EndingEquity:       req.EndingEquity,
		FloatingProfitLoss: req.FloatingProfitLoss,
		RealizedProfitLoss: req.RealizedProfitLoss,
		Deposit:            req.Deposit,
		Withdrawal:         req.Withdrawal,
		ProfitLoss:         profitLoss,
	}

	result := h.db.Where("year_month = ?", ym).Assign(stmt).FirstOrCreate(&stmt)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save futures statement"})
		return
	}
	if result.RowsAffected == 0 {
		h.db.Model(&stmt).Updates(stmt)
	}

	c.JSON(http.StatusOK, stmt)
}

// --- Stock Statement ---

type UpsertStocksRequest struct {
	AccountBalance int                   `json:"account_balance"`
	Deposit        int                   `json:"deposit"`
	Withdrawal     int                   `json:"withdrawal"`
	Holdings       []StockHoldingRequest `json:"holdings"`
}

type StockHoldingRequest struct {
	Symbol       string          `json:"symbol" binding:"required"`
	Shares       int             `json:"shares"`
	ClosingPrice decimal.Decimal `json:"closing_price"`
}

func (h *InvestmentHandler) UpsertStocks(c *gin.Context) {
	ym := c.Param("ym")

	var settlement models.InvSettlement
	if err := h.db.First(&settlement, "year_month = ?", ym).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Settlement not found"})
		return
	}
	if settlement.Status != "draft" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Settlement is not in draft status"})
		return
	}

	var req UpsertStocksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate market value
	marketValue := 0
	for _, h := range req.Holdings {
		mv := h.ClosingPrice.Mul(decimal.NewFromInt(int64(h.Shares))).IntPart()
		marketValue += int(mv)
	}

	// Calculate profit_loss
	// 本期總權益 = 帳戶餘額 + 庫存現值
	currentEquity := req.AccountBalance + marketValue

	// Get previous month
	prevYM := prevYearMonth(ym)
	prevEquity := 0
	var prevStocks models.InvStockStatement
	if err := h.db.First(&prevStocks, "year_month = ?", prevYM).Error; err == nil {
		prevEquity = prevStocks.AccountBalance + prevStocks.MarketValue
	}

	// 股票損益 = 本期總權益 - 前期總權益 - 入金 + 出金
	profitLoss := currentEquity - prevEquity - req.Deposit + req.Withdrawal

	tx := h.db.Begin()

	stmt := models.InvStockStatement{
		YearMonth:      ym,
		AccountBalance: req.AccountBalance,
		MarketValue:    marketValue,
		Deposit:        req.Deposit,
		Withdrawal:     req.Withdrawal,
		ProfitLoss:     profitLoss,
	}

	// Upsert statement
	var existing models.InvStockStatement
	if err := tx.First(&existing, "year_month = ?", ym).Error; err == nil {
		tx.Model(&existing).Updates(map[string]interface{}{
			"account_balance": stmt.AccountBalance,
			"market_value":    stmt.MarketValue,
			"deposit":         stmt.Deposit,
			"withdrawal":      stmt.Withdrawal,
			"profit_loss":     stmt.ProfitLoss,
		})
	} else {
		tx.Create(&stmt)
	}

	// Replace holdings
	tx.Where("year_month = ?", ym).Delete(&models.InvStockHolding{})
	for _, hr := range req.Holdings {
		mv := hr.ClosingPrice.Mul(decimal.NewFromInt(int64(hr.Shares))).IntPart()
		holding := models.InvStockHolding{
			ID:           uuid.New(),
			YearMonth:    ym,
			Symbol:       hr.Symbol,
			Shares:       hr.Shares,
			ClosingPrice: hr.ClosingPrice,
			MarketValue:  int(mv),
		}
		tx.Create(&holding)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save stock statement"})
		return
	}

	// Reload
	h.db.Preload("Holdings").First(&stmt, "year_month = ?", ym)
	c.JSON(http.StatusOK, stmt)
}

// --- Member Transactions ---

func (h *InvestmentHandler) ListMemberTransactions(c *gin.Context) {
	query := h.db.Preload("Member").Order("date DESC")

	from := c.Query("from")
	to := c.Query("to")
	if from != "" {
		if t, err := time.Parse("2006-01-02", from); err == nil {
			query = query.Where("date >= ?", t)
		}
	}
	if to != "" {
		if t, err := time.Parse("2006-01-02", to); err == nil {
			query = query.Where("date <= ?", t)
		}
	}

	var txns []models.InvMemberTransaction
	if err := query.Find(&txns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}
	c.JSON(http.StatusOK, txns)
}

type CreateMemberTransactionRequest struct {
	MemberID string `json:"member_id" binding:"required"`
	Date     string `json:"date" binding:"required"`
	Type     string `json:"type" binding:"required,oneof=deposit withdrawal"`
	Amount   int    `json:"amount" binding:"required,gt=0"`
	Note     string `json:"note"`
}

func (h *InvestmentHandler) CreateMemberTransaction(c *gin.Context) {
	var req CreateMemberTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	// Verify member exists
	var member models.InvMember
	if err := h.db.First(&member, "id = ?", req.MemberID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Member not found"})
		return
	}

	txn := models.InvMemberTransaction{
		ID:       uuid.New(),
		MemberID: req.MemberID,
		Date:     date,
		Type:     req.Type,
		Amount:   req.Amount,
		Note:     req.Note,
	}

	if err := h.db.Create(&txn).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	h.db.Preload("Member").First(&txn, "id = ?", txn.ID)
	c.JSON(http.StatusCreated, txn)
}

type UpdateMemberTransactionRequest struct {
	MemberID *string `json:"member_id"`
	Date     *string `json:"date"`
	Type     *string `json:"type"`
	Amount   *int    `json:"amount"`
	Note     *string `json:"note"`
}

func (h *InvestmentHandler) UpdateMemberTransaction(c *gin.Context) {
	id := c.Param("id")

	var txn models.InvMemberTransaction
	if err := h.db.First(&txn, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	var req UpdateMemberTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.MemberID != nil {
		updates["member_id"] = *req.MemberID
	}
	if req.Date != nil {
		if d, err := time.Parse("2006-01-02", *req.Date); err == nil {
			updates["date"] = d
		}
	}
	if req.Type != nil {
		if *req.Type != "deposit" && *req.Type != "withdrawal" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Type must be deposit or withdrawal"})
			return
		}
		updates["type"] = *req.Type
	}
	if req.Amount != nil {
		updates["amount"] = *req.Amount
	}
	if req.Note != nil {
		updates["note"] = *req.Note
	}

	if len(updates) > 0 {
		if err := h.db.Model(&txn).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction"})
			return
		}
	}

	h.db.Preload("Member").First(&txn, "id = ?", id)
	c.JSON(http.StatusOK, txn)
}

func (h *InvestmentHandler) DeleteMemberTransaction(c *gin.Context) {
	id := c.Param("id")

	var txn models.InvMemberTransaction
	if err := h.db.First(&txn, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if err := h.db.Delete(&txn).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted"})
}

// --- Stock Trades ---

func (h *InvestmentHandler) ListStockTrades(c *gin.Context) {
	query := h.db.Order("trade_date DESC")

	from := c.Query("from")
	to := c.Query("to")
	if from != "" {
		if t, err := time.Parse("2006-01-02", from); err == nil {
			query = query.Where("trade_date >= ?", t)
		}
	}
	if to != "" {
		if t, err := time.Parse("2006-01-02", to); err == nil {
			query = query.Where("trade_date <= ?", t)
		}
	}

	var trades []models.InvStockTrade
	if err := query.Find(&trades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trades"})
		return
	}
	c.JSON(http.StatusOK, trades)
}

type CreateStockTradeRequest struct {
	TradeDate string          `json:"trade_date" binding:"required"`
	Symbol    string          `json:"symbol" binding:"required"`
	Shares    int             `json:"shares" binding:"required"`
	Price     decimal.Decimal `json:"price" binding:"required"`
	Fee       int             `json:"fee"`
	Tax       int             `json:"tax"`
}

func (h *InvestmentHandler) CreateStockTrade(c *gin.Context) {
	var req CreateStockTradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tradeDate, err := time.Parse("2006-01-02", req.TradeDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trade_date format"})
		return
	}

	trade := models.InvStockTrade{
		ID:        uuid.New(),
		TradeDate: tradeDate,
		Symbol:    req.Symbol,
		Shares:    req.Shares,
		Price:     req.Price,
		Fee:       req.Fee,
		Tax:       req.Tax,
	}

	if err := h.db.Create(&trade).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trade"})
		return
	}
	c.JSON(http.StatusCreated, trade)
}

type UpdateStockTradeRequest struct {
	TradeDate *string          `json:"trade_date"`
	Symbol    *string          `json:"symbol"`
	Shares    *int             `json:"shares"`
	Price     *decimal.Decimal `json:"price"`
	Fee       *int             `json:"fee"`
	Tax       *int             `json:"tax"`
}

func (h *InvestmentHandler) UpdateStockTrade(c *gin.Context) {
	id := c.Param("id")

	var trade models.InvStockTrade
	if err := h.db.First(&trade, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trade not found"})
		return
	}

	var req UpdateStockTradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.TradeDate != nil {
		if d, err := time.Parse("2006-01-02", *req.TradeDate); err == nil {
			updates["trade_date"] = d
		}
	}
	if req.Symbol != nil {
		updates["symbol"] = *req.Symbol
	}
	if req.Shares != nil {
		updates["shares"] = *req.Shares
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.Fee != nil {
		updates["fee"] = *req.Fee
	}
	if req.Tax != nil {
		updates["tax"] = *req.Tax
	}

	if len(updates) > 0 {
		if err := h.db.Model(&trade).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update trade"})
			return
		}
	}

	h.db.First(&trade, "id = ?", id)
	c.JSON(http.StatusOK, trade)
}

func (h *InvestmentHandler) DeleteStockTrade(c *gin.Context) {
	id := c.Param("id")

	var trade models.InvStockTrade
	if err := h.db.First(&trade, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trade not found"})
		return
	}

	if err := h.db.Delete(&trade).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete trade"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Trade deleted"})
}

// --- Allocations (read-only for non-owners) ---

func (h *InvestmentHandler) ListAllocations(c *gin.Context) {
	member := c.MustGet("invMember").(*models.InvMember)
	isOwner := member.IsOwner

	query := h.db.Preload("Member").Order("year_month DESC")

	from := c.Query("from")
	to := c.Query("to")
	if from != "" {
		query = query.Where("year_month >= ?", from)
	}
	if to != "" {
		query = query.Where("year_month <= ?", to)
	}

	if !isOwner {
		query = query.Where("member_id = ?", member.ID)
	}

	var allocs []models.InvSettlementAllocation
	if err := query.Find(&allocs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch allocations"})
		return
	}
	c.JSON(http.StatusOK, allocs)
}

// --- Helpers ---

func prevYearMonth(ym string) string {
	t, err := time.Parse("2006-01", ym)
	if err != nil {
		return ""
	}
	prev := t.AddDate(0, -1, 0)
	return prev.Format("2006-01")
}

func monthRange(ym string) (time.Time, time.Time) {
	t, _ := time.Parse("2006-01", ym)
	start := t
	end := t.AddDate(0, 1, -1)
	return start, end
}

func lastDayOfMonth(ym string) time.Time {
	t, _ := time.Parse("2006-01", ym)
	return t.AddDate(0, 1, -1)
}
