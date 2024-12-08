package handlers

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "product-management-system/internal/models"
    "product-management-system/internal/service"
)

type ProductHandler struct {
    productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
    return &ProductHandler{
        productService: productService,
    }
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
    var product models.Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.productService.CreateProduct(c.Request.Context(), &product); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    product, err := h.productService.GetProduct(c.Request.Context(), uint(id))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if product == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }

    c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
    userID, err := strconv.ParseUint(c.Query("user_id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    minPrice, _ := strconv.ParseFloat(c.Query("min_price"), 64)
    maxPrice, _ := strconv.ParseFloat(c.Query("max_price"), 64)
    productName := c.Query("product_name")

    products, err := h.productService.GetProducts(c.Request.Context(), uint(userID), minPrice, maxPrice, productName)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, products)
}