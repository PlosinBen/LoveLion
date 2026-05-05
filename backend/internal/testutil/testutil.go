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
		setupDB.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schemaName))
	})

	var appDSN string
	if len(dsn) > 11 && strings.HasPrefix(dsn, "postgres://") {
		separator := "?"
		if strings.Contains(dsn, "?") {
			separator = "&"
		}
		appDSN = fmt.Sprintf("%s%ssearch_path=%s,public", dsn, separator, schemaName)
	} else {
		appDSN = fmt.Sprintf("%s search_path=%s,public", dsn, schemaName)
	}

	db, err := gorm.Open(postgres.Open(appDSN), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database with schema: %v", err)
	}

	// Auto-migrate Image first to establish entity_id as varchar before
	// polymorphic associations in Space/Transaction try to override it.
	err = db.AutoMigrate(&models.Image{})
	if err != nil {
		t.Fatalf("Failed to migrate Image model: %v", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Space{},
		&models.SpaceMember{},
		&models.SpaceInvite{},
		&models.Transaction{},
		&models.TransactionExpense{},
		&models.TransactionExpenseItem{},
		&models.TransactionDebt{},
		&models.ComparisonStore{},
		&models.ComparisonProduct{},
		&models.InvMember{},
		&models.InvSettlement{},
		&models.InvMemberTransaction{},
		&models.InvSettlementAllocation{},
		&models.InvFuturesStatement{},
		&models.InvStockStatement{},
		&models.InvStockHolding{},
		&models.InvStockTrade{},
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
