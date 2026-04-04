package main

import (
	"log"
	"time"

	"lovelion/internal/config"
	"lovelion/internal/database"
	"lovelion/internal/handlers"
	"lovelion/internal/middleware"
	"lovelion/internal/repositories"
	"lovelion/internal/services"

	"github.com/gin-contrib/cors"
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

	// CORS
	if len(cfg.CORSOrigins) > 0 {
		r.Use(cors.New(cors.Config{
			AllowOrigins:     cfg.CORSOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Authorization", "Content-Type"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}

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
			users.PUT("/me", middleware.AuthRequiredWithDB(cfg.JWTSecret, db), authHandler.UpdateMe)
		}

		// Repositories
		inviteRepo := repositories.NewInviteRepo(db)
		memberRepo := repositories.NewMemberRepo(db)
		txnRepo := repositories.NewTransactionRepo(db)
		expenseRepo := repositories.NewTransactionExpenseRepo(db)
		expenseItemRepo := repositories.NewTransactionExpenseItemRepo(db)
		debtRepo := repositories.NewTransactionDebtRepo(db)

		// Services
		inviteService := services.NewInviteService(db, inviteRepo, memberRepo)
		txnService := services.NewTransactionService(db, txnRepo, expenseRepo, expenseItemRepo, debtRepo)

		// Sharing routes (Public Info)
		sharingHandler := handlers.NewSpaceSharingHandler(inviteService, memberRepo)
		api.GET("/invites/:token", sharingHandler.GetInviteInfo)
		api.POST("/invites/:token/join", middleware.AuthRequiredWithDB(cfg.JWTSecret, db), sharingHandler.JoinSpace)

		// Unified Space routes
		spaces := api.Group("/spaces")
		spaces.Use(middleware.AuthRequiredWithDB(cfg.JWTSecret, db))
		{
			spaceHandler := handlers.NewSpaceHandler(db)
			spaces.GET("", spaceHandler.List)
			spaces.POST("", spaceHandler.Create)

			// Single space operations
			spaceGroup := spaces.Group("/:id")
			spaceGroup.Use(middleware.SpaceAccess(db))
			{
				spaceGroup.GET("", spaceHandler.Get)
				spaceGroup.POST("/leave", spaceHandler.Leave)

				// Owner only operations
				ownerGroup := spaceGroup.Group("")
				ownerGroup.Use(middleware.SpaceOwnerOnly())
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

				// Transaction routes (shared: list, get, delete)
				transactionHandler := handlers.NewTransactionHandler(txnService)
				spaceGroup.GET("/transactions", transactionHandler.List)
				spaceGroup.GET("/transactions/:txn_id", transactionHandler.Get)
				spaceGroup.DELETE("/transactions/:txn_id", transactionHandler.Delete)

				// Expense routes
				expenseHandler := handlers.NewExpenseHandler(txnService)
				spaceGroup.POST("/expenses", expenseHandler.Create)
				spaceGroup.PUT("/expenses/:txn_id", expenseHandler.Update)

				// Payment routes
				paymentHandler := handlers.NewPaymentHandler(txnService)
				spaceGroup.POST("/payments", paymentHandler.Create)
				spaceGroup.PUT("/payments/:txn_id", paymentHandler.Update)
			}
		}

		// Image routes
		images := api.Group("/images")
		images.Use(middleware.AuthRequiredWithDB(cfg.JWTSecret, db))
		{
			imageHandler, err := handlers.NewImageHandler(db)
			if err != nil {
				log.Fatalf("Failed to initialize image handler: %v", err)
			}
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
