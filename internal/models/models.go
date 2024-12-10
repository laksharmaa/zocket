package models

import (
    "time"
    // "gorm.io/gorm"
)

type User struct {
    ID        uint      `gorm:"primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    Name      string    `gorm:"not null"`
    Email     string    `gorm:"unique;not null"`
}

type Product struct {
    ID                     uint      `gorm:"primaryKey"`
    CreatedAt             time.Time
    UpdatedAt             time.Time
    UserID                uint      `gorm:"not null"`
    ProductName           string    `gorm:"not null"`
    ProductDescription    string    `gorm:"type:text"`
    ProductImages         []string  `gorm:"type:text[]"`
    CompressedProductImages []string `gorm:"type:text[]"`
    ProductPrice          float64   `gorm:"type:decimal(10,2);not null"`
}