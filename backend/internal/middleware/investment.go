package middleware

import (
	"net/http"

	"lovelion/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InvestmentAccess verifies the user is an active inv_member.
// Sets "invMember" in context. If is_owner=true, also sets "invIsOwner".
func InvestmentAccess(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		var member models.InvMember
		err := db.Where("user_id = ? AND active = ?", userID.(uuid.UUID), true).First(&member).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusForbidden, gin.H{"error": "No investment access"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			c.Abort()
			return
		}

		c.Set("invMember", &member)
		c.Set("invIsOwner", member.IsOwner)
		c.Next()
	}
}

// InvestmentOwnerOnly restricts to is_owner=true members.
// Must be used after InvestmentAccess.
func InvestmentOwnerOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		isOwner, exists := c.Get("invIsOwner")
		if !exists || !isOwner.(bool) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Owner access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
