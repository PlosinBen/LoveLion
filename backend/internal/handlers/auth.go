package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"lovelion/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db        *gorm.DB
	jwtSecret string
	jwtExpiry time.Duration
}

func NewAuthHandler(db *gorm.DB, jwtSecret string, jwtExpiry time.Duration) *AuthHandler {
	return &AuthHandler{db: db, jwtSecret: jwtSecret, jwtExpiry: jwtExpiry}
}

type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`
	Password    string `json:"password" binding:"required,min=6"`
	DisplayName string `json:"display_name" binding:"required,min=1,max=50"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

type UpdateMeRequest struct {
	DisplayName     string `json:"display_name" binding:"omitempty,min=1,max=50"`
	CurrentPassword string `json:"current_password" binding:"omitempty,min=6"`
	NewPassword     string `json:"new_password" binding:"omitempty,min=6"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username already exists
	var existingUser models.User
	if err := h.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Start a transaction to ensure both user and space are created
	err := h.db.Transaction(func(tx *gorm.DB) error {
		// 1. Create new user
		user := &models.User{
			ID:          uuid.New(),
			Username:    req.Username,
			DisplayName: req.DisplayName,
		}

		if err := user.SetPassword(req.Password); err != nil {
			return err
		}

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// 2. Create default personal space
		defaultCategories := []string{"餐飲", "交通", "購物", "娛樂", "生活", "其他"}
		categoriesJSON, _ := json.Marshal(defaultCategories)

		defaultCurrencies := []string{"TWD"}
		currenciesJSON, _ := json.Marshal(defaultCurrencies)

		space := &models.Space{
			ID:             uuid.New(),
			UserID:         user.ID,
			Name:           "我的帳本",
			Type:           "personal",
			BaseCurrency:   "TWD",
			Currencies:     datatypes.JSON(currenciesJSON),
			Categories:     datatypes.JSON(categoriesJSON),
			SplitMembers:   datatypes.JSON("[]"),
			PaymentMethods: datatypes.JSON("[]"),
		}

		if err := tx.Create(space).Error; err != nil {
			return err
		}

		// Pass the user back out of the closure
		c.Set("registeredUser", user)

		// Add user to the members table of their own space
		member := &models.SpaceMember{
			ID:      uuid.New(),
			SpaceID: space.ID,
			UserID:  user.ID,
			Role:    "owner",
		}

		if err := tx.Create(member).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user and create default space: " + err.Error()})
		return
	}

	// Get the user from context
	userValue, _ := c.Get("registeredUser")
	user := userValue.(*models.User)

	// Generate token
	token, err := h.generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, TokenResponse{Token: token, User: user})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by username
	var user models.User
	if err := h.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate token
	token, err := h.generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{Token: token, User: &user})
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	var user models.User
	if err := h.db.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var invMember models.InvMember
	invAccess := h.db.Where("user_id = ? AND active = ?", userID, true).First(&invMember).Error == nil
	invIsOwner := invAccess && invMember.IsOwner

	type meResponse struct {
		models.User
		InvAccess  bool `json:"inv_access"`
		InvIsOwner bool `json:"inv_is_owner"`
	}

	c.JSON(http.StatusOK, meResponse{
		User:       user,
		InvAccess:  invAccess,
		InvIsOwner: invIsOwner,
	})
}

func (h *AuthHandler) UpdateMe(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var req UpdateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// If updating password, verify current password
	if req.NewPassword != "" {
		if req.CurrentPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is required to set a new password"})
			return
		}
		if !user.CheckPassword(req.CurrentPassword) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid current password"})
			return
		}
		if err := user.SetPassword(req.NewPassword); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set new password"})
			return
		}
	}

	// Update display name if provided
	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) generateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(h.jwtExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}
