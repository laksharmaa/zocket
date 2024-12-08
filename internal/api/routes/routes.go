package routes

import (
    "github.com/gin-gonic/gin"
    "product-management-system/internal/api/handlers"
    "product-management-system/internal/service"
    "log"
    "net/http"
)

func SetupRoutes(router *gin.Engine, productService *service.ProductService, logger *log.Logger) {
    // Root route
    router.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Welcome to Product Management System API",
            "status": "running",
        })
    })

    // Create handlers
    productHandler := handlers.NewProductHandler(productService)

    // API routes group
    api := router.Group("/api/v1")
    {
        // Product routes
        api.POST("/products", productHandler.CreateProduct)
        api.GET("/products/:id", productHandler.GetProduct)
        api.GET("/products", productHandler.GetProducts)
    }
}