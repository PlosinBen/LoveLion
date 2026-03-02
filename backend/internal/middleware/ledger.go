package middleware

import (
	"net/http"

	"lovelion/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LedgerAccess verifies if the user is a member of the ledger.
// It preloads the Ledger and the User's LedgerMember record into the context.
func LedgerAccess(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		ledgerIDStr := c.Param("id")
		if ledgerIDStr == "" {
			// Try to get from txn_id parent if needed, but standard is :id
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ledger ID is required"})
			c.Abort()
			return
		}

		ledgerID, err := uuid.Parse(ledgerIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Ledger ID format"})
			c.Abort()
			return
		}

		// Verify membership
		var member models.LedgerMember
		if err := db.Where("ledger_id = ? AND user_id = ?", ledgerID, userID.(uuid.UUID)).First(&member).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to this ledger"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			c.Abort()
			return
		}

		// Fetch ledger to ensure it exists and provide it to handlers
		var ledger models.Ledger
		if err := db.First(&ledger, "id = ?", ledgerID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ledger not found"})
			c.Abort()
			return
		}

		// Store in context for handlers to use
		c.Set("ledger", &ledger)
		c.Set("member", &member)
		c.Next()
	}
}

// LedgerOwnerOnly restricts access to the ledger owner.
// MUST be used after LedgerAccess middleware.
func LedgerOwnerOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		memberVal, exists := c.Get("member")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ledger context missing"})
			c.Abort()
			return
		}

		member := memberVal.(*models.LedgerMember)
		if member.Role != "owner" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only the ledger owner can perform this action"})
			c.Abort()
			return
		}

		c.Next()
	}
}
