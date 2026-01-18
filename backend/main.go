package main

import (
	"log"

	"lovelion/internal/config"
	"lovelion/internal/database"
	"lovelion/internal/handlers"
	"lovelion/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Setup Gin router
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api")
	{
		// Public routes
		users := api.Group("/users")
		{
			authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret)
			users.POST("/register", authHandler.Register)
			users.POST("/login", authHandler.Login)

			// Protected routes
			users.GET("/me", middleware.AuthRequired(cfg.JWTSecret), authHandler.GetMe)
		}
	}

	// Start server
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
