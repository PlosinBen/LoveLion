package handlers

import (
	"net/http/httptest"
	"testing"

	"lovelion/internal/testutil"
)

func TestComparisonHandler_ListStores(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	spaceHandler := NewSpaceHandler(db)
	comparisonHandler := NewComparisonHandler(db)

	router := testutil.TestRouter()
	router.POST("/api/spaces", testutil.AuthContext(user.ID), spaceHandler.Create)
	router.GET("/api/spaces/:id/stores", testutil.AuthContext(user.ID), comparisonHandler.ListStores)

	// Create a space
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/spaces", map[string]interface{}{"name": "Personal Space"})
	router.ServeHTTP(w, req)
	var space map[string]interface{}
	testutil.ParseResponse(t, w, &space)
	spaceID := space["id"].(string)

	// List stores
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("GET", "/api/spaces/"+spaceID+"/stores", nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)
}

func TestComparisonHandler_CreateStore(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	spaceHandler := NewSpaceHandler(db)
	comparisonHandler := NewComparisonHandler(db)

	router := testutil.TestRouter()
	router.POST("/api/spaces", testutil.AuthContext(user.ID), spaceHandler.Create)
	router.POST("/api/spaces/:id/stores", testutil.AuthContext(user.ID), comparisonHandler.CreateStore)

	// Create a space
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/spaces", map[string]interface{}{"name": "Project Space"})
	router.ServeHTTP(w, req)
	var space map[string]interface{}
	testutil.ParseResponse(t, w, &space)
	spaceID := space["id"].(string)

	// Create store
	storeBody := map[string]interface{}{
		"name":     "Don Quijote",
		"location": "Tokyo",
	}
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("POST", "/api/spaces/"+spaceID+"/stores", storeBody)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 201)
}

func TestComparisonHandler_UpdateStore(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	spaceHandler := NewSpaceHandler(db)
	comparisonHandler := NewComparisonHandler(db)

	router := testutil.TestRouter()
	router.POST("/api/spaces", testutil.AuthContext(user.ID), spaceHandler.Create)
	router.POST("/api/spaces/:id/stores", testutil.AuthContext(user.ID), comparisonHandler.CreateStore)
	router.PUT("/api/spaces/:id/stores/:store_id", testutil.AuthContext(user.ID), comparisonHandler.UpdateStore)

	// Create a space
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/spaces", map[string]interface{}{"name": "Trip Space"})
	router.ServeHTTP(w, req)
	var space map[string]interface{}
	testutil.ParseResponse(t, w, &space)
	spaceID := space["id"].(string)

	// Create store
	storeBody := map[string]interface{}{"name": "Original Store"}
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("POST", "/api/spaces/"+spaceID+"/stores", storeBody)
	router.ServeHTTP(w, req)
	var store map[string]interface{}
	testutil.ParseResponse(t, w, &store)
	storeID := store["id"].(string)

	// Update store
	updateBody := map[string]interface{}{
		"name":           "Updated Store",
		"google_map_url": "https://maps.google.com/?q=Tokyo",
		"location":       "Tokyo, Japan",
	}
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("PUT", "/api/spaces/"+spaceID+"/stores/"+storeID, updateBody)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)

	// Verify update
	var updatedStore map[string]interface{}
	testutil.ParseResponse(t, w, &updatedStore)

	if updatedStore["google_map_url"] != "https://maps.google.com/?q=Tokyo" {
		t.Errorf("Expected google_map_url to be updated, got %v", updatedStore["google_map_url"])
	}
}

func TestComparisonHandler_CreateProduct(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	spaceHandler := NewSpaceHandler(db)
	comparisonHandler := NewComparisonHandler(db)

	router := testutil.TestRouter()
	router.POST("/api/spaces", testutil.AuthContext(user.ID), spaceHandler.Create)
	router.POST("/api/spaces/:id/stores", testutil.AuthContext(user.ID), comparisonHandler.CreateStore)
	router.POST("/api/spaces/:id/stores/:store_id/products", testutil.AuthContext(user.ID), comparisonHandler.CreateProduct)

	// Create a space
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/spaces", map[string]interface{}{"name": "Trip Space"})
	router.ServeHTTP(w, req)
	var space map[string]interface{}
	testutil.ParseResponse(t, w, &space)
	spaceID := space["id"].(string)

	// Create store
	storeBody := map[string]interface{}{"name": "Store"}
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("POST", "/api/spaces/"+spaceID+"/stores", storeBody)
	router.ServeHTTP(w, req)
	var store map[string]interface{}
	testutil.ParseResponse(t, w, &store)
	storeID := store["id"].(string)

	// Create product
	productBody := map[string]interface{}{
		"name":     "Snack",
		"price":    150,
		"currency": "JPY",
	}
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("POST", "/api/spaces/"+spaceID+"/stores/"+storeID+"/products", productBody)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 201)
}

func TestComparisonHandler_DeleteStore(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	spaceHandler := NewSpaceHandler(db)
	comparisonHandler := NewComparisonHandler(db)

	router := testutil.TestRouter()
	router.POST("/api/spaces", testutil.AuthContext(user.ID), spaceHandler.Create)
	router.POST("/api/spaces/:id/stores", testutil.AuthContext(user.ID), comparisonHandler.CreateStore)
	router.DELETE("/api/spaces/:id/stores/:store_id", testutil.AuthContext(user.ID), comparisonHandler.DeleteStore)

	// Create a space
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/spaces", map[string]interface{}{"name": "Trip Space"})
	router.ServeHTTP(w, req)
	var space map[string]interface{}
	testutil.ParseResponse(t, w, &space)
	spaceID := space["id"].(string)

	// Create store
	storeBody := map[string]interface{}{"name": "Store to Delete"}
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("POST", "/api/spaces/"+spaceID+"/stores", storeBody)
	router.ServeHTTP(w, req)
	var store map[string]interface{}
	testutil.ParseResponse(t, w, &store)
	storeID := store["id"].(string)

	// Delete store
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("DELETE", "/api/spaces/"+spaceID+"/stores/"+storeID, nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)
}
