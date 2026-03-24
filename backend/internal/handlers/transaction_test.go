package handlers

import (
	"net/http/httptest"
	"testing"

	"lovelion/internal/middleware"
	"lovelion/internal/repositories"
	"lovelion/internal/services"
	"lovelion/internal/testutil"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func newTestTransactionService(db *gorm.DB) *services.TransactionService {
	txnRepo := repositories.NewTransactionRepo(db)
	expenseRepo := repositories.NewTransactionExpenseRepo(db)
	expenseItemRepo := repositories.NewTransactionExpenseItemRepo(db)
	debtRepo := repositories.NewTransactionDebtRepo(db)
	return services.NewTransactionService(db, txnRepo, expenseRepo, expenseItemRepo, debtRepo)
}

// createTestSpace is a helper that creates a space and returns its ID
func createTestSpace(t *testing.T, db *gorm.DB, userID uuid.UUID) string {
	t.Helper()
	spaceRouter := testutil.TestRouter()
	spaceHandler := NewSpaceHandler(db)
	spaceRouter.POST("/api/spaces", testutil.AuthContext(userID), spaceHandler.Create)

	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/spaces", map[string]interface{}{
		"name":          "Test Space",
		"base_currency": "TWD",
	})
	spaceRouter.ServeHTTP(w, req)

	var space map[string]interface{}
	testutil.ParseResponse(t, w, &space)
	return space["id"].(string)
}

func TestExpenseHandler_Create(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)
	spaceID := createTestSpace(t, db, user.ID)

	svc := newTestTransactionService(db)
	handler := NewExpenseHandler(svc)

	router := testutil.TestRouter()
	router.POST("/api/spaces/:id/expenses", testutil.AuthContext(user.ID), middleware.SpaceAccess(db), handler.Create)

	tests := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid expense with items",
			body: map[string]interface{}{
				"title":    "Lunch",
				"currency": "TWD",
				"expense": map[string]interface{}{
					"category":       "Food",
					"exchange_rate":  1,
					"payment_method": "Cash",
					"items": []map[string]interface{}{
						{"name": "Ramen", "unit_price": 150, "quantity": 1},
					},
				},
			},
			wantStatus: 201,
		},
		{
			name: "expense with multiple items",
			body: map[string]interface{}{
				"title":    "Shopping",
				"currency": "TWD",
				"expense": map[string]interface{}{
					"category":       "Shopping",
					"exchange_rate":  1,
					"payment_method": "Credit Card",
					"items": []map[string]interface{}{
						{"name": "Item 1", "unit_price": 100, "quantity": 2},
						{"name": "Item 2", "unit_price": 50, "quantity": 1},
					},
				},
			},
			wantStatus: 201,
		},
		{
			name: "expense with debts",
			body: map[string]interface{}{
				"title":    "Dinner",
				"currency": "TWD",
				"expense": map[string]interface{}{
					"category":      "Food",
					"exchange_rate": 1,
					"items": []map[string]interface{}{
						{"name": "Set meal", "unit_price": 300, "quantity": 2},
					},
				},
				"debts": []map[string]interface{}{
					{"payer_name": "Bob", "payee_name": "Alice", "amount": 300},
				},
			},
			wantStatus: 201,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := testutil.JSONRequest("POST", "/api/spaces/"+spaceID+"/expenses", tt.body)
			router.ServeHTTP(w, req)
			testutil.ExpectStatus(t, w, tt.wantStatus)
		})
	}
}

func TestPaymentHandler_Create(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)
	spaceID := createTestSpace(t, db, user.ID)

	svc := newTestTransactionService(db)
	handler := NewPaymentHandler(svc)

	router := testutil.TestRouter()
	router.POST("/api/spaces/:id/payments", testutil.AuthContext(user.ID), middleware.SpaceAccess(db), handler.Create)

	tests := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid payment",
			body: map[string]interface{}{
				"title":        "Settle up",
				"total_amount": 500,
				"payer_name":   "Bob",
				"payee_name":   "Alice",
			},
			wantStatus: 201,
		},
		{
			name: "payment same payer payee",
			body: map[string]interface{}{
				"title":        "Self pay",
				"total_amount": 100,
				"payer_name":   "Alice",
				"payee_name":   "Alice",
			},
			wantStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := testutil.JSONRequest("POST", "/api/spaces/"+spaceID+"/payments", tt.body)
			router.ServeHTTP(w, req)
			testutil.ExpectStatus(t, w, tt.wantStatus)
		})
	}
}

func TestTransactionHandler_List(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)
	spaceID := createTestSpace(t, db, user.ID)

	svc := newTestTransactionService(db)
	txnHandler := NewTransactionHandler(svc)

	router := testutil.TestRouter()
	router.GET("/api/spaces/:id/transactions", testutil.AuthContext(user.ID), middleware.SpaceAccess(db), txnHandler.List)

	w := httptest.NewRecorder()
	req := testutil.JSONRequest("GET", "/api/spaces/"+spaceID+"/transactions", nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)
}

func TestTransactionHandler_GetAndDelete(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)
	spaceID := createTestSpace(t, db, user.ID)

	svc := newTestTransactionService(db)
	expenseHandler := NewExpenseHandler(svc)
	txnHandler := NewTransactionHandler(svc)

	router := testutil.TestRouter()
	router.POST("/api/spaces/:id/expenses", testutil.AuthContext(user.ID), middleware.SpaceAccess(db), expenseHandler.Create)
	router.GET("/api/spaces/:id/transactions/:txn_id", testutil.AuthContext(user.ID), middleware.SpaceAccess(db), txnHandler.Get)
	router.DELETE("/api/spaces/:id/transactions/:txn_id", testutil.AuthContext(user.ID), middleware.SpaceAccess(db), txnHandler.Delete)

	// Create an expense
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/spaces/"+spaceID+"/expenses", map[string]interface{}{
		"title":    "Test",
		"currency": "TWD",
		"expense": map[string]interface{}{
			"category":      "Food",
			"exchange_rate": 1,
			"items": []map[string]interface{}{
				{"name": "Item", "unit_price": 100, "quantity": 1},
			},
		},
	})
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 201)

	var txn map[string]interface{}
	testutil.ParseResponse(t, w, &txn)
	txnID := txn["id"].(string)

	// Get transaction
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("GET", "/api/spaces/"+spaceID+"/transactions/"+txnID, nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)

	// Delete transaction
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("DELETE", "/api/spaces/"+spaceID+"/transactions/"+txnID, nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)

	// Verify deleted
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("GET", "/api/spaces/"+spaceID+"/transactions/"+txnID, nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 404)
}
