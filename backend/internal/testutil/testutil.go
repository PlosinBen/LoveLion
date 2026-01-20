package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"lovelion/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDBCounter int

// TestDB creates a PostgreSQL test database with isolated schema
func TestDB(t *testing.T) *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=postgres user=postgres password=postgres dbname=lovelion port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Create unique schema for test isolation
	testDBCounter++
	schemaName := fmt.Sprintf("test_%d_%d", time.Now().UnixNano(), testDBCounter)
	db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schemaName))
	db.Exec(fmt.Sprintf("SET search_path TO %s, public", schemaName))

	// Enable UUID extension (from public schema)
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\" SCHEMA public")

	// Auto-migrate all models
	err = db.AutoMigrate(
		&models.User{},
		&models.Ledger{},
		&models.Transaction{},
		&models.TransactionItem{},
		&models.Trip{},
		&models.TripMember{},
		&models.ComparisonStore{},
		&models.ComparisonProduct{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Cleanup after test
	t.Cleanup(func() {
		db.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schemaName))
	})

	return db
}

// TestRouter creates a Gin router in test mode
func TestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// CreateTestUser creates a test user and returns it
func CreateTestUser(t *testing.T, db *gorm.DB) *models.User {
	user := &models.User{
		ID:          uuid.New(),
		Username:    fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		DisplayName: "Test User",
	}
	if err := user.SetPassword("password123"); err != nil {
		t.Fatalf("Failed to set password: %v", err)
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

// AuthContext sets userID in gin context for testing protected routes
func AuthContext(userID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	}
}

// JSONRequest creates a test HTTP request with JSON body
func JSONRequest(method, path string, body interface{}) *http.Request {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// ParseResponse parses JSON response body into target struct
func ParseResponse(t *testing.T, w *httptest.ResponseRecorder, target interface{}) {
	if err := json.Unmarshal(w.Body.Bytes(), target); err != nil {
		t.Fatalf("Failed to parse response: %v, body: %s", err, w.Body.String())
	}
}

// ExpectStatus asserts response has expected status code
func ExpectStatus(t *testing.T, w *httptest.ResponseRecorder, expected int) {
	if w.Code != expected {
		t.Errorf("Expected status %d, got %d. Body: %s", expected, w.Code, w.Body.String())
	}
}
