package handlers

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"lovelion/internal/models"
	"lovelion/internal/testutil"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// mockS3Client creates a client pointing to a dummy server
func mockS3Client(t *testing.T) (*s3.Client, *httptest.Server) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: server.URL,
		}, nil
	})

	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithEndpointResolverWithOptions(resolver),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")),
		awsconfig.WithRegion("us-east-1"),
	)
	if err != nil {
		t.Fatalf("Failed to load aws config: %v", err)
	}

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	}), server
}

func TestImageHandler_List(t *testing.T) {
	db := testutil.TestDB(t)

	// Create dummy data
	entityID := "trip_123"
	entityType := "trip"

	images := []models.Image{
		{
			ID:         uuid.New(),
			EntityID:   entityID,
			EntityType: entityType,
			FilePath:   "https://example.com/img1.jpg",
			SortOrder:  1,
		},
		{
			ID:         uuid.New(),
			EntityID:   entityID,
			EntityType: entityType,
			FilePath:   "https://example.com/img2.jpg",
			SortOrder:  0,
		},
	}

	for _, img := range images {
		if err := db.Create(&img).Error; err != nil {
			t.Fatalf("Failed to create image: %v", err)
		}
	}

	// Setup handler
	handler := &ImageHandler{db: db}
	router := testutil.TestRouter()
	router.GET("/api/images", handler.List)

	// Test case 1: List by entity_id and entity_type
	w := httptest.NewRecorder()
	req := testutil.JSONRequest("GET", "/api/images?entity_id="+entityID+"&entity_type="+entityType, nil)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)

	var result []models.Image
	testutil.ParseResponse(t, w, &result)

	if len(result) != 2 {
		t.Errorf("Expected 2 images, got %d", len(result))
	}

	// Check sorting (should be by SortOrder ASC, so img2 (0) then img1 (1))
	if result[0].SortOrder != 0 {
		t.Error("Images should be sorted by sort_order ASC")
	}
}

func TestImageHandler_Reorder(t *testing.T) {
	db := testutil.TestDB(t)

	// Create dummy data
	entityID := "trip_123"
	entityType := "trip"

	img1 := models.Image{ID: uuid.New(), EntityID: entityID, EntityType: entityType, SortOrder: 0}
	img2 := models.Image{ID: uuid.New(), EntityID: entityID, EntityType: entityType, SortOrder: 1}

	db.Create(&img1)
	db.Create(&img2)

	handler := &ImageHandler{db: db}
	router := testutil.TestRouter()
	router.POST("/api/images/reorder", handler.Reorder)

	// Reorder swap
	reqBody := ReorderRequest{
		IDs: []string{img2.ID.String(), img1.ID.String()},
	}

	w := httptest.NewRecorder()
	req := testutil.JSONRequest("POST", "/api/images/reorder", reqBody)
	router.ServeHTTP(w, req)
	testutil.ExpectStatus(t, w, 200)

	// Verify new order
	var loadedImg1, loadedImg2 models.Image
	db.First(&loadedImg1, "id = ?", img1.ID)
	db.First(&loadedImg2, "id = ?", img2.ID)

	if loadedImg1.SortOrder != 1 {
		t.Errorf("Img1 should have sort_order 1, got %d", loadedImg1.SortOrder)
	}
	if loadedImg2.SortOrder != 0 {
		t.Errorf("Img2 should have sort_order 0, got %d", loadedImg2.SortOrder)
	}
}

func TestImageHandler_Upload(t *testing.T) {
	db := testutil.TestDB(t)
	s3Client, server := mockS3Client(t)
	defer server.Close()

	handler := &ImageHandler{
		db:           db,
		s3Client:     s3Client,
		bucket:       "test-bucket",
		publicDomain: "https://r2.example.com",
	}

	router := testutil.TestRouter()
	router.POST("/api/images", handler.Upload)

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "test.jpg")
	if err != nil {
		t.Fatal(err)
	}
	io.Copy(part, bytes.NewBufferString("fake image content"))

	writer.WriteField("entity_id", "trip_123")
	writer.WriteField("entity_type", "trip")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/images", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	testutil.ExpectStatus(t, w, 201)

	var created models.Image
	testutil.ParseResponse(t, w, &created)

	if created.EntityID != "trip_123" {
		t.Errorf("Expected EntityID trip_123, got %s", created.EntityID)
	}
}
