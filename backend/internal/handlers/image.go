package handlers

import (
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"path/filepath"
	"strings"

	"lovelion/internal/config"
	"lovelion/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bbrks/go-blurhash"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ImageHandler struct {
	db           *gorm.DB
	s3Client     *s3.Client
	bucket       string
	publicDomain string
}

func NewImageHandler(db *gorm.DB) *ImageHandler {
	cfg := config.Load()

	// Setup S3 Client for R2
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.R2AccountID),
		}, nil
	})

	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithEndpointResolverWithOptions(r2Resolver),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.R2AccessKey, cfg.R2SecretKey, "")),
		awsconfig.WithRegion("auto"),
	)
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := s3.NewFromConfig(awsCfg)

	return &ImageHandler{
		db:           db,
		s3Client:     client,
		bucket:       cfg.R2Bucket,
		publicDomain: cfg.R2PublicDomain,
	}
}

// Upload handles file upload and creates a database record
func (h *ImageHandler) Upload(c *gin.Context) {
	// Parse form
	entityID := c.PostForm("entity_id")
	entityType := c.PostForm("entity_type")

	if entityID == "" || entityType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "entity_id and entity_type are required"})
		return
	}

	// Get file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// Validation
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only jpg, jpeg, and png are allowed"})
		return
	}

	// Open file
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	// Generate BlurHash
	// We need to decode the image first. Since we need to upload the original file reader later,
	// and decoding consumes the reader, we might need to handle this carefully.
	// However, multipart.File implements Seek, so we can rewind.
	imgData, _, err := image.Decode(f)
	if err != nil {
		// Log warning but don't fail upload? Or fail?
		// For now, if not an image we can decode, just skip blurhash or error.
		// User only uploads jpg/png validated above, so should be fine.
		fmt.Printf("Failed to decode image for blurhash: %v\n", err)
	}

	var blurHashStr string
	if imgData != nil {
		// Components x=4, y=3 usually good compromise
		blurHashStr, err = blurhash.Encode(4, 3, imgData)
		if err != nil {
			fmt.Printf("Failed to encode blurhash: %v\n", err)
		}
	}

	// Rewind file for S3 upload
	if _, err := f.Seek(0, 0); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file"})
		return
	}

	// Generate key
	fileUUID := uuid.New()
	key := fmt.Sprintf("%s/%s%s", entityType, fileUUID.String(), ext)

	// Upload to R2
	_, err = h.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(h.bucket),
		Key:         aws.String(key),
		Body:        f,
		ContentType: aws.String(file.Header.Get("Content-Type")),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to storage: " + err.Error()})
		return
	}

	// Create DB Record
	fullURL := fmt.Sprintf("%s/%s", h.publicDomain, key)
	image := models.Image{
		ID:         fileUUID,
		EntityID:   entityID,
		EntityType: entityType,
		FilePath:   fullURL, // Store Full URL
		BlurHash:   blurHashStr,
	}

	// Determine sort order
	var maxOrder int
	var count int64
	h.db.Model(&models.Image{}).Where("entity_id = ? AND entity_type = ?", entityID, entityType).Count(&count)
	if count > 0 {
		h.db.Model(&models.Image{}).Where("entity_id = ? AND entity_type = ?", entityID, entityType).Select("MAX(sort_order)").Scan(&maxOrder)
		image.SortOrder = maxOrder + 1
	} else {
		image.SortOrder = 0
	}

	if err := h.db.Create(&image).Error; err != nil {
		// Attempt to delete from R2 if DB insert fails
		h.s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Bucket: aws.String(h.bucket),
			Key:    aws.String(key),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image record"})
		return
	}

	c.JSON(http.StatusCreated, image)
}

// List images for an entity
func (h *ImageHandler) List(c *gin.Context) {
	entityID := c.Query("entity_id")
	entityType := c.Query("entity_type")
	entityIDs := c.Query("entity_ids") // comma separated

	query := h.db.Model(&models.Image{})

	if entityIDs != "" {
		ids := strings.Split(entityIDs, ",")
		query = query.Where("entity_id IN ?", ids)
	} else if entityID != "" {
		query = query.Where("entity_id = ?", entityID)
	}

	if entityType != "" {
		query = query.Where("entity_type = ?", entityType)
	}

	var images []models.Image
	if err := query.Order("sort_order ASC").Find(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch images"})
		return
	}

	c.JSON(http.StatusOK, images)
}

// ReorderRequest
type ReorderRequest struct {
	IDs []string `json:"ids"`
}

// Reorder images
func (h *ImageHandler) Reorder(c *gin.Context) {
	var req ReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		for i, id := range req.IDs {
			if err := tx.Model(&models.Image{}).Where("id = ?", id).Update("sort_order", i).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reorder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reordered successfully"})
}

// Delete image
func (h *ImageHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	var image models.Image
	if err := h.db.First(&image, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	// Delete from R2
	// Extract key from full URL
	// URL: https://domain.com/entity/uuid.ext
	// Key: entity/uuid.ext
	key := strings.TrimPrefix(image.FilePath, h.publicDomain+"/")

	_, err := h.s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(h.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		// Log error but continue to delete from DB?
		// Or fail? Usually better to fail or log.
		// For now, fail to alert user.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete from storage: " + err.Error()})
		return
	}

	if err := h.db.Delete(&image).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image deleted"})
}
