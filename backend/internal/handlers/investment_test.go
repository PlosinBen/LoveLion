package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"lovelion/internal/middleware"
	"lovelion/internal/models"
	"lovelion/internal/testutil"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupInvestmentTest(t *testing.T) (*gin.Engine, *InvestmentHandler, *models.User, *models.InvMember) {
	db := testutil.TestDB(t)
	handler := NewInvestmentHandler(db)

	user := testutil.CreateTestUser(t, db)
	member := &models.InvMember{
		ID:      "test-owner-001",
		Name:    "Owner",
		UserID:  &user.ID,
		IsOwner: true,
		Active:  true,
	}
	db.Create(member)

	r := testutil.TestRouter()
	inv := r.Group("/api/investments")
	inv.Use(testutil.AuthContext(user.ID), middleware.InvestmentAccess(db))
	{
		inv.GET("/allocations", handler.ListAllocations)

		owner := inv.Group("")
		owner.Use(middleware.InvestmentOwnerOnly())
		{
			owner.GET("/members", handler.ListMembers)
			owner.POST("/members", handler.CreateMember)
			owner.PUT("/members/:id", handler.UpdateMember)

			owner.GET("/settlements", handler.ListSettlements)
			owner.POST("/settlements", handler.CreateSettlement)
			owner.GET("/settlements/:ym", handler.GetSettlement)
			owner.PUT("/settlements/:ym/complete", handler.CompleteSettlement)
			owner.PUT("/settlements/:ym/reopen", handler.ReopenSettlement)
			owner.DELETE("/settlements/:ym", handler.DeleteSettlement)

			owner.PUT("/settlements/:ym/futures", handler.UpsertFutures)
			owner.PUT("/settlements/:ym/stocks", handler.UpsertStocks)

			owner.GET("/members/transactions", handler.ListMemberTransactions)
			owner.POST("/members/transactions", handler.CreateMemberTransaction)
			owner.PUT("/members/transactions/:id", handler.UpdateMemberTransaction)
			owner.DELETE("/members/transactions/:id", handler.DeleteMemberTransaction)

			owner.GET("/stocks/trades", handler.ListStockTrades)
			owner.POST("/stocks/trades", handler.CreateStockTrade)
			owner.PUT("/stocks/trades/:id", handler.UpdateStockTrade)
			owner.DELETE("/stocks/trades/:id", handler.DeleteStockTrade)
		}
	}

	return r, handler, user, member
}

func TestInvestmentAccess_NoMember(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	r := testutil.TestRouter()
	r.GET("/test", testutil.AuthContext(user.ID), middleware.InvestmentAccess(db), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestInvestmentAccess_InactiveMember(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	member := &models.InvMember{
		ID:     "inactive-001",
		Name:   "Inactive",
		UserID: &user.ID,
		Active: true,
	}
	db.Create(member)
	db.Model(member).Update("active", false)

	r := testutil.TestRouter()
	r.GET("/test", testutil.AuthContext(user.ID), middleware.InvestmentAccess(db), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestInvestmentOwnerOnly_NonOwner(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	db.Create(&models.InvMember{
		ID:      "member-001",
		Name:    "Member",
		UserID:  &user.ID,
		IsOwner: false,
		Active:  true,
	})

	handler := NewInvestmentHandler(db)
	r := testutil.TestRouter()
	inv := r.Group("/api/investments")
	inv.Use(testutil.AuthContext(user.ID), middleware.InvestmentAccess(db))
	{
		owner := inv.Group("")
		owner.Use(middleware.InvestmentOwnerOnly())
		owner.GET("/members", handler.ListMembers)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/investments/members", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// --- Member CRUD ---

func TestCreateMember(t *testing.T) {
	r, _, _, _ := setupInvestmentTest(t)

	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/investments/members", map[string]interface{}{
		"name": "New Member",
	})
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var member models.InvMember
	json.Unmarshal(w.Body.Bytes(), &member)
	assert.Equal(t, "New Member", member.Name)
	assert.True(t, member.Active)
}

func TestUpdateMember(t *testing.T) {
	r, _, _, owner := setupInvestmentTest(t)

	w := httptest.NewRecorder()
	active := false
	req := testutil.JSONRequest("PUT", "/api/investments/members/"+owner.ID, map[string]interface{}{
		"name":   "Renamed",
		"active": active,
	})
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var member models.InvMember
	json.Unmarshal(w.Body.Bytes(), &member)
	assert.Equal(t, "Renamed", member.Name)
	assert.False(t, member.Active)
}

// --- Settlements ---

func TestCreateSettlement(t *testing.T) {
	r, _, _, _ := setupInvestmentTest(t)

	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/investments/settlements", map[string]interface{}{
		"year_month": "2026-05",
	})
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var s models.InvSettlement
	json.Unmarshal(w.Body.Bytes(), &s)
	assert.Equal(t, "2026-05", s.YearMonth)
	assert.Equal(t, "draft", s.Status)
}

func TestCreateSettlement_Duplicate(t *testing.T) {
	r, _, _, _ := setupInvestmentTest(t)

	req := testutil.JSONRequest("POST", "/api/investments/settlements", map[string]interface{}{
		"year_month": "2026-05",
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Duplicate
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("POST", "/api/investments/settlements", map[string]interface{}{
		"year_month": "2026-05",
	})
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestDeleteSettlement_OnlyDraft(t *testing.T) {
	r, handler, _, _ := setupInvestmentTest(t)

	// Create and complete a settlement
	handler.db.Create(&models.InvSettlement{YearMonth: "2026-04", Status: "completed"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/investments/settlements/2026-04", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// --- Futures ---

func TestUpsertFutures(t *testing.T) {
	r, handler, _, _ := setupInvestmentTest(t)

	handler.db.Create(&models.InvSettlement{YearMonth: "2026-05", Status: "draft"})

	w := httptest.NewRecorder()
	req := testutil.JSONRequest("PUT", "/api/investments/settlements/2026-05/futures", map[string]interface{}{
		"ending_equity":        120000,
		"floating_profit_loss": 5000,
		"realized_profit_loss": 15000,
		"deposit":              0,
		"withdrawal":           0,
	})
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var stmt models.InvFuturesStatement
	json.Unmarshal(w.Body.Bytes(), &stmt)
	assert.Equal(t, "2026-05", stmt.YearMonth)
	assert.NotZero(t, stmt.ProfitLoss)
}

func TestUpsertFutures_NotDraft(t *testing.T) {
	r, handler, _, _ := setupInvestmentTest(t)

	handler.db.Create(&models.InvSettlement{YearMonth: "2026-05", Status: "completed"})

	w := httptest.NewRecorder()
	req := testutil.JSONRequest("PUT", "/api/investments/settlements/2026-05/futures", map[string]interface{}{
		"ending_equity": 100000,
	})
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// --- Stocks ---

func TestUpsertStocks(t *testing.T) {
	r, handler, _, _ := setupInvestmentTest(t)

	handler.db.Create(&models.InvSettlement{YearMonth: "2026-05", Status: "draft"})

	w := httptest.NewRecorder()
	req := testutil.JSONRequest("PUT", "/api/investments/settlements/2026-05/stocks", map[string]interface{}{
		"account_balance": 80000,
		"deposit":         50000,
		"withdrawal":      0,
		"holdings": []map[string]interface{}{
			{"symbol": "2330", "shares": 1000, "closing_price": 650.5},
			{"symbol": "0056", "shares": 2000, "closing_price": 38.5},
		},
	})
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var stmt models.InvStockStatement
	json.Unmarshal(w.Body.Bytes(), &stmt)
	assert.Equal(t, "2026-05", stmt.YearMonth)
	assert.Equal(t, 727500, stmt.MarketValue) // 650.5*1000 + 38.5*2000 = 650500 + 77000
}

// --- Complete/Reopen ---

func TestCompleteSettlement(t *testing.T) {
	r, handler, _, _ := setupInvestmentTest(t)

	handler.db.Create(&models.InvSettlement{YearMonth: "2026-05", Status: "draft"})
	handler.db.Create(&models.InvFuturesStatement{
		YearMonth: "2026-05", ProfitLoss: 10000,
	})
	handler.db.Create(&models.InvStockStatement{
		YearMonth: "2026-05", ProfitLoss: 5000,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/investments/settlements/2026-05/complete", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var s models.InvSettlement
	json.Unmarshal(w.Body.Bytes(), &s)
	assert.Equal(t, "completed", s.Status)
	assert.Equal(t, 15000, s.TotalProfitLoss)

	// Verify allocations created
	var allocs []models.InvSettlementAllocation
	handler.db.Where("year_month = ?", "2026-05").Find(&allocs)
	assert.NotEmpty(t, allocs)
}

func TestCompleteSettlement_MissingStatements(t *testing.T) {
	r, handler, _, _ := setupInvestmentTest(t)

	handler.db.Create(&models.InvSettlement{YearMonth: "2026-05", Status: "draft"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/investments/settlements/2026-05/complete", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "All statement types must be filled")
}

func TestReopenSettlement(t *testing.T) {
	r, handler, _, _ := setupInvestmentTest(t)

	handler.db.Create(&models.InvSettlement{YearMonth: "2026-05", Status: "completed"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/investments/settlements/2026-05/reopen", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var s models.InvSettlement
	json.Unmarshal(w.Body.Bytes(), &s)
	assert.Equal(t, "draft", s.Status)
}

// --- Member Transactions ---

func TestCreateMemberTransaction(t *testing.T) {
	r, _, _, owner := setupInvestmentTest(t)

	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/investments/members/transactions", map[string]interface{}{
		"member_id": owner.ID,
		"date":      "2026-05-15",
		"type":      "deposit",
		"amount":    10000,
		"note":      "加碼",
	})
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var txn models.InvMemberTransaction
	json.Unmarshal(w.Body.Bytes(), &txn)
	assert.Equal(t, owner.ID, txn.MemberID)
	assert.Equal(t, "deposit", txn.Type)
	assert.Equal(t, 10000, txn.Amount)
}

func TestCreateMemberTransaction_InvalidType(t *testing.T) {
	r, _, _, owner := setupInvestmentTest(t)

	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/investments/members/transactions", map[string]interface{}{
		"member_id": owner.ID,
		"date":      "2026-05-15",
		"type":      "invalid",
		"amount":    10000,
	})
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// --- Stock Trades ---

func TestCreateStockTrade(t *testing.T) {
	r, _, _, _ := setupInvestmentTest(t)

	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/investments/stocks/trades", map[string]interface{}{
		"trade_date": "2026-05-02",
		"symbol":     "2330",
		"shares":     1000,
		"price":      650.5,
		"fee":        928,
		"tax":        0,
	})
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var trade models.InvStockTrade
	json.Unmarshal(w.Body.Bytes(), &trade)
	assert.Equal(t, "2330", trade.Symbol)
	assert.Equal(t, 1000, trade.Shares)
}

func TestDeleteStockTrade(t *testing.T) {
	r, handler, _, _ := setupInvestmentTest(t)

	trade := models.InvStockTrade{
		ID:     uuid.New(),
		Symbol: "2330",
		Shares: 1000,
	}
	handler.db.Create(&trade)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/investments/stocks/trades/"+trade.ID.String(), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// --- Allocations (non-owner access) ---

func TestListAllocations_NonOwner(t *testing.T) {
	db := testutil.TestDB(t)
	handler := NewInvestmentHandler(db)

	owner := testutil.CreateTestUser(t, db)
	nonOwner := testutil.CreateTestUser(t, db)

	ownerMember := &models.InvMember{ID: "owner-001", Name: "Owner", UserID: &owner.ID, IsOwner: true, Active: true}
	nonOwnerMember := &models.InvMember{ID: "member-002", Name: "Member", UserID: &nonOwner.ID, IsOwner: false, Active: true}
	db.Create(ownerMember)
	db.Create(nonOwnerMember)

	// Create allocations for both
	db.Create(&models.InvSettlement{YearMonth: "2026-05", Status: "completed"})
	db.Create(&models.InvSettlementAllocation{YearMonth: "2026-05", MemberID: "owner-001", Amount: 8000, Balance: 58000})
	db.Create(&models.InvSettlementAllocation{YearMonth: "2026-05", MemberID: "member-002", Amount: 2000, Balance: 17000})

	r := testutil.TestRouter()
	inv := r.Group("/api/investments")
	inv.Use(testutil.AuthContext(nonOwner.ID), middleware.InvestmentAccess(db))
	inv.GET("/allocations", handler.ListAllocations)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/investments/allocations", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var allocs []models.InvSettlementAllocation
	json.Unmarshal(w.Body.Bytes(), &allocs)
	require.Len(t, allocs, 1)
	assert.Equal(t, "member-002", allocs[0].MemberID)
}

// --- Allocation Preview (profit distribution) ---

func TestAllocationPreview_OwnerAbsorbsRemainder(t *testing.T) {
	db := testutil.TestDB(t)
	handler := NewInvestmentHandler(db)

	user := testutil.CreateTestUser(t, db)
	ownerMember := &models.InvMember{ID: "owner-001", Name: "Owner", UserID: &user.ID, IsOwner: true, Active: true}
	memberA := &models.InvMember{ID: "member-a", Name: "A", Active: true}
	memberB := &models.InvMember{ID: "member-b", Name: "B", Active: true}
	db.Create(ownerMember)
	db.Create(memberA)
	db.Create(memberB)

	// Previous month allocations (balance determines weight)
	db.Create(&models.InvSettlement{YearMonth: "2026-04", Status: "completed"})
	db.Create(&models.InvSettlementAllocation{YearMonth: "2026-04", MemberID: "owner-001", Balance: 50000})
	db.Create(&models.InvSettlementAllocation{YearMonth: "2026-04", MemberID: "member-a", Balance: 15000})
	db.Create(&models.InvSettlementAllocation{YearMonth: "2026-04", MemberID: "member-b", Balance: 5000})

	// Current month
	db.Create(&models.InvSettlement{YearMonth: "2026-05", Status: "draft"})
	db.Create(&models.InvFuturesStatement{YearMonth: "2026-05", ProfitLoss: 15000})
	db.Create(&models.InvStockStatement{YearMonth: "2026-05", ProfitLoss: -3197})

	// Get settlement detail which includes preview
	r := testutil.TestRouter()
	inv := r.Group("/api/investments")
	inv.Use(testutil.AuthContext(user.ID), middleware.InvestmentAccess(db), middleware.InvestmentOwnerOnly())
	inv.GET("/settlements/:ym", handler.GetSettlement)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/investments/settlements/2026-05", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var detail SettlementDetail
	json.Unmarshal(w.Body.Bytes(), &detail)

	// Total PL = 15000 - 3197 = 11803
	// Weights: owner=floor(50000/5000)=10, A=floor(15000/5000)=3, B=floor(5000/5000)=1
	// PL per weight = floor(11803/14) = 843
	// A = 843*3 = 2529, B = 843*1 = 843
	// Owner = 11803 - 2529 - 843 = 8431
	totalPL := 11803
	require.Len(t, detail.Allocations, 3)

	var ownerAlloc, aAlloc, bAlloc AllocationPreview
	for _, a := range detail.Allocations {
		switch a.MemberID {
		case "owner-001":
			ownerAlloc = a
		case "member-a":
			aAlloc = a
		case "member-b":
			bAlloc = a
		}
	}

	assert.Equal(t, 10, ownerAlloc.Weight)
	assert.Equal(t, 3, aAlloc.Weight)
	assert.Equal(t, 1, bAlloc.Weight)

	assert.Equal(t, 2529, aAlloc.Amount)
	assert.Equal(t, 843, bAlloc.Amount)
	assert.Equal(t, totalPL-2529-843, ownerAlloc.Amount) // 8431
}

func TestCompleteSettlement_WritesProfitLossTransactions(t *testing.T) {
	r, handler, _, _ := setupInvestmentTest(t)

	handler.db.Create(&models.InvSettlement{YearMonth: "2026-05", Status: "draft"})
	handler.db.Create(&models.InvFuturesStatement{YearMonth: "2026-05", ProfitLoss: 8000})
	handler.db.Create(&models.InvStockStatement{YearMonth: "2026-05", ProfitLoss: 2000})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/investments/settlements/2026-05/complete", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify profit_loss member_transaction was created for the owner
	var txns []models.InvMemberTransaction
	handler.db.Where("type = ? AND member_id = ?", "profit_loss", "test-owner-001").Find(&txns)
	require.Len(t, txns, 1)
	assert.Equal(t, 10000, txns[0].Amount)
	// Date should be last day of month
	assert.Equal(t, "2026-05-31", txns[0].Date.Format("2006-01-02"))
}

func TestListAllocations_RangeFilter(t *testing.T) {
	db := testutil.TestDB(t)
	handler := NewInvestmentHandler(db)

	user := testutil.CreateTestUser(t, db)
	member := &models.InvMember{ID: "owner-001", Name: "Owner", UserID: &user.ID, IsOwner: true, Active: true}
	db.Create(member)

	db.Create(&models.InvSettlement{YearMonth: "2026-03", Status: "completed"})
	db.Create(&models.InvSettlement{YearMonth: "2026-04", Status: "completed"})
	db.Create(&models.InvSettlement{YearMonth: "2026-05", Status: "completed"})
	db.Create(&models.InvSettlementAllocation{YearMonth: "2026-03", MemberID: "owner-001", Amount: 1000, Balance: 51000})
	db.Create(&models.InvSettlementAllocation{YearMonth: "2026-04", MemberID: "owner-001", Amount: 2000, Balance: 53000})
	db.Create(&models.InvSettlementAllocation{YearMonth: "2026-05", MemberID: "owner-001", Amount: 3000, Balance: 56000})

	r := testutil.TestRouter()
	inv := r.Group("/api/investments")
	inv.Use(testutil.AuthContext(user.ID), middleware.InvestmentAccess(db))
	inv.GET("/allocations", handler.ListAllocations)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/investments/allocations?from=2026-04&to=2026-05", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var allocs []models.InvSettlementAllocation
	json.Unmarshal(w.Body.Bytes(), &allocs)
	require.Len(t, allocs, 2)
	assert.Equal(t, "2026-05", allocs[0].YearMonth)
	assert.Equal(t, "2026-04", allocs[1].YearMonth)
}

func TestListStockTrades_DateFilter(t *testing.T) {
	r, handler, _, _ := setupInvestmentTest(t)

	handler.db.Create(&models.InvStockTrade{
		TradeDate: time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC),
		Symbol:    "2330", Shares: 1000, Price: decimal.NewFromFloat(600),
	})
	handler.db.Create(&models.InvStockTrade{
		TradeDate: time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC),
		Symbol:    "2317", Shares: -500, Price: decimal.NewFromFloat(120),
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/investments/stocks/trades?from=2026-05-01&to=2026-05-31", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var trades []models.InvStockTrade
	json.Unmarshal(w.Body.Bytes(), &trades)
	require.Len(t, trades, 1)
	assert.Equal(t, "2317", trades[0].Symbol)
}

func TestReopenSettlement_ClearsAllocations(t *testing.T) {
	r, handler, _, _ := setupInvestmentTest(t)

	handler.db.Create(&models.InvSettlement{YearMonth: "2026-05", Status: "completed", TotalProfitLoss: 10000})
	handler.db.Create(&models.InvSettlementAllocation{YearMonth: "2026-05", MemberID: "test-owner-001", Amount: 10000, Balance: 60000})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/investments/settlements/2026-05/reopen", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var allocs []models.InvSettlementAllocation
	handler.db.Where("year_month = ?", "2026-05").Find(&allocs)
	assert.Empty(t, allocs)
}
