package main

import (
	"fmt"
	"log"
	"os"

	"lovelion/internal/models"
	"lovelion/internal/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: inv-init <username> [display_name]\n")
		fmt.Fprintf(os.Stderr, "  Creates the first investment owner for the given user.\n")
		os.Exit(1)
	}
	username := os.Args[1]

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("connect to db: %v", err)
	}

	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		log.Fatalf("user %q not found: %v", username, err)
	}

	var existing models.InvMember
	if err := db.Where("user_id = ?", user.ID).First(&existing).Error; err == nil {
		fmt.Printf("User %q is already an investment member (id=%s, is_owner=%v)\n", username, existing.ID, existing.IsOwner)
		return
	}

	displayName := user.DisplayName
	if len(os.Args) >= 3 {
		displayName = os.Args[2]
	}

	id := utils.MustNewShortID(db, "inv_members", "id")
	member := models.InvMember{
		ID:      id,
		Name:    displayName,
		UserID:  &user.ID,
		IsOwner: true,
		Active:  true,
	}

	if err := db.Create(&member).Error; err != nil {
		log.Fatalf("create inv_member: %v", err)
	}

	fmt.Printf("Created investment owner: %s (%s) → inv_member id=%s\n", username, displayName, id)
}
