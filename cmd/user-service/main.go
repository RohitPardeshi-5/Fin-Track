package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fintrack/internal/common"
	"fintrack/internal/user"
	"fintrack/pkg/config"
	"fintrack/pkg/database"
	"fintrack/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Connect to database
	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Auto-migrate models
	if err := db.AutoMigrate(&common.User{}); err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}

	// Initialize services
	userService := user.NewService(db, cfg.JWTSecret)
	userHandler := user.NewHandler(userService, logger)

	// Setup router
	router := gin.New()
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())

	// Health check
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "user-service"})
	})

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes
	api := router.Group("/api/v1")
	userHandler.SetupRoutes(api.Group("/users"))

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		logger.Info("Starting user service", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}