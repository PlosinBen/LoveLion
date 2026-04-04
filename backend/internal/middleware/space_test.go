package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"lovelion/internal/models"
	"lovelion/internal/testutil"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSpaceAccess_NoAuth(t *testing.T) {
	db := testutil.TestDB(t)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// No auth context — userID not set
	r.GET("/spaces/:id", SpaceAccess(db), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/spaces/"+uuid.New().String(), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authentication required")
}

func TestSpaceAccess_InvalidSpaceID(t *testing.T) {
	db := testutil.TestDB(t)
	userID := uuid.New()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/spaces/:id", testutil.AuthContext(userID), SpaceAccess(db), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/spaces/not-a-uuid", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid Space ID format")
}

func TestSpaceAccess_NotMember(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	// Create a space owned by someone else
	otherUser := testutil.CreateTestUser(t, db)
	space := models.Space{
		ID:           uuid.New(),
		UserID:       otherUser.ID,
		Name:         "Private Space",
		Type:         "personal",
		BaseCurrency: "TWD",
	}
	db.Create(&space)
	db.Create(&models.SpaceMember{
		ID:      uuid.New(),
		SpaceID: space.ID,
		UserID:  otherUser.ID,
		Role:    "owner",
	})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/spaces/:id", testutil.AuthContext(user.ID), SpaceAccess(db), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/spaces/"+space.ID.String(), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "You do not have access")
}

func TestSpaceAccess_MemberAllowed(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	space := models.Space{
		ID:           uuid.New(),
		UserID:       user.ID,
		Name:         "My Space",
		Type:         "personal",
		BaseCurrency: "TWD",
	}
	db.Create(&space)
	db.Create(&models.SpaceMember{
		ID:      uuid.New(),
		SpaceID: space.ID,
		UserID:  user.ID,
		Role:    "owner",
	})

	gin.SetMode(gin.TestMode)
	r := gin.New()

	var ctxSpace *models.Space
	var ctxMember *models.SpaceMember
	r.GET("/spaces/:id", testutil.AuthContext(user.ID), SpaceAccess(db), func(c *gin.Context) {
		ctxSpace = c.MustGet("space").(*models.Space)
		ctxMember = c.MustGet("member").(*models.SpaceMember)
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/spaces/"+space.ID.String(), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, space.ID, ctxSpace.ID)
	assert.Equal(t, "owner", ctxMember.Role)
}

func TestSpaceAccess_NonexistentSpace(t *testing.T) {
	db := testutil.TestDB(t)
	user := testutil.CreateTestUser(t, db)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/spaces/:id", testutil.AuthContext(user.ID), SpaceAccess(db), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// Valid UUID but no space or membership exists
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/spaces/"+uuid.New().String(), nil)
	r.ServeHTTP(w, req)

	// Middleware checks membership first, so returns 403 (not 404)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestSpaceOwnerOnly_NotOwner(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		// Simulate SpaceAccess having set a "member" role member
		c.Set("member", &models.SpaceMember{Role: "member"})
		c.Next()
	}, SpaceOwnerOnly(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Only the space owner")
}

func TestSpaceOwnerOnly_Owner(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		c.Set("member", &models.SpaceMember{Role: "owner"})
		c.Next()
	}, SpaceOwnerOnly(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSpaceOwnerOnly_MissingContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// No SpaceAccess middleware — member not in context
	r.GET("/test", SpaceOwnerOnly(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Space context missing")
}
