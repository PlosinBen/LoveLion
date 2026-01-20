package handlers

import (
	"net/http"
	"time"

	"lovelion/internal/models"
	"lovelion/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TripHandler struct {
	db *gorm.DB
}

func NewTripHandler(db *gorm.DB) *TripHandler {
	return &TripHandler{db: db}
}

type CreateTripRequest struct {
	Name         string     `json:"name" binding:"required,min=1,max=100"`
	Description  string     `json:"description"`
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	BaseCurrency string     `json:"base_currency"`
	Members      []string   `json:"members"` // Member names
}

type UpdateTripRequest struct {
	Name         string     `json:"name" binding:"omitempty,min=1,max=100"`
	Description  string     `json:"description"`
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	BaseCurrency string     `json:"base_currency"`
}

type AddMemberRequest struct {
	Name   string     `json:"name" binding:"required,min=1,max=50"`
	UserID *uuid.UUID `json:"user_id"`
}

// Helper to verify trip access
func (h *TripHandler) verifyTripAccess(tripID string, userID uuid.UUID) (*models.Trip, error) {
	var trip models.Trip
	if err := h.db.Where("id = ?", tripID).Preload("Members").First(&trip).Error; err != nil {
		return nil, err
	}

	// Check if user is creator or member
	if trip.CreatedBy == userID {
		return &trip, nil
	}

	for _, m := range trip.Members {
		if m.UserID != nil && *m.UserID == userID {
			return &trip, nil
		}
	}

	return nil, gorm.ErrRecordNotFound
}

// List user's trips
func (h *TripHandler) List(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var trips []models.Trip
	if err := h.db.
		Joins("LEFT JOIN trip_members ON trip_members.trip_id = trips.id").
		Where("trips.created_by = ? OR trip_members.user_id = ?", userID, userID).
		Group("trips.id").
		Preload("Members").
		Order("trips.created_at DESC").
		Find(&trips).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trips"})
		return
	}

	c.JSON(http.StatusOK, trips)
}

// Create a new trip with members
func (h *TripHandler) Create(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var req CreateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate short ID
	tripID, err := utils.NewShortID(h.db, "trips", "id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate ID"})
		return
	}

	trip := &models.Trip{
		ID:           tripID,
		Name:         req.Name,
		Description:  req.Description,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		BaseCurrency: req.BaseCurrency,
		CreatedBy:    userID,
	}

	if trip.BaseCurrency == "" {
		trip.BaseCurrency = "TWD"
	}

	// Create trip with transaction
	if err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(trip).Error; err != nil {
			return err
		}

		// Add creator as owner member
		ownerMember := &models.TripMember{
			ID:      uuid.New(),
			TripID:  tripID,
			UserID:  &userID,
			Name:    "Me", // Will be updated with actual name
			IsOwner: true,
		}
		if err := tx.Create(ownerMember).Error; err != nil {
			return err
		}

		// Add other members
		for _, name := range req.Members {
			member := &models.TripMember{
				ID:     uuid.New(),
				TripID: tripID,
				Name:   name,
			}
			if err := tx.Create(member).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trip"})
		return
	}

	// Reload with members
	h.db.Preload("Members").First(trip, "id = ?", tripID)

	c.JSON(http.StatusCreated, trip)
}

// Get a single trip
func (h *TripHandler) Get(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	tripID := c.Param("id")

	trip, err := h.verifyTripAccess(tripID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	c.JSON(http.StatusOK, trip)
}

// Update a trip
func (h *TripHandler) Update(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	tripID := c.Param("id")

	trip, err := h.verifyTripAccess(tripID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	var req UpdateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		trip.Name = req.Name
	}
	trip.Description = req.Description
	if req.StartDate != nil {
		trip.StartDate = req.StartDate
	}
	if req.EndDate != nil {
		trip.EndDate = req.EndDate
	}
	if req.BaseCurrency != "" {
		trip.BaseCurrency = req.BaseCurrency
	}

	if err := h.db.Save(trip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update trip"})
		return
	}

	c.JSON(http.StatusOK, trip)
}

// Delete a trip
func (h *TripHandler) Delete(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	tripID := c.Param("id")

	// Only creator can delete
	var trip models.Trip
	if err := h.db.Where("id = ? AND created_by = ?", tripID, userID).First(&trip).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found or not authorized"})
		return
	}

	if err := h.db.Delete(&trip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete trip"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Trip deleted"})
}

// Add member to trip
func (h *TripHandler) AddMember(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	tripID := c.Param("id")

	_, err := h.verifyTripAccess(tripID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member := &models.TripMember{
		ID:     uuid.New(),
		TripID: tripID,
		UserID: req.UserID,
		Name:   req.Name,
	}

	if err := h.db.Create(member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add member"})
		return
	}

	c.JSON(http.StatusCreated, member)
}

// Remove member from trip
func (h *TripHandler) RemoveMember(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	tripID := c.Param("id")
	memberID, err := uuid.Parse(c.Param("member_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid member ID"})
		return
	}

	_, err = h.verifyTripAccess(tripID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	// Cannot remove owner
	var member models.TripMember
	if err := h.db.Where("id = ? AND trip_id = ?", memberID, tripID).First(&member).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}

	if member.IsOwner {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot remove trip owner"})
		return
	}

	if err := h.db.Delete(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed"})
}

// List trip members
func (h *TripHandler) ListMembers(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	tripID := c.Param("id")

	_, err := h.verifyTripAccess(tripID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	var members []models.TripMember
	if err := h.db.Where("trip_id = ?", tripID).Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch members"})
		return
	}

	c.JSON(http.StatusOK, members)
}
