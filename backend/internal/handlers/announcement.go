package handlers

import (
	"net/http"
	"time"

	"lovelion/internal/models"
	"lovelion/internal/services"
	"lovelion/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AnnouncementHandler serves public announcement endpoints.
type AnnouncementHandler struct {
	db *gorm.DB
}

func NewAnnouncementHandler(db *gorm.DB) *AnnouncementHandler {
	return &AnnouncementHandler{db: db}
}

// List returns all published announcements.
func (h *AnnouncementHandler) List(c *gin.Context) {
	var announcements []models.Announcement
	if err := h.db.Where("status = ?", "published").Order("created_at DESC").Find(&announcements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch announcements"})
		return
	}
	c.JSON(http.StatusOK, announcements)
}

// Get returns a single published announcement.
func (h *AnnouncementHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var announcement models.Announcement
	if err := h.db.Where("id = ? AND status = ?", id, "published").First(&announcement).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Announcement not found"})
		return
	}
	c.JSON(http.StatusOK, announcement)
}

// Broadcast returns the latest currently broadcasting announcement.
func (h *AnnouncementHandler) Broadcast(c *gin.Context) {
	now := time.Now()
	var announcement models.Announcement
	err := h.db.
		Where("status = ? AND broadcast_start <= ? AND broadcast_end >= ?", "published", now, now).
		Order("created_at DESC").
		First(&announcement).Error
	if err != nil {
		// No active broadcast — return null instead of error
		c.JSON(http.StatusOK, nil)
		return
	}
	c.JSON(http.StatusOK, announcement)
}

// AdminAnnouncementHandler serves admin-only announcement endpoints.
type AdminAnnouncementHandler struct {
	db        *gorm.DB
	generator *services.AnnouncementGenerator
}

func NewAdminAnnouncementHandler(db *gorm.DB, generator *services.AnnouncementGenerator) *AdminAnnouncementHandler {
	return &AdminAnnouncementHandler{db: db, generator: generator}
}

// List returns all announcements (including drafts) for admin.
func (h *AdminAnnouncementHandler) List(c *gin.Context) {
	var announcements []models.Announcement
	if err := h.db.Order("created_at DESC").Find(&announcements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch announcements"})
		return
	}
	c.JSON(http.StatusOK, announcements)
}

type createAnnouncementRequest struct {
	Title          string     `json:"title" binding:"required,min=1,max=255"`
	Content        string     `json:"content"`
	Status         string     `json:"status" binding:"required,oneof=draft published"`
	BroadcastStart *time.Time `json:"broadcast_start"`
	BroadcastEnd   *time.Time `json:"broadcast_end"`
}

// Create creates a new announcement.
func (h *AdminAnnouncementHandler) Create(c *gin.Context) {
	var req createAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := utils.NewShortID(h.db, "announcements", "id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate ID"})
		return
	}

	announcement := models.Announcement{
		ID:             id,
		Title:          req.Title,
		Content:        req.Content,
		Status:         req.Status,
		BroadcastStart: req.BroadcastStart,
		BroadcastEnd:   req.BroadcastEnd,
	}

	if err := h.db.Create(&announcement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create announcement"})
		return
	}

	c.JSON(http.StatusCreated, announcement)
}

type updateAnnouncementRequest struct {
	Title          string     `json:"title" binding:"required,min=1,max=255"`
	Content        string     `json:"content"`
	Status         string     `json:"status" binding:"required,oneof=draft published"`
	BroadcastStart *time.Time `json:"broadcast_start"`
	BroadcastEnd   *time.Time `json:"broadcast_end"`
}

// Update updates an existing announcement.
func (h *AdminAnnouncementHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var announcement models.Announcement
	if err := h.db.First(&announcement, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Announcement not found"})
		return
	}

	var req updateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	announcement.Title = req.Title
	announcement.Content = req.Content
	announcement.Status = req.Status
	announcement.BroadcastStart = req.BroadcastStart
	announcement.BroadcastEnd = req.BroadcastEnd

	if err := h.db.Save(&announcement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update announcement"})
		return
	}

	c.JSON(http.StatusOK, announcement)
}

// Delete hard-deletes an announcement.
func (h *AdminAnnouncementHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	result := h.db.Where("id = ?", id).Delete(&models.Announcement{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete announcement"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Announcement not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Announcement deleted"})
}

type generateAnnouncementRequest struct {
	Description string `json:"description" binding:"required,min=1"`
}

// Generate uses AI to generate announcement title and content from a description.
func (h *AdminAnnouncementHandler) Generate(c *gin.Context) {
	if h.generator == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "AI generation not available"})
		return
	}

	var req generateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.generator.Generate(c.Request.Context(), req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI generation failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
