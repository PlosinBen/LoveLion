package handlers

import (
	"net/http/httptest"
	"testing"

	"lovelion/internal/middleware"
	"lovelion/internal/repositories"
	"lovelion/internal/services"
	"lovelion/internal/testutil"

	"gorm.io/gorm"
)

func newTestTransactionHandler(db *gorm.DB) *TransactionHandler {
	txnRepo := repositories.NewTransactionRepo(db)
	itemRepo := repositories.NewTransactionItemRepo(db)
	splitRepo := repositories.NewTransactionSplitRepo(db)
	svc := services.NewTransactionService(db, txnRepo, itemRepo, splitRepo)
	return NewTransactionHandler(svc)
}

func TestTransactionHandler_List(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	// Create a space first
	spaceRouter := testutil.TestRouter()
	spaceHandler := NewSpaceHandler(db)
	spaceRouter.POST("/api/spaces", testutil.AuthContext(user.ID), spaceHandler.Create)

	body := map[string]interface{}{"name": "Test Space"}
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/spaces", body)
	spaceRouter.ServeHTTP(w, req)

	var space map[string]interface{}
	testutil.ParseResponse(t, w, &space)
	spaceID := space["id"].(string)

	// Test list transactions
	router := testutil.TestRouter()
	handler := newTestTransactionHandler(db)
	router.GET("/api/spaces/:id/transactions", testutil.AuthContext(user.ID), middleware.SpaceAccess(db), handler.List)

	w = httptest.NewRecorder()
	req = testutil.JSONRequest("GET", "/api/spaces/"+spaceID+"/transactions", nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)
}

func TestTransactionHandler_Create(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	// Create a space first
	spaceRouter := testutil.TestRouter()
	spaceHandler := NewSpaceHandler(db)
	spaceRouter.POST("/api/spaces", testutil.AuthContext(user.ID), spaceHandler.Create)

	body := map[string]interface{}{"name": "Test Space"}
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/spaces", body)
	spaceRouter.ServeHTTP(w, req)

	var space map[string]interface{}
	testutil.ParseResponse(t, w, &space)
	spaceID := space["id"].(string)

	// Test create transaction
	router := testutil.TestRouter()
	handler := newTestTransactionHandler(db)
	router.POST("/api/spaces/:id/transactions", testutil.AuthContext(user.ID), middleware.SpaceAccess(db), handler.Create)

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
			req := testutil.JSONRequest("POST", "/api/spaces/"+spaceID+"/transactions", tt.body)
			router.ServeHTTP(w, req)
			testutil.ExpectStatus(t, w, tt.wantStatus)
		})
	}
}

func TestTransactionHandler_Get(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	// Setup
	spaceHandler := NewSpaceHandler(db)
	txnHandler := newTestTransactionHandler(db)

	router := testutil.TestRouter()
	router.POST("/api/spaces", testutil.AuthContext(user.ID), spaceHandler.Create)
	router.POST("/api/spaces/:id/transactions", testutil.AuthContext(user.ID), middleware.SpaceAccess(db), txnHandler.Create)
	router.GET("/api/spaces/:id/transactions/:txn_id", testutil.AuthContext(user.ID), middleware.SpaceAccess(db), txnHandler.Get)

	// Create space
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/spaces", map[string]interface{}{"name": "Space"})
	router.ServeHTTP(w, req)
	var space map[string]interface{}
	testutil.ParseResponse(t, w, &space)
	spaceID := space["id"].(string)

	// Create transaction
	txnBody := map[string]interface{}{
		"payer": "Me",
		"items": []map[string]interface{}{{"name": "Test", "unit_price": 100}},
	}
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("POST", "/api/spaces/"+spaceID+"/transactions", txnBody)
	router.ServeHTTP(w, req)
	var txn map[string]interface{}
	testutil.ParseResponse(t, w, &txn)
	txnID := txn["id"].(string)

	// Get transaction
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("GET", "/api/spaces/"+spaceID+"/transactions/"+txnID, nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)
}

func TestTransactionHandler_Delete(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	spaceHandler := NewSpaceHandler(db)
	txnHandler := newTestTransactionHandler(db)

	router := testutil.TestRouter()
	router.POST("/api/spaces", testutil.AuthContext(user.ID), spaceHandler.Create)
	router.POST("/api/spaces/:id/transactions", testutil.AuthContext(user.ID), middleware.SpaceAccess(db), txnHandler.Create)
	router.DELETE("/api/spaces/:id/transactions/:txn_id", testutil.AuthContext(user.ID), middleware.SpaceAccess(db), txnHandler.Delete)

	// Create space
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/spaces", map[string]interface{}{"name": "Space"})
	router.ServeHTTP(w, req)
	var space map[string]interface{}
	testutil.ParseResponse(t, w, &space)
	spaceID := space["id"].(string)

	// Create transaction
	txnBody := map[string]interface{}{
		"payer": "Me",
		"items": []map[string]interface{}{},
	}
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("POST", "/api/spaces/"+spaceID+"/transactions", txnBody)
	router.ServeHTTP(w, req)
	var txn map[string]interface{}
	testutil.ParseResponse(t, w, &txn)
	txnID := txn["id"].(string)

	// Delete transaction
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("DELETE", "/api/spaces/"+spaceID+"/transactions/"+txnID, nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)
}
