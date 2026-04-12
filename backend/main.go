package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"lovelion/internal/config"
	"lovelion/internal/database"
	"lovelion/internal/handlers"
	"lovelion/internal/middleware"
	"lovelion/internal/repositories"
	"lovelion/internal/services"
	"lovelion/internal/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Setup structured logging
	if os.Getenv("GIN_MODE") == "release" {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	}

	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	// Setup Gin router
	r := gin.Default()

	// Request logging
	r.Use(middleware.RequestLogger())

	// CORS
	if len(cfg.CORSOrigins) > 0 {
		r.Use(cors.New(cors.Config{
			AllowOrigins:     cfg.CORSOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Authorization", "Content-Type"},
			ExposeHeaders:    []string{"Content-Length", "X-Total-Count"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// aiWorker is assigned inside the api block (where dependencies are in
	// scope). It stays nil when receipt extraction is disabled.
	var aiWorker *services.AIWorker

	// API routes
	api := r.Group("/api")
	{
		// Public routes
		authRateLimit := middleware.RateLimit(30, time.Minute)
		users := api.Group("/users")
		{
			authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret, cfg.JWTExpiry)
			users.POST("/register", authRateLimit, authHandler.Register)
			users.POST("/login", authRateLimit, authHandler.Login)

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

		// Shared R2 storage (used by ImageHandler and TransactionService)
		r2Storage, err := storage.NewR2Storage(cfg)
		if err != nil {
			slog.Error("failed to initialize R2 storage", "error", err)
			os.Exit(1)
		}

		// Services
		inviteService := services.NewInviteService(db, inviteRepo, memberRepo)
		txnService := services.NewTransactionService(db, txnRepo, expenseRepo, expenseItemRepo, debtRepo, r2Storage)

		// AI receipt extraction rate limiter (per-user daily cap).
		// A zero/negative cap disables the check entirely.
		aiRateLimiter := middleware.NewAIRateLimiter(cfg.ReceiptRateLimitPerDay)

		// AI worker — polls transactions with ai_status='pending' and calls
		// the Gemini extractor. Disabled if the feature flag is off or the
		// API key is missing.
		if cfg.ReceiptExtractEnabled {
			if cfg.GeminiAPIKey == "" {
				slog.Warn("receipt extraction enabled but GEMINI_API_KEY is empty — worker not started")
			} else {
				extractor := services.NewGeminiReceiptExtractor(cfg.GeminiAPIKey, cfg.GeminiModel, cfg.GeminiBaseURL)
				aiWorker = services.NewAIWorker(db, extractor, r2Storage, services.AIWorkerConfig{})
			}
		} else {
			slog.Info("receipt extraction disabled", "RECEIPT_EXTRACT_ENABLED", "false")
		}

		// Public announcement routes
		announcementHandler := handlers.NewAnnouncementHandler(db)
		announcements := api.Group("/announcements")
		{
			announcements.GET("", announcementHandler.List)
			announcements.GET("/broadcast", announcementHandler.Broadcast)
			announcements.GET("/:id", announcementHandler.Get)
		}

		// Admin announcement routes
		var announcementGenerator *services.AnnouncementGenerator
		if cfg.GeminiAPIKey != "" {
			announcementGenerator = services.NewAnnouncementGenerator(cfg.GeminiAPIKey, cfg.GeminiModel, cfg.GeminiBaseURL)
		}
		adminAnnouncementHandler := handlers.NewAdminAnnouncementHandler(db, announcementGenerator)
		adminGroup := api.Group("/admin")
		adminGroup.Use(middleware.AuthRequiredWithDB(cfg.JWTSecret, db), middleware.AdminOnly(db))
		{
			adminAnnouncements := adminGroup.Group("/announcements")
			{
				adminAnnouncements.GET("", adminAnnouncementHandler.List)
				adminAnnouncements.GET("/config", adminAnnouncementHandler.Config)
				adminAnnouncements.POST("", adminAnnouncementHandler.Create)
				adminAnnouncements.PUT("/:id", adminAnnouncementHandler.Update)
				adminAnnouncements.DELETE("/:id", adminAnnouncementHandler.Delete)
				adminAnnouncements.POST("/generate", adminAnnouncementHandler.Generate)
			}
		}

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

				// Transaction routes (shared: list, get, delete, ai-cancel)
				transactionHandler := handlers.NewTransactionHandler(txnService)
				spaceGroup.GET("/transactions", transactionHandler.List)
				spaceGroup.GET("/transactions/:txn_id", transactionHandler.Get)
				spaceGroup.DELETE("/transactions/:txn_id", transactionHandler.Delete)
				spaceGroup.POST("/transactions/:txn_id/ai-cancel", transactionHandler.AICancel)

				// Expense routes
				expenseHandler := handlers.NewExpenseHandler(txnService, aiRateLimiter)
				spaceGroup.POST("/expenses", expenseHandler.Create)
				spaceGroup.PUT("/expenses/:txn_id", expenseHandler.Update)

				// Payment routes
				paymentHandler := handlers.NewPaymentHandler(txnService)
				spaceGroup.POST("/payments", paymentHandler.Create)
				spaceGroup.PUT("/payments/:txn_id", paymentHandler.Update)

				// Expense template routes
				templateHandler := handlers.NewExpenseTemplateHandler(db)
				spaceGroup.GET("/expense-templates", templateHandler.List)
				spaceGroup.POST("/expense-templates", templateHandler.Create)
				spaceGroup.DELETE("/expense-templates/:template_id", templateHandler.Delete)
			}
		}

		// Image routes
		images := api.Group("/images")
		images.Use(middleware.AuthRequiredWithDB(cfg.JWTSecret, db))
		{
			imageHandler := handlers.NewImageHandler(db, r2Storage)
			images.POST("", imageHandler.Upload)
			images.GET("", imageHandler.List)
			images.PUT("/order", imageHandler.Reorder)
			images.DELETE("/:id", imageHandler.Delete)
		}
	}
	// Start server with graceful shutdown
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start AI worker (if enabled) with its own cancellable context so we can
	// stop it independently of the HTTP server. We also hold a WaitGroup so
	// shutdown blocks until the worker has returned from Run().
	workerCtx, cancelWorker := context.WithCancel(context.Background())
	defer cancelWorker()
	var workerWG sync.WaitGroup
	if aiWorker != nil {
		workerWG.Add(1)
		go func() {
			defer workerWG.Done()
			aiWorker.Run(workerCtx)
		}()
	}

	go func() {
		slog.Info("server starting", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server...")

	// Stop the worker first so no new LLM calls or DB writes start during the
	// HTTP drain window. A processing row at this point will be re-queued on
	// the next boot via recoverStuck.
	cancelWorker()
	if aiWorker != nil {
		workerWG.Wait()
		slog.Info("ai worker stopped")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server exited gracefully")
}
