package config

import (
    "fmt"
    "github.com/spf13/viper"
)

type Config struct {
    DBHost           string
    DBPort           string
    DBUser           string
    DBPassword       string
    DBName           string
    CloudinaryName   string
    CloudinaryKey    string
    CloudinarySecret string
    RedisURL         string
    RabbitMQURL      string
}

// LoadConfig initializes configuration from .env or environment variables
func LoadConfig() (*Config, error) {
    // Set default file name and enable automatic environment variable overrides
    viper.SetConfigFile(".env")
    viper.AutomaticEnv()

    // Set defaults for optional configuration values
    viper.SetDefault("DB_PORT", "5432")
    viper.SetDefault("REDIS_URL", "redis://localhost:6379")
    viper.SetDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672")

    // Read configuration file
    if err := viper.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("failed to read configuration: %w", err)
    }

    // Load values into the Config struct
    config := &Config{
        DBHost:           viper.GetString("DB_HOST"),
        DBPort:           viper.GetString("DB_PORT"),
        DBUser:           viper.GetString("DB_USER"),
        DBPassword:       viper.GetString("DB_PASSWORD"),
        DBName:           viper.GetString("DB_NAME"),
        CloudinaryName:   viper.GetString("CLOUDINARY_CLOUD_NAME"),
        CloudinaryKey:    viper.GetString("CLOUDINARY_API_KEY"),
        CloudinarySecret: viper.GetString("CLOUDINARY_API_SECRET"),
        RedisURL:         viper.GetString("REDIS_URL"),
        RabbitMQURL:      viper.GetString("RABBITMQ_URL"),
    }

    // Validate required fields
    if config.DBHost == "" || config.DBUser == "" || config.DBPassword == "" || config.DBName == "" {
        return nil, fmt.Errorf("missing required database configuration")
    }

    return config, nil
}
