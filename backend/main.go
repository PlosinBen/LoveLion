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
			users.GET("/me", middleware.AuthRequiredWithDB(cfg.JWTSecret, db), authHandler.GetMe)
		}

		// Sharing routes (Public Info)
		sharingHandler := handlers.NewLedgerSharingHandler(db)
		api.GET("/invites/:token", sharingHandler.GetInviteInfo)
		api.POST("/invites/:token/join", middleware.AuthRequiredWithDB(cfg.JWTSecret, db), sharingHandler.JoinLedger)

		// Unified Space (Ledger) routes
		spaces := api.Group("/spaces")
		spaces.Use(middleware.AuthRequiredWithDB(cfg.JWTSecret, db))
		{
			spaceHandler := handlers.NewSpaceHandler(db)
			spaces.GET("", spaceHandler.List)
			spaces.POST("", spaceHandler.Create)

			// Single space operations
			spaceGroup := spaces.Group("/:id")
			spaceGroup.Use(middleware.LedgerAccess(db))
			{
				spaceGroup.GET("", spaceHandler.Get)
				
				// Owner only operations
				ownerGroup := spaceGroup.Group("")
				ownerGroup.Use(middleware.LedgerOwnerOnly())
				{
					ownerGroup.PUT("", spaceHandler.Update)
					ownerGroup.PATCH("", spaceHandler.Update) // Add PATCH support
					ownerGroup.DELETE("", spaceHandler.Delete)
					
					// Invitation management
					ownerGroup.POST("/invites", sharingHandler.CreateInvite)
					ownerGroup.GET("/invites", sharingHandler.ListInvites)
					ownerGroup.DELETE("/invites/:invite_id", sharingHandler.RevokeInvite)
				}

				// Member management
				spaceGroup.GET("/members", sharingHandler.ListMembers)
				spaceGroup.PATCH("/members/:user_id", sharingHandler.UpdateMemberAlias)
				spaceGroup.DELETE("/members/:user_id", sharingHandler.RemoveMember)

				// Comparison routes (Integrated into space)
				comparisonHandler := handlers.NewComparisonHandler(db)
				spaceGroup.GET("/stores", comparisonHandler.ListStores)
				spaceGroup.POST("/stores", comparisonHandler.CreateStore)
				spaceGroup.GET("/stores/:store_id", comparisonHandler.GetStore)
				spaceGroup.PUT("/stores/:store_id", comparisonHandler.UpdateStore)
				spaceGroup.DELETE("/stores/:store_id", comparisonHandler.DeleteStore)
				spaceGroup.GET("/products", comparisonHandler.ListAllProducts)
				spaceGroup.POST("/stores/:store_id/products", comparisonHandler.CreateProduct)
				spaceGroup.GET("/stores/:store_id/products/:product_id", comparisonHandler.GetProduct)
				spaceGroup.PUT("/stores/:store_id/products/:product_id", comparisonHandler.UpdateProduct)
				spaceGroup.DELETE("/stores/:store_id/products/:product_id", comparisonHandler.DeleteProduct)

				// Transaction routes nested under space
				transactionHandler := handlers.NewTransactionHandler(db)
				spaceGroup.GET("/transactions", transactionHandler.List)
				spaceGroup.POST("/transactions", transactionHandler.Create)
				spaceGroup.GET("/transactions/:txn_id", transactionHandler.Get)
				spaceGroup.PUT("/transactions/:txn_id", transactionHandler.Update)
				spaceGroup.DELETE("/transactions/:txn_id", transactionHandler.Delete)
			}
		}

		// Backward compatibility: Aliasing /ledgers and /trips to /spaces
		ledgerHandler := handlers.NewSpaceHandler(db)
		ledgers := api.Group("/ledgers")
		ledgers.Use(middleware.AuthRequiredWithDB(cfg.JWTSecret, db))
		{
			ledgers.GET("", ledgerHandler.List)
			ledgers.POST("", ledgerHandler.Create)
			
			ledgerGroup := ledgers.Group("/:id")
			ledgerGroup.Use(middleware.LedgerAccess(db))
			{
				ledgerGroup.GET("", ledgerHandler.Get)
				ownerGroup := ledgerGroup.Group("")
				ownerGroup.Use(middleware.LedgerOwnerOnly())
				{
					ownerGroup.PUT("", ledgerHandler.Update)
					ownerGroup.PATCH("", ledgerHandler.Update)
					ownerGroup.DELETE("", ledgerHandler.Delete)
					ownerGroup.POST("/invites", sharingHandler.CreateInvite)
					ownerGroup.GET("/invites", sharingHandler.ListInvites)
					ownerGroup.DELETE("/invites/:invite_id", sharingHandler.RevokeInvite)
				}
				ledgerGroup.GET("/members", sharingHandler.ListMembers)
				ledgerGroup.PATCH("/members/:user_id", sharingHandler.UpdateMemberAlias)
				ledgerGroup.DELETE("/members/:user_id", sharingHandler.RemoveMember)

				transactionHandler := handlers.NewTransactionHandler(db)
				ledgerGroup.GET("/transactions", transactionHandler.List)
				ledgerGroup.POST("/transactions", transactionHandler.Create)
				ledgerGroup.GET("/transactions/:txn_id", transactionHandler.Get)
				ledgerGroup.PUT("/transactions/:txn_id", transactionHandler.Update)
				ledgerGroup.DELETE("/transactions/:txn_id", transactionHandler.Delete)
			}
		}

		trips := api.Group("/trips")
		trips.Use(middleware.AuthRequiredWithDB(cfg.JWTSecret, db))
		{
			// For List, we can filter by type=trip automatically in the handler or here
			// Let's use the SpaceHandler's List and rely on query params from frontend if they update, 
			// or we could wrap it. For compatibility, we just map them.
			trips.GET("", ledgerHandler.List) 
			trips.POST("", ledgerHandler.Create)
			trips.GET("/:id", ledgerHandler.Get)
			trips.PUT("/:id", ledgerHandler.Update)
			trips.DELETE("/:id", ledgerHandler.Delete)
			
			// For trips, members and comparisons are now under the same ID
			comparisonHandler := handlers.NewComparisonHandler(db)
			trips.GET("/:id/members", sharingHandler.ListMembers)
			trips.GET("/:id/stores", comparisonHandler.ListStores)
			trips.POST("/:id/stores", comparisonHandler.CreateStore)
			trips.GET("/:id/stores/:store_id", comparisonHandler.GetStore)
			trips.PUT("/:id/stores/:store_id", comparisonHandler.UpdateStore)
			trips.DELETE("/:id/stores/:store_id", comparisonHandler.DeleteStore)
			trips.GET("/:id/products", comparisonHandler.ListAllProducts)
			trips.POST("/:id/stores/:store_id/products", comparisonHandler.CreateProduct)
			trips.GET("/:id/stores/:store_id/products/:product_id", comparisonHandler.GetProduct)
			trips.PUT("/:id/stores/:store_id/products/:product_id", comparisonHandler.UpdateProduct)
			trips.DELETE("/:id/stores/:store_id/products/:product_id", comparisonHandler.DeleteProduct)
		}

		// Image routes
		images := api.Group("/images")
		images.Use(middleware.AuthRequiredWithDB(cfg.JWTSecret, db))
		{
			imageHandler := handlers.NewImageHandler(db)
			images.POST("", imageHandler.Upload)
			images.GET("", imageHandler.List)
			images.PUT("/order", imageHandler.Reorder)
			images.DELETE("/:id", imageHandler.Delete)
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
