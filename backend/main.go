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

		// Protected ledger routes
		ledgers := api.Group("/ledgers")
		ledgers.Use(middleware.AuthRequired(cfg.JWTSecret))
		{
			ledgerHandler := handlers.NewLedgerHandler(db)
			ledgers.GET("", ledgerHandler.List)
			ledgers.POST("", ledgerHandler.Create)
			ledgers.GET("/:id", ledgerHandler.Get)
			ledgers.PUT("/:id", ledgerHandler.Update)
			ledgers.DELETE("/:id", ledgerHandler.Delete)

			// Transaction routes nested under ledger
			transactionHandler := handlers.NewTransactionHandler(db)
			ledgers.GET("/:id/transactions", transactionHandler.List)
			ledgers.POST("/:id/transactions", transactionHandler.Create)
			ledgers.GET("/:id/transactions/:txn_id", transactionHandler.Get)
			ledgers.PUT("/:id/transactions/:txn_id", transactionHandler.Update)
			ledgers.DELETE("/:id/transactions/:txn_id", transactionHandler.Delete)
		}

		// Protected trip routes
		trips := api.Group("/trips")
		trips.Use(middleware.AuthRequired(cfg.JWTSecret))
		{
			tripHandler := handlers.NewTripHandler(db)
			trips.GET("", tripHandler.List)
			trips.POST("", tripHandler.Create)
			trips.GET("/:id", tripHandler.Get)
			trips.PUT("/:id", tripHandler.Update)
			trips.DELETE("/:id", tripHandler.Delete)

			// Member routes
			trips.GET("/:id/members", tripHandler.ListMembers)
			trips.POST("/:id/members", tripHandler.AddMember)
			trips.DELETE("/:id/members/:member_id", tripHandler.RemoveMember)

			// Comparison routes
			comparisonHandler := handlers.NewComparisonHandler(db)
			trips.GET("/:id/stores", comparisonHandler.ListStores)
			trips.POST("/:id/stores", comparisonHandler.CreateStore)
			trips.GET("/:id/stores/:store_id", comparisonHandler.GetStore)
			trips.DELETE("/:id/stores/:store_id", comparisonHandler.DeleteStore)
			trips.GET("/:id/products", comparisonHandler.ListAllProducts)
			trips.POST("/:id/stores/:store_id/products", comparisonHandler.CreateProduct)
			trips.PUT("/:id/stores/:store_id/products/:product_id", comparisonHandler.UpdateProduct)
			trips.DELETE("/:id/stores/:store_id/products/:product_id", comparisonHandler.DeleteProduct)
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
