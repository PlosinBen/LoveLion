package handlers

import (
	"net/http"

	"lovelion/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ComparisonHandler struct {
	db *gorm.DB
}

func NewComparisonHandler(db *gorm.DB) *ComparisonHandler {
	return &ComparisonHandler{db: db}
}

type CreateStoreRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=100"`
	Location string `json:"location"`
}

type CreateProductRequest struct {
	Name     string          `json:"name" binding:"required,min=1,max=100"`
	Price    decimal.Decimal `json:"price" binding:"required"`
	Currency string          `json:"currency"`
	Unit     string          `json:"unit"`
	Note     string          `json:"note"`
}

type UpdateProductRequest struct {
	Name     string           `json:"name"`
	Price    *decimal.Decimal `json:"price"`
	Currency string           `json:"currency"`
	Unit     string           `json:"unit"`
	Note     string           `json:"note"`
}

// Helper to verify trip access
func (h *ComparisonHandler) verifyTripAccess(tripID string, userID uuid.UUID) error {
	var trip models.Trip
	if err := h.db.Where("id = ?", tripID).Preload("Members").First(&trip).Error; err != nil {
		return err
	}

	if trip.CreatedBy == userID {
		return nil
	}

	for _, m := range trip.Members {
		if m.UserID != nil && *m.UserID == userID {
			return nil
		}
	}

	return gorm.ErrRecordNotFound
}

// List stores for a trip
func (h *ComparisonHandler) ListStores(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	tripID := c.Param("id")

	if err := h.verifyTripAccess(tripID, userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	var stores []models.ComparisonStore
	if err := h.db.Where("trip_id = ?", tripID).Preload("Products").Find(&stores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stores"})
		return
	}

	c.JSON(http.StatusOK, stores)
}

// Create a store
func (h *ComparisonHandler) CreateStore(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	tripID := c.Param("id")

	if err := h.verifyTripAccess(tripID, userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	var req CreateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	storeID, err := gonanoid.New(21)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate ID"})
		return
	}

	store := &models.ComparisonStore{
		ID:       storeID,
		TripID:   tripID,
		Name:     req.Name,
		Location: req.Location,
	}

	if err := h.db.Create(store).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create store"})
		return
	}

	c.JSON(http.StatusCreated, store)
}

// Get a store with products
func (h *ComparisonHandler) GetStore(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	tripID := c.Param("id")
	storeID := c.Param("store_id")

	if err := h.verifyTripAccess(tripID, userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	var store models.ComparisonStore
	if err := h.db.Where("id = ? AND trip_id = ?", storeID, tripID).Preload("Products").First(&store).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}

	c.JSON(http.StatusOK, store)
}

// Delete a store
func (h *ComparisonHandler) DeleteStore(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	tripID := c.Param("id")
	storeID := c.Param("store_id")

	if err := h.verifyTripAccess(tripID, userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	result := h.db.Where("id = ? AND trip_id = ?", storeID, tripID).Delete(&models.ComparisonStore{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete store"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Store deleted"})
}

// List all products across stores for a trip (for comparison view)
func (h *ComparisonHandler) ListAllProducts(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	tripID := c.Param("id")

	if err := h.verifyTripAccess(tripID, userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	var products []models.ComparisonProduct
	if err := h.db.
		Joins("JOIN trip_comparison_stores ON trip_comparison_stores.id = trip_comparison_products.store_id").
		Where("trip_comparison_stores.trip_id = ?", tripID).
		Preload("Store").
		Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// Create a product in a store
func (h *ComparisonHandler) CreateProduct(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	tripID := c.Param("id")
	storeID := c.Param("store_id")

	if err := h.verifyTripAccess(tripID, userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	// Verify store belongs to trip
	var store models.ComparisonStore
	if err := h.db.Where("id = ? AND trip_id = ?", storeID, tripID).First(&store).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := &models.ComparisonProduct{
		ID:       uuid.New(),
		StoreID:  storeID,
		Name:     req.Name,
		Price:    req.Price,
		Currency: req.Currency,
		Unit:     req.Unit,
		Note:     req.Note,
	}

	if product.Currency == "" {
		product.Currency = "TWD"
	}

	if err := h.db.Create(product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// Update a product
func (h *ComparisonHandler) UpdateProduct(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	tripID := c.Param("id")
	storeID := c.Param("store_id")
	productID, err := uuid.Parse(c.Param("product_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := h.verifyTripAccess(tripID, userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	var product models.ComparisonProduct
	if err := h.db.Where("id = ? AND store_id = ?", productID, storeID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Currency != "" {
		product.Currency = req.Currency
	}
	product.Unit = req.Unit
	product.Note = req.Note

	if err := h.db.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// Delete a product
func (h *ComparisonHandler) DeleteProduct(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	tripID := c.Param("id")
	storeID := c.Param("store_id")
	productID, err := uuid.Parse(c.Param("product_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := h.verifyTripAccess(tripID, userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	result := h.db.Where("id = ? AND store_id = ?", productID, storeID).Delete(&models.ComparisonProduct{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}
