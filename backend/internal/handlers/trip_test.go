package handlers

import (
	"net/http/httptest"
	"testing"

	"lovelion/internal/models"
	"lovelion/internal/testutil"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func createTestTrip(t *testing.T, db interface{}, userID uuid.UUID) *models.Trip {
	trip := &models.Trip{
		ID:           "test1",
		Name:         "Test Trip",
		BaseCurrency: "TWD",
		CreatedBy:    userID,
	}

	// Use type assertion to get *gorm.DB
	gormDB := db.(interface {
		Create(value interface{}) interface{ Error() error }
	})

	if err := gormDB.Create(trip).(interface{ Error() error }).Error(); err != nil {
		t.Fatalf("Failed to create test trip: %v", err)
	}

	return trip
}

func TestTripHandler_List(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	router := testutil.TestRouter()
	handler := NewTripHandler(db)
	router.GET("/api/trips", testutil.AuthContext(user.ID), handler.List)

	w := httptest.NewRecorder()
	req := testutil.JSONRequest("GET", "/api/trips", nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)
}

func TestTripHandler_Create(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	router := testutil.TestRouter()
	handler := NewTripHandler(db)
	router.POST("/api/trips", testutil.AuthContext(user.ID), handler.Create)

	tests := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid trip",
			body: map[string]interface{}{
				"name": "Japan Trip",
			},
			wantStatus: 201,
		},
		{
			name: "with members",
			body: map[string]interface{}{
				"name":          "Korea Trip",
				"base_currency": "KRW",
				"members":       []string{"Alice", "Bob"},
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
			req := testutil.JSONRequest("POST", "/api/trips", tt.body)
			router.ServeHTTP(w, req)
			testutil.ExpectStatus(t, w, tt.wantStatus)
		})
	}
}

func TestTripHandler_Get(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	router := testutil.TestRouter()
	handler := NewTripHandler(db)
	router.POST("/api/trips", testutil.AuthContext(user.ID), handler.Create)
	router.GET("/api/trips/:id", testutil.AuthContext(user.ID), handler.Get)

	// Create a trip first
	body := map[string]interface{}{"name": "Test Trip"}
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/trips", body)
	router.ServeHTTP(w, req)

	var created map[string]interface{}
	testutil.ParseResponse(t, w, &created)
	tripID := created["id"].(string)

	// Manually update the ledger with some data
	var trip models.Trip
	if err := db.Where("id = ?", tripID).First(&trip).Error; err != nil {
		t.Fatalf("Failed to fetch trip: %v", err)
	}

	if trip.LedgerID == nil {
		t.Fatal("Trip should have a LedgerID")
	}

	// Update Ledger directly using SQL or GORM to ensure data exists
	// Note: We need to use valid JSON for jsonb/json columns.
	// In GORM/SQLite (which testutil likely uses?), it might store as text.
	// datatypes.JSON usually handles this, but let's be explicit.
	err := db.Model(&models.Ledger{}).Where("id = ?", *trip.LedgerID).Updates(map[string]interface{}{
		"categories":      datatypes.JSON(`["Food", "Travel"]`),
		"payment_methods": datatypes.JSON(`["Cash"]`),
	}).Error
	if err != nil {
		t.Fatalf("Failed to update ledger: %v", err)
	}

	// Get the trip
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("GET", "/api/trips/"+tripID, nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)

	var response map[string]interface{}
	testutil.ParseResponse(t, w, &response)

	if response["ledger"] == nil {
		t.Fatal("Response should contain ledger")
	}

	ledger := response["ledger"].(map[string]interface{})

	// Debug print
	t.Logf("Ledger response: %+v", ledger)

	categories, ok := ledger["categories"].([]interface{})
	if !ok {
		t.Errorf("categories should be a list, got %T: %v", ledger["categories"], ledger["categories"])
	} else if len(categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(categories))
	}
}

func TestTripHandler_Get_NotFound(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	router := testutil.TestRouter()
	handler := NewTripHandler(db)
	router.GET("/api/trips/:id", testutil.AuthContext(user.ID), handler.Get)

	w := httptest.NewRecorder()
	req := testutil.JSONRequest("GET", "/api/trips/nonexistent", nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 404)
}

func TestTripHandler_Update(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	router := testutil.TestRouter()
	handler := NewTripHandler(db)
	router.POST("/api/trips", testutil.AuthContext(user.ID), handler.Create)
	router.PUT("/api/trips/:id", testutil.AuthContext(user.ID), handler.Update)

	// Create a trip first
	createBody := map[string]interface{}{"name": "Original Trip"}
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/trips", createBody)
	router.ServeHTTP(w, req)

	var created map[string]interface{}
	testutil.ParseResponse(t, w, &created)
	tripID := created["id"].(string)

	// Update the trip
	updateBody := map[string]interface{}{
		"name":            "Updated Trip",
		"currencies":      []string{"USD", "EUR"},
		"categories":      []string{"Food", "Transport"},
		"payment_methods": []string{"Cash", "Credit Card"},
	}
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("PUT", "/api/trips/"+tripID, updateBody)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)

	var updatedTrip map[string]interface{}
	testutil.ParseResponse(t, w, &updatedTrip)

	if updatedTrip["ledger"] == nil {
		t.Error("Expected ledger in response")
	} else {
		ledger := updatedTrip["ledger"].(map[string]interface{})

		// Helper to check slice content
		checkSlice := func(key string, expected []string) {
			val, ok := ledger[key].([]interface{})
			if !ok {
				t.Errorf("Expected %s to be array, got %T", key, ledger[key])
				return
			}
			if len(val) != len(expected) {
				t.Errorf("Expected %d %s, got %d", len(expected), key, len(val))
			}
		}

		checkSlice("currencies", []string{"USD", "EUR"})
		checkSlice("categories", []string{"Food", "Transport"})
		checkSlice("payment_methods", []string{"Cash", "Credit Card"})
	}
}

func TestTripHandler_Delete(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	router := testutil.TestRouter()
	handler := NewTripHandler(db)
	router.POST("/api/trips", testutil.AuthContext(user.ID), handler.Create)
	router.GET("/api/trips/:id/members", testutil.AuthContext(user.ID), handler.ListMembers)
	router.DELETE("/api/trips/:id/members/:member_id", testutil.AuthContext(user.ID), handler.RemoveMember)
	router.DELETE("/api/trips/:id", testutil.AuthContext(user.ID), handler.Delete)

	// Create a trip first
	createBody := map[string]interface{}{"name": "To Delete"}
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/trips", createBody)
	router.ServeHTTP(w, req)

	var created map[string]interface{}
	testutil.ParseResponse(t, w, &created)
	tripID := created["id"].(string)

	// Get members and remove them first (owner member is auto-created)
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("GET", "/api/trips/"+tripID+"/members", nil)
	router.ServeHTTP(w, req)

	var members []map[string]interface{}
	testutil.ParseResponse(t, w, &members)

	// Remove all non-owner members (owners can't be removed, but we need to handle FK)
	// Actually, delete the trip directly by first deleting members via DB
	db.Exec("DELETE FROM trip_members WHERE trip_id = ?", tripID)

	// Now delete the trip
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("DELETE", "/api/trips/"+tripID, nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)
}

func TestTripHandler_AddMember(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	router := testutil.TestRouter()
	handler := NewTripHandler(db)
	router.POST("/api/trips", testutil.AuthContext(user.ID), handler.Create)
	router.POST("/api/trips/:id/members", testutil.AuthContext(user.ID), handler.AddMember)

	// Create a trip first
	createBody := map[string]interface{}{"name": "Trip with Members"}
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/trips", createBody)
	router.ServeHTTP(w, req)

	var created map[string]interface{}
	testutil.ParseResponse(t, w, &created)
	tripID := created["id"].(string)

	// Add member
	memberBody := map[string]interface{}{"name": "New Member"}
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("POST", "/api/trips/"+tripID+"/members", memberBody)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 201)
}
