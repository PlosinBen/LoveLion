package handlers

import (
	"net/http/httptest"
	"testing"

	"lovelion/internal/testutil"
)

func TestAuthHandler_Register(t *testing.T) {
	db := testutil.TestDB(t)
	router := testutil.TestRouter()

	handler := NewAuthHandler(db, "test-secret")
	router.POST("/api/users/register", handler.Register)

	tests := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid registration",
			body: map[string]interface{}{
				"username":     "newuser",
				"password":     "password123",
				"display_name": "New User",
			},
			wantStatus: 201,
		},
		{
			name: "missing username",
			body: map[string]interface{}{
				"password":     "password123",
				"display_name": "New User",
			},
			wantStatus: 400,
		},
		{
			name: "password too short",
			body: map[string]interface{}{
				"username":     "user2",
				"password":     "123",
				"display_name": "New User",
			},
			wantStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := testutil.JSONRequest("POST", "/api/users/register", tt.body)
			router.ServeHTTP(w, req)
			testutil.ExpectStatus(t, w, tt.wantStatus)
		})
	}
}

func TestAuthHandler_Register_DuplicateUsername(t *testing.T) {
	db := testutil.TestDB(t)
	router := testutil.TestRouter()

	handler := NewAuthHandler(db, "test-secret")
	router.POST("/api/users/register", handler.Register)

	// First registration
	body := map[string]interface{}{
		"username":     "duplicate",
		"password":     "password123",
		"display_name": "First User",
	}
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/users/register", body)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 201)

	// Second registration with same username
	w = httptest.NewRecorder()
	req = testutil.JSONRequest("POST", "/api/users/register", body)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 409)
}

func TestAuthHandler_Login(t *testing.T) {
	db := testutil.TestDB(t)
	router := testutil.TestRouter()

	handler := NewAuthHandler(db, "test-secret")
	router.POST("/api/users/register", handler.Register)
	router.POST("/api/users/login", handler.Login)

	// Register a user first
	registerBody := map[string]interface{}{
		"username":     "loginuser",
		"password":     "password123",
		"display_name": "Login User",
	}
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/users/register", registerBody)
	router.ServeHTTP(w, req)

	tests := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid login",
			body: map[string]interface{}{
				"username": "loginuser",
				"password": "password123",
			},
			wantStatus: 200,
		},
		{
			name: "wrong password",
			body: map[string]interface{}{
				"username": "loginuser",
				"password": "wrongpassword",
			},
			wantStatus: 401,
		},
		{
			name: "user not found",
			body: map[string]interface{}{
				"username": "nonexistent",
				"password": "password123",
			},
			wantStatus: 401,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := testutil.JSONRequest("POST", "/api/users/login", tt.body)
			router.ServeHTTP(w, req)
			testutil.ExpectStatus(t, w, tt.wantStatus)
		})
	}
}

func TestAuthHandler_GetMe(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	router := testutil.TestRouter()
	handler := NewAuthHandler(db, "test-secret")

	// With auth context
	router.GET("/api/users/me", testutil.AuthContext(user.ID), handler.GetMe)

	w := httptest.NewRecorder()
	req := testutil.JSONRequest("GET", "/api/users/me", nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)
}
