package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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

	// Connect to default DB to create schema
	// We use a temporary connection for setup
	setupDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to setup database: %v", err)
	}

	// Create unique schema for test isolation
	testDBCounter++
	schemaName := fmt.Sprintf("test_%d_%d", time.Now().UnixNano(), testDBCounter)

	if err := setupDB.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schemaName)).Error; err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	// Enable UUID extension (from public schema)
	setupDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\" SCHEMA public")

	// Cleanup using setupDB
	t.Cleanup(func() {
		// Re-establish connection or use setupDB if it's still valid/safe
		// Note: setupDB is a pool, so it's fine.
		setupDB.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schemaName))
	})

	// Connect to the specific schema using search_path in DSN
	// We append search_path to the DSN.
	// If DSN is a URL, this might be tricky, but assuming KV or standard libpq support.
	// For libpq KV: append " search_path=schema,public"
	// For URL: append "?search_path=schema,public" (if no query) or "&..."

	var appDSN string
	if len(dsn) > 11 && strings.HasPrefix(dsn, "postgres://") {
		// URL format - rigorous parsing would be better but simple append often works if no existing params
		// or if we trust the structure.
		// For robustness in this specific project context:
		separator := "?"
		if strings.Contains(dsn, "?") {
			separator = "&"
		}
		appDSN = fmt.Sprintf("%s%ssearch_path=%s,public", dsn, separator, schemaName)
	} else {
		// Key-Value format
		appDSN = fmt.Sprintf("%s search_path=%s,public", dsn, schemaName)
	}

	db, err := gorm.Open(postgres.Open(appDSN), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database with schema: %v", err)
	}

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
		&models.Image{}, // Added Image model to migration
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

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
