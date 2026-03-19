package middleware

import (
	"net/http"

	"lovelion/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SpaceAccess verifies if the user is a member of the space.
// It preloads the Space and the User's SpaceMember record into the context.
func SpaceAccess(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		spaceIDStr := c.Param("id")
		if spaceIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Space ID is required"})
			c.Abort()
			return
		}

		spaceID, err := uuid.Parse(spaceIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Space ID format"})
			c.Abort()
			return
		}

		// Verify membership
		var member models.SpaceMember
		if err := db.Where("space_id = ? AND user_id = ?", spaceID, userID.(uuid.UUID)).First(&member).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to this space"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			c.Abort()
			return
		}

		// Fetch space to ensure it exists and provide it to handlers
		var space models.Space
		if err := db.Preload("Images", "entity_type = ?", "space").First(&space, "id = ?", spaceID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Space not found"})
			c.Abort()
			return
		}
		space.PopulateCoverImage()

		// Store in context for handlers to use
		c.Set("space", &space)
		c.Set("member", &member)
		c.Next()
	}
}

// SpaceOwnerOnly restricts access to the space owner.
// MUST be used after SpaceAccess middleware.
func SpaceOwnerOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		memberVal, exists := c.Get("member")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Space context missing"})
			c.Abort()
			return
		}

		member := memberVal.(*models.SpaceMember)
		if member.Role != "owner" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only the space owner can perform this action"})
			c.Abort()
			return
		}

		c.Next()
	}
}
