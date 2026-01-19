package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// The migrations are in the "migrations" directory relative to the project root
	// (where we run the command from)
	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No migrations to apply")
		} else {
			log.Fatalf("Migration failed: %v", err)
		}
	} else {
		log.Println("Migrations applied successfully")
	}
}
