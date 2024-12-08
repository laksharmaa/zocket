package main

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"product-management-system/internal/api/middleware"
	"product-management-system/internal/api/routes"
	"product-management-system/internal/config"
	"product-management-system/internal/service"
	"product-management-system/pkg/cache"
	"product-management-system/pkg/cloudinary"
	"product-management-system/pkg/database"
	"product-management-system/pkg/queue"
	"syscall"
	"time"
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "[PRODUCT-API] ", log.LstdFlags)

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Cannot load config:", err)
	}

	// Initialize database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("Failed to get database instance:", err)
	}
	defer sqlDB.Close()

	// Initialize Redis cache
	redisCache, err := cache.NewRedisCache(cfg.RedisURL)
	if err != nil {
		logger.Fatal("Failed to connect to Redis:", err)
	}

	// Initialize RabbitMQ
	rabbitMQ, err := queue.NewRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		logger.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitMQ.Close()

	// Initialize Cloudinary client
	cloudinaryClient, err := cloudinary.NewCloudinaryClient(
		cfg.CloudinaryName,
		cfg.CloudinaryKey,
		cfg.CloudinarySecret,
	)
	if err != nil {
		logger.Fatal("Failed to initialize Cloudinary client:", err)
	}

	// Initialize services
	imageProcessor := service.NewImageProcessor(db, cloudinaryClient, rabbitMQ, logger)
	productService := service.NewProductService(db, redisCache, rabbitMQ, logger)

	// Start the image processing consumer in a goroutine
	go func() {
		logger.Println("Starting image processor...")
		if err := startImageProcessor(imageProcessor, rabbitMQ, logger); err != nil {
			logger.Printf("Image processor stopped: %v", err)
		}
	}()

	// Initialize Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(logger))

	// Setup routes
	routes.SetupRoutes(router, productService, logger)

	// Create server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Graceful shutdown handling
	go func() {
		logger.Printf("Server starting on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutting down server...")

	// Give outstanding operations 5 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", err)
	}

	logger.Println("Server exiting")
}

func startImageProcessor(processor *service.ImageProcessor, q *queue.RabbitMQ, logger *log.Logger) error {
	msgs, err := q.ConsumeMessages("image_processing")
	if err != nil {
		return err
	}

	for msg := range msgs {
		var processingMsg queue.ImageProcessingMessage
		if err := json.Unmarshal(msg.Body, &processingMsg); err != nil {
			logger.Printf("Error unmarshaling message: %v", err)
			msg.Nack(false, false)
			continue
		}

		if err := processor.ProcessImages(processingMsg); err != nil {
			logger.Printf("Error processing images: %v", err)
			msg.Nack(false, true) // Requeue the message
			continue
		}

		msg.Ack(false)
		logger.Printf("Successfully processed images for product ID: %d", processingMsg.ProductID)
	}

	return nil
}
