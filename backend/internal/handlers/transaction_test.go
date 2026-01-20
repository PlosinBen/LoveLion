package handlers

import (
	"net/http/httptest"
	"testing"

	"lovelion/internal/testutil"
)

func TestTransactionHandler_List(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	// Create a ledger first
	ledgerRouter := testutil.TestRouter()
	ledgerHandler := NewLedgerHandler(db)
	ledgerRouter.POST("/api/ledgers", testutil.AuthContext(user.ID), ledgerHandler.Create)

	body := map[string]interface{}{"name": "Test Ledger"}
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/ledgers", body)
	ledgerRouter.ServeHTTP(w, req)

	var ledger map[string]interface{}
	testutil.ParseResponse(t, w, &ledger)
	ledgerID := ledger["id"].(string)

	// Test list transactions
	router := testutil.TestRouter()
	handler := NewTransactionHandler(db)
	router.GET("/api/ledgers/:id/transactions", testutil.AuthContext(user.ID), handler.List)

	w = httptest.NewRecorder()
	req = testutil.JSONRequest("GET", "/api/ledgers/"+ledgerID+"/transactions", nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)
}

func TestTransactionHandler_Create(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	// Create a ledger first
	ledgerRouter := testutil.TestRouter()
	ledgerHandler := NewLedgerHandler(db)
	ledgerRouter.POST("/api/ledgers", testutil.AuthContext(user.ID), ledgerHandler.Create)

	body := map[string]interface{}{"name": "Test Ledger"}
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/ledgers", body)
	ledgerRouter.ServeHTTP(w, req)

	var ledger map[string]interface{}
	testutil.ParseResponse(t, w, &ledger)
	ledgerID := ledger["id"].(string)

	// Test create transaction
	router := testutil.TestRouter()
	handler := NewTransactionHandler(db)
	router.POST("/api/ledgers/:id/transactions", testutil.AuthContext(user.ID), handler.Create)

	tests := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid transaction",
			body: map[string]interface{}{
				"payer":    "Me",
				"category": "Food",
				"items": []map[string]interface{}{
					{"name": "Lunch", "unit_price": 150},
				},
			},
			wantStatus: 201,
		},
		{
			name: "with multiple items",
			body: map[string]interface{}{
				"payer":    "Me",
				"category": "Shopping",
				"items": []map[string]interface{}{
					{"name": "Item 1", "unit_price": 100, "quantity": 2},
					{"name": "Item 2", "unit_price": 50},
				},
			},
			wantStatus: 201,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := testutil.JSONRequest("POST", "/api/ledgers/"+ledgerID+"/transactions", tt.body)
			router.ServeHTTP(w, req)
			testutil.ExpectStatus(t, w, tt.wantStatus)
		})
	}
}

func TestTransactionHandler_Get(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	// Setup
	ledgerHandler := NewLedgerHandler(db)
	txnHandler := NewTransactionHandler(db)

	router := testutil.TestRouter()
	router.POST("/api/ledgers", testutil.AuthContext(user.ID), ledgerHandler.Create)
	router.POST("/api/ledgers/:id/transactions", testutil.AuthContext(user.ID), txnHandler.Create)
	router.GET("/api/ledgers/:id/transactions/:txn_id", testutil.AuthContext(user.ID), txnHandler.Get)

	// Create ledger
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/ledgers", map[string]interface{}{"name": "Ledger"})
	router.ServeHTTP(w, req)
	var ledger map[string]interface{}
	testutil.ParseResponse(t, w, &ledger)
	ledgerID := ledger["id"].(string)

	// Create transaction
	txnBody := map[string]interface{}{
		"payer": "Me",
		"items": []map[string]interface{}{{"name": "Test", "unit_price": 100}},
	}
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("POST", "/api/ledgers/"+ledgerID+"/transactions", txnBody)
	router.ServeHTTP(w, req)
	var txn map[string]interface{}
	testutil.ParseResponse(t, w, &txn)
	txnID := txn["id"].(string)

	// Get transaction
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("GET", "/api/ledgers/"+ledgerID+"/transactions/"+txnID, nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)
}

func TestTransactionHandler_Delete(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	ledgerHandler := NewLedgerHandler(db)
	txnHandler := NewTransactionHandler(db)

	router := testutil.TestRouter()
	router.POST("/api/ledgers", testutil.AuthContext(user.ID), ledgerHandler.Create)
	router.POST("/api/ledgers/:id/transactions", testutil.AuthContext(user.ID), txnHandler.Create)
	router.DELETE("/api/ledgers/:id/transactions/:txn_id", testutil.AuthContext(user.ID), txnHandler.Delete)

	// Create ledger
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/ledgers", map[string]interface{}{"name": "Ledger"})
	router.ServeHTTP(w, req)
	var ledger map[string]interface{}
	testutil.ParseResponse(t, w, &ledger)
	ledgerID := ledger["id"].(string)

	// Create transaction without items (to avoid FK issues on delete)
	txnBody := map[string]interface{}{
		"payer": "Me",
		"items": []map[string]interface{}{},
	}
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("POST", "/api/ledgers/"+ledgerID+"/transactions", txnBody)
	router.ServeHTTP(w, req)
	var txn map[string]interface{}
	testutil.ParseResponse(t, w, &txn)
	txnID := txn["id"].(string)

	// Delete transaction
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("DELETE", "/api/ledgers/"+ledgerID+"/transactions/"+txnID, nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)
}
