package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BradenHooton/pinecone-api/internal/auth"
	"github.com/BradenHooton/pinecone-api/internal/config"
	"github.com/BradenHooton/pinecone-api/internal/mealplan"
	"github.com/BradenHooton/pinecone-api/internal/middleware"
	"github.com/BradenHooton/pinecone-api/internal/nutrition"
	"github.com/BradenHooton/pinecone-api/internal/recipe"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup structured logging
	logLevel := slog.LevelInfo
	switch cfg.LogLevel {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	logger.Info("Starting Pinecone API server",
		slog.String("version", "1.0.0"),
		slog.String("port", cfg.ServerPort),
		slog.String("log_level", cfg.LogLevel),
	)

	// Connect to database
	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		logger.Error("Unable to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer dbpool.Close()

	// Test database connection
	if err := dbpool.Ping(context.Background()); err != nil {
		logger.Error("Unable to ping database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Info("Database connection established")

	// Initialize repositories
	authRepo := auth.NewPostgresRepository(dbpool)
	recipeRepo := recipe.NewPostgresRepository(dbpool)
	nutritionRepo := nutrition.NewPostgresRepository(dbpool)
	mealPlanRepo := mealplan.NewPostgresRepository(dbpool)

	// Initialize USDA client (stub for now)
	usdaClient := nutrition.NewStubUSDAClient()

	// Initialize services
	authService := auth.NewService(authRepo, cfg.JWTSecret, cfg.JWTExpiryHours)
	recipeService := recipe.NewService(recipeRepo)
	nutritionService := nutrition.NewService(nutritionRepo, usdaClient)
	mealPlanService := mealplan.NewService(mealPlanRepo)

	// Initialize handlers
	authHandler := auth.NewHandler(authService)
	recipeHandler := recipe.NewHandler(recipeService, cfg.UploadDir)
	nutritionHandler := nutrition.NewHandler(nutritionService)
	mealPlanHandler := mealplan.NewHandler(mealPlanService)

	// Create router
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger(logger))
	r.Use(middleware.CORS(cfg.AllowedOrigins))
	r.Use(httprate.LimitByIP(100, 1*time.Minute)) // 100 requests per minute per IP

	// Health check endpoint (no auth required)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes (no auth required)
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", authHandler.HandleRegister)
			r.Post("/auth/login", authHandler.HandleLogin)
		})

		// Protected routes (auth required)
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(cfg.JWTSecret))

			r.Post("/auth/logout", authHandler.HandleLogout)

			// Recipe routes
			r.Route("/recipes", func(r chi.Router) {
				r.Get("/", recipeHandler.HandleList)
				r.Post("/", recipeHandler.HandleCreate)
				r.Post("/upload-image", recipeHandler.HandleUploadImage)
				r.Get("/{id}", recipeHandler.HandleGetByID)
				r.Put("/{id}", recipeHandler.HandleUpdate)
				r.Delete("/{id}", recipeHandler.HandleDelete)
			})

			// Nutrition routes
			r.Get("/nutrition/search", nutritionHandler.HandleSearch)

			// Meal plan routes
			r.Route("/meal-plans", func(r chi.Router) {
				r.Get("/", mealPlanHandler.HandleGetByDateRange)
				r.Get("/date", mealPlanHandler.HandleGetByDate)
				r.Put("/date", mealPlanHandler.HandleUpdate)
			})
		})
	})

	// Create HTTP server
	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server listening", slog.String("address", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("Server exited gracefully")
}
