package database

import (
    "fmt"
    "product-management-system/internal/config"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func NewPostgresDB(config *config.Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        config.DBHost,
        config.DBUser,
        config.DBPassword,
        config.DBName,
        config.DBPort,
    )

    // Initialize the database connection
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to the database: %w", err)
    }

    // Test the database connection
    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("failed to get DB from gorm: %w", err)
    }

    if err := sqlDB.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping the database: %w", err)
    }

    return db, nil
}