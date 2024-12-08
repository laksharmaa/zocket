package service

import (
    "context"
    "log"
    "product-management-system/internal/models"
    "product-management-system/pkg/cloudinary"
    "product-management-system/pkg/queue"
    "gorm.io/gorm"
)

type ImageProcessor struct {
    db         *gorm.DB
    cloudinary *cloudinary.Client
    queue      *queue.RabbitMQ
    logger     *log.Logger
}

func NewImageProcessor(db *gorm.DB, cloudinary *cloudinary.Client, queue *queue.RabbitMQ, logger *log.Logger) *ImageProcessor {
    return &ImageProcessor{
        db:         db,
        cloudinary: cloudinary,
        queue:      queue,
        logger:     logger,
    }
}

func (ip *ImageProcessor) ProcessImages(msg queue.ImageProcessingMessage) error {
    var compressedUrls []string
    
    for _, imageURL := range msg.Images {
        ip.logger.Printf("Processing image: %s", imageURL)
        
        // Upload and compress image using Cloudinary
        compressedURL, err := ip.cloudinary.UploadAndCompressImage(context.Background(), imageURL)
        if err != nil {
            ip.logger.Printf("Failed to process image %s: %v", imageURL, err)
            continue
        }
        
        ip.logger.Printf("Successfully compressed image. New URL: %s", compressedURL)
        compressedUrls = append(compressedUrls, compressedURL)
    }

    // Update product with compressed images
    if err := ip.db.Model(&models.Product{}).
        Where("id = ?", msg.ProductID).
        Update("compressed_product_images", compressedUrls).
        Error; err != nil {
        ip.logger.Printf("Failed to update product %d with compressed images: %v", msg.ProductID, err)
        return err
    }

    ip.logger.Printf("Successfully updated product %d with %d compressed images", msg.ProductID, len(compressedUrls))
    return nil
}