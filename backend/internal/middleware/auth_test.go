package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"lovelion/internal/testutil"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const testSecret = "test-jwt-secret"

func makeToken(claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := token.SignedString([]byte(testSecret))
	return s
}

func validUserToken(userID uuid.UUID) string {
	return makeToken(jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
}

func setupAuthRouter(secret string) (*gin.Engine, *uuid.UUID) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	var captured *uuid.UUID
	r.GET("/test", AuthRequired(secret), func(c *gin.Context) {
		uid := c.MustGet("userID").(uuid.UUID)
		captured = &uid
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	return r, captured
}

func TestAuthRequired_NoHeader(t *testing.T) {
	r, _ := setupAuthRouter(testSecret)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header required")
}

func TestAuthRequired_InvalidFormat(t *testing.T) {
	r, _ := setupAuthRouter(testSecret)

	tests := []struct {
		name   string
		header string
	}{
		{"no bearer prefix", "just-a-token"},
		{"wrong prefix", "Basic some-token"},
		{"too many parts", "Bearer token extra"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", tt.header)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}

func TestAuthRequired_InvalidToken(t *testing.T) {
	r, _ := setupAuthRouter(testSecret)

	tests := []struct {
		name  string
		token string
	}{
		{"garbage token", "not-a-jwt"},
		{"wrong secret", func() string {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"user_id": uuid.New().String(),
				"exp":     time.Now().Add(time.Hour).Unix(),
			})
			s, _ := token.SignedString([]byte("wrong-secret"))
			return s
		}()},
		{"expired token", makeToken(jwt.MapClaims{
			"user_id": uuid.New().String(),
			"exp":     time.Now().Add(-time.Hour).Unix(),
		})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", "Bearer "+tt.token)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}

func TestAuthRequired_MissingUserID(t *testing.T) {
	r, _ := setupAuthRouter(testSecret)

	// Token without user_id claim
	token := makeToken(jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid user ID in token")
}

func TestAuthRequired_InvalidUserIDFormat(t *testing.T) {
	r, _ := setupAuthRouter(testSecret)

	token := makeToken(jwt.MapClaims{
		"user_id": "not-a-uuid",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid user ID format")
}

func TestAuthRequired_ValidToken(t *testing.T) {
	userID := uuid.New()
	gin.SetMode(gin.TestMode)
	r := gin.New()

	var capturedID uuid.UUID
	r.GET("/test", AuthRequired(testSecret), func(c *gin.Context) {
		capturedID = c.MustGet("userID").(uuid.UUID)
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	token := validUserToken(userID)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, userID, capturedID)
}

func TestAuthRequiredWithDB_UserNotFound(t *testing.T) {
	db := testutil.TestDB(t)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/test", AuthRequiredWithDB(testSecret, db), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// Valid JWT but user doesn't exist in DB
	token := validUserToken(uuid.New())
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "User not found")
}

func TestAuthRequiredWithDB_UserExists(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/test", AuthRequiredWithDB(testSecret, db), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	token := validUserToken(user.ID)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
