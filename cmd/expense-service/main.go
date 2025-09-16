package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fintrack/internal/common"
	"fintrack/internal/expense"
	"fintrack/pkg/config"
	"fintrack/pkg/database"
	"fintrack/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	if err := db.AutoMigrate(&common.User{}, &common.Expense{}); err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}

	expenseService := expense.NewService(db)
	expenseHandler := expense.NewHandler(expenseService, logger)

	router := gin.New()
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "expense-service"})
	})

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	api := router.Group("/api/v1")
	protected := api.Group("/expenses")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	expenseHandler.SetupRoutes(protected)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		logger.Info("Starting expense service", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

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