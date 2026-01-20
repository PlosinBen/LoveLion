package handlers

import (
	"net/http/httptest"
	"testing"

	"lovelion/internal/testutil"

	"github.com/google/uuid"
)

func TestLedgerHandler_List(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	router := testutil.TestRouter()
	handler := NewLedgerHandler(db)
	router.GET("/api/ledgers", testutil.AuthContext(user.ID), handler.List)

	w := httptest.NewRecorder()
	req := testutil.JSONRequest("GET", "/api/ledgers", nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)
}

func TestLedgerHandler_Create(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	router := testutil.TestRouter()
	handler := NewLedgerHandler(db)
	router.POST("/api/ledgers", testutil.AuthContext(user.ID), handler.Create)

	tests := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid ledger",
			body: map[string]interface{}{
				"name": "My Ledger",
			},
			wantStatus: 201,
		},
		{
			name: "with currencies",
			body: map[string]interface{}{
				"name":       "Trip Ledger",
				"type":       "trip",
				"currencies": []string{"TWD", "JPY"},
			},
			wantStatus: 201,
		},
		{
			name:       "missing name",
			body:       map[string]interface{}{},
			wantStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := testutil.JSONRequest("POST", "/api/ledgers", tt.body)
			router.ServeHTTP(w, req)
			testutil.ExpectStatus(t, w, tt.wantStatus)
		})
	}
}

func TestLedgerHandler_Get(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	router := testutil.TestRouter()
	handler := NewLedgerHandler(db)
	router.POST("/api/ledgers", testutil.AuthContext(user.ID), handler.Create)
	router.GET("/api/ledgers/:id", testutil.AuthContext(user.ID), handler.Get)

	// Create a ledger first
	body := map[string]interface{}{"name": "Test Ledger"}
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/ledgers", body)
	router.ServeHTTP(w, req)

	var created map[string]interface{}
	testutil.ParseResponse(t, w, &created)
	ledgerID := created["id"].(string)

	// Get the ledger
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("GET", "/api/ledgers/"+ledgerID, nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)
}

func TestLedgerHandler_Get_NotFound(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	router := testutil.TestRouter()
	handler := NewLedgerHandler(db)
	router.GET("/api/ledgers/:id", testutil.AuthContext(user.ID), handler.Get)

	w := httptest.NewRecorder()
	req := testutil.JSONRequest("GET", "/api/ledgers/"+uuid.New().String(), nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 404)
}

func TestLedgerHandler_Update(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	router := testutil.TestRouter()
	handler := NewLedgerHandler(db)
	router.POST("/api/ledgers", testutil.AuthContext(user.ID), handler.Create)
	router.PUT("/api/ledgers/:id", testutil.AuthContext(user.ID), handler.Update)

	// Create a ledger first
	createBody := map[string]interface{}{"name": "Original Name"}
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/ledgers", createBody)
	router.ServeHTTP(w, req)

	var created map[string]interface{}
	testutil.ParseResponse(t, w, &created)
	ledgerID := created["id"].(string)

	// Update the ledger
	updateBody := map[string]interface{}{"name": "Updated Name"}
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("PUT", "/api/ledgers/"+ledgerID, updateBody)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)

	var updated map[string]interface{}
	testutil.ParseResponse(t, w, &updated)
	if updated["name"] != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got '%v'", updated["name"])
	}
}

func TestLedgerHandler_Delete(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	router := testutil.TestRouter()
	handler := NewLedgerHandler(db)
	router.POST("/api/ledgers", testutil.AuthContext(user.ID), handler.Create)
	router.DELETE("/api/ledgers/:id", testutil.AuthContext(user.ID), handler.Delete)

	// Create a ledger first
	createBody := map[string]interface{}{"name": "To Delete"}
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/ledgers", createBody)
	router.ServeHTTP(w, req)

	var created map[string]interface{}
	testutil.ParseResponse(t, w, &created)
	ledgerID := created["id"].(string)

	// Delete the ledger
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("DELETE", "/api/ledgers/"+ledgerID, nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)
}
