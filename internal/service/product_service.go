package service

import (
    "context"
    "log"
    "product-management-system/internal/models"
    "product-management-system/pkg/cache"
    "product-management-system/pkg/queue"
    "gorm.io/gorm"
)

type ProductService struct {
    db      *gorm.DB
    cache   *cache.RedisCache
    queue   *queue.RabbitMQ
    logger  *log.Logger
}

func NewProductService(db *gorm.DB, cache *cache.RedisCache, queue *queue.RabbitMQ, logger *log.Logger) *ProductService {
    return &ProductService{
        db:     db,
        cache:  cache,
        queue:  queue,
        logger: logger,
    }
}

func (s *ProductService) CreateProduct(ctx context.Context, product *models.Product) error {
    return s.db.Create(product).Error
}

func (s *ProductService) GetProduct(ctx context.Context, id uint) (*models.Product, error) {
    var product models.Product
    if err := s.db.First(&product, id).Error; err != nil {
        return nil, err
    }
    return &product, nil
}

func (s *ProductService) GetProducts(ctx context.Context, userID uint, minPrice, maxPrice float64, productName string) ([]models.Product, error) {
    var products []models.Product
    query := s.db.Where("user_id = ?", userID)

    if minPrice > 0 {
        query = query.Where("product_price >= ?", minPrice)
    }
    if maxPrice > 0 {
        query = query.Where("product_price <= ?", maxPrice)
    }
    if productName != "" {
        query = query.Where("product_name LIKE ?", "%"+productName+"%")
    }

    if err := query.Find(&products).Error; err != nil {
        return nil, err
    }
    return products, nil
}