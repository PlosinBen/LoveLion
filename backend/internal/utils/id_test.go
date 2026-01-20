package utils

import (
	"fmt"
	"os"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testCounter int

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=postgres user=postgres password=postgres dbname=lovelion port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Create unique schema for test isolation
	testCounter++
	schemaName := fmt.Sprintf("test_utils_%d_%d", time.Now().UnixNano(), testCounter)
	db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schemaName))
	db.Exec(fmt.Sprintf("SET search_path TO %s", schemaName))

	// Create test table
	db.Exec("CREATE TABLE IF NOT EXISTS test_ids (id VARCHAR(50) PRIMARY KEY)")

	t.Cleanup(func() {
		db.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schemaName))
	})

	return db
}

func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"length 5", 5},
		{"length 10", 10},
		{"length 1", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := generateRandomString(tt.length)
			if err != nil {
				t.Fatalf("generateRandomString failed: %v", err)
			}

			if len(id) != tt.length {
				t.Errorf("Expected length %d, got %d", tt.length, len(id))
			}

			// Check all characters are in charset
			for _, c := range id {
				if !contains(IDCharset, byte(c)) {
					t.Errorf("Character %c not in charset", c)
				}
			}
		})
	}
}

func TestGenerateRandomString_Uniqueness(t *testing.T) {
	ids := make(map[string]bool)
	count := 1000

	for i := 0; i < count; i++ {
		id, err := generateRandomString(DefaultIDLength)
		if err != nil {
			t.Fatalf("generateRandomString failed: %v", err)
		}
		if ids[id] {
			t.Errorf("Duplicate ID generated: %s", id)
		}
		ids[id] = true
	}
}

func TestNewShortID(t *testing.T) {
	db := setupTestDB(t)

	id, err := NewShortID(db, "test_ids", "id")
	if err != nil {
		t.Fatalf("NewShortID failed: %v", err)
	}

	if len(id) != DefaultIDLength {
		t.Errorf("Expected length %d, got %d", DefaultIDLength, len(id))
	}
}

func TestNewShortID_NoCollision(t *testing.T) {
	db := setupTestDB(t)

	// Generate multiple IDs, they should all be unique
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id, err := NewShortID(db, "test_ids", "id")
		if err != nil {
			t.Fatalf("NewShortID failed: %v", err)
		}
		if ids[id] {
			t.Errorf("Duplicate ID: %s", id)
		}
		ids[id] = true

		// Insert ID into table
		db.Exec("INSERT INTO test_ids (id) VALUES (?)", id)
	}
}

func TestMustNewShortID(t *testing.T) {
	db := setupTestDB(t)

	// Should not panic
	id := MustNewShortID(db, "test_ids", "id")
	if len(id) < DefaultIDLength {
		t.Errorf("Expected length at least %d, got %d", DefaultIDLength, len(id))
	}
}

// Helper function
func contains(s string, c byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return true
		}
	}
	return false
}
