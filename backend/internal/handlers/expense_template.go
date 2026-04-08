package handlers

import (
	"net/http"

	"lovelion/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExpenseTemplateHandler struct {
	db *gorm.DB
}

func NewExpenseTemplateHandler(db *gorm.DB) *ExpenseTemplateHandler {
	return &ExpenseTemplateHandler{db: db}
}

type CreateExpenseTemplateRequest struct {
	Name string                     `json:"name" binding:"required,min=1,max=100"`
	Data models.ExpenseTemplateData `json:"data" binding:"required"`
}

func (h *ExpenseTemplateHandler) List(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Space)

	var templates []models.ExpenseTemplate
	if err := h.db.Where("space_id = ?", space.ID).Order("created_at DESC").Find(&templates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch templates"})
		return
	}

	c.JSON(http.StatusOK, templates)
}

func (h *ExpenseTemplateHandler) Create(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Space)

	var req CreateExpenseTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template := models.ExpenseTemplate{
		SpaceID: space.ID,
		Name:    req.Name,
		Data:    req.Data,
	}

	if err := h.db.Create(&template).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create template"})
		return
	}

	c.JSON(http.StatusCreated, template)
}

func (h *ExpenseTemplateHandler) Delete(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Space)

	templateID, err := uuid.Parse(c.Param("template_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	result := h.db.Where("id = ? AND space_id = ?", templateID, space.ID).Delete(&models.ExpenseTemplate{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete template"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template deleted"})
}
