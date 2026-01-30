package main

import (
	"fmt"
	"log"
	"lovelion/internal/config"
	"lovelion/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	dsn := cfg.DatabaseURL

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Connected to database")

	var images []models.Image
	err = db.Find(&images).Error
	if err != nil {
		log.Fatalf("Failed to query images: %v", err)
	}

	fmt.Printf("Total images found: %d\n", len(images))
	for _, img := range images {
		fmt.Printf("- Image: ID=%s, EntityID=%s, EntityType=%s, FilePath=%s\n", img.ID, img.EntityID, img.EntityType, img.FilePath)
	}

	// Check specific transaction
	txnID := "trip_txn01"
	var txnImages []models.Image
	db.Where("entity_id = ? AND entity_type = ?", txnID, "transaction").Find(&txnImages)
	fmt.Printf("\nImages for transaction %s: %d\n", txnID, len(txnImages))
}
