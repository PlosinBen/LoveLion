package handlers

import (
	"net/http/httptest"
	"testing"

	"lovelion/internal/middleware"
	"lovelion/internal/testutil"
)

func TestSpaceHandler_List(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	router := testutil.TestRouter()
	handler := NewSpaceHandler(db)
	router.GET("/api/spaces", testutil.AuthContext(user.ID), handler.List)

	w := httptest.NewRecorder()
	req := testutil.JSONRequest("GET", "/api/spaces", nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)
}

func TestSpaceHandler_Create(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	router := testutil.TestRouter()
	handler := NewSpaceHandler(db)
	router.POST("/api/spaces", testutil.AuthContext(user.ID), handler.Create)

	tests := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid space",
			body: map[string]interface{}{
				"name": "My Space",
			},
			wantStatus: 201,
		},
		{
			name: "with currencies and type",
			body: map[string]interface{}{
				"name":       "Japan Trip",
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
			req := testutil.JSONRequest("POST", "/api/spaces", tt.body)
			router.ServeHTTP(w, req)
			testutil.ExpectStatus(t, w, tt.wantStatus)
		})
	}
}

func TestSpaceHandler_Get(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	router := testutil.TestRouter()
	handler := NewSpaceHandler(db)
	router.POST("/api/spaces", testutil.AuthContext(user.ID), handler.Create)
	router.GET("/api/spaces/:id", testutil.AuthContext(user.ID), middleware.SpaceAccess(db), handler.Get)

	// Create a space first
	body := map[string]interface{}{"name": "Test Space"}
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/spaces", body)
	router.ServeHTTP(w, req)

	var created map[string]interface{}
	testutil.ParseResponse(t, w, &created)
	spaceID := created["id"].(string)

	// Get the space
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("GET", "/api/spaces/"+spaceID, nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)
}

func TestSpaceHandler_Update(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	router := testutil.TestRouter()
	handler := NewSpaceHandler(db)
	router.POST("/api/spaces", testutil.AuthContext(user.ID), handler.Create)
	router.PUT("/api/spaces/:id", testutil.AuthContext(user.ID), middleware.SpaceAccess(db), middleware.SpaceOwnerOnly(), handler.Update)

	// Create a space first
	createBody := map[string]interface{}{"name": "Original Name"}
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/spaces", createBody)
	router.ServeHTTP(w, req)

	var created map[string]interface{}
	testutil.ParseResponse(t, w, &created)
	spaceID := created["id"].(string)

	// Update the space
	updateBody := map[string]interface{}{"name": "Updated Name", "is_pinned": true}
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("PUT", "/api/spaces/"+spaceID, updateBody)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)

	var updated map[string]interface{}
	testutil.ParseResponse(t, w, &updated)
	if updated["name"] != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got '%v'", updated["name"])
	}
	if updated["is_pinned"] != true {
		t.Errorf("Expected is_pinned true, got %v", updated["is_pinned"])
	}
}

func TestSpaceHandler_Delete(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	router := testutil.TestRouter()
	handler := NewSpaceHandler(db)
	router.POST("/api/spaces", testutil.AuthContext(user.ID), handler.Create)
	router.DELETE("/api/spaces/:id", testutil.AuthContext(user.ID), middleware.SpaceAccess(db), middleware.SpaceOwnerOnly(), handler.Delete)

	// Create a space first
	createBody := map[string]interface{}{"name": "To Delete"}
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/spaces", createBody)
	router.ServeHTTP(w, req)

	var created map[string]interface{}
	testutil.ParseResponse(t, w, &created)
	spaceID := created["id"].(string)

	// Delete the space
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("DELETE", "/api/spaces/"+spaceID, nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)
}
