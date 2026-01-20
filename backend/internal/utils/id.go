package utils

import (
	"crypto/rand"
	"math/big"

	"gorm.io/gorm"
)

const (
	// DefaultIDLength is the initial length for generated IDs
	DefaultIDLength = 5

	// MaxRetries is the number of collision retries before increasing length
	MaxRetries = 3

	// IDCharset contains allowed characters for ID generation (a-zA-Z0-9)
	IDCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// generateRandomString creates a random string of specified length using IDCharset
func generateRandomString(length int) (string, error) {
	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(IDCharset)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		result[i] = IDCharset[num.Int64()]
	}

	return string(result), nil
}

// NewShortID generates a short ID with automatic length increase on collision.
// It starts with DefaultIDLength and increases by 1 after MaxRetries collisions.
//
// Parameters:
//   - db: GORM database instance
//   - table: table name to check for collisions
//   - column: column name for the ID field
//
// Returns the generated ID or an error.
func NewShortID(db *gorm.DB, table string, column string) (string, error) {
	length := DefaultIDLength

	for attempts := 0; attempts < MaxRetries; attempts++ {
		id, err := generateRandomString(length)
		if err != nil {
			return "", err
		}

		// Check if ID already exists
		var count int64
		if err := db.Table(table).Where(column+" = ?", id).Count(&count).Error; err != nil {
			return "", err
		}

		if count == 0 {
			return id, nil
		}
	}

	// After MaxRetries collisions, try with longer length
	length++
	id, err := generateRandomString(length)
	if err != nil {
		return "", err
	}

	return id, nil
}

// MustNewShortID is like NewShortID but panics on error.
// Use this only when you're certain the operation won't fail.
func MustNewShortID(db *gorm.DB, table string, column string) string {
	id, err := NewShortID(db, table, column)
	if err != nil {
		panic(err)
	}
	return id
}
