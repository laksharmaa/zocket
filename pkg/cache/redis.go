package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"product-management-system/internal/models"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
    client *redis.Client
}

func NewRedisCache(url string) (*RedisCache, error) {
    client := redis.NewClient(&redis.Options{
        Addr: url,
    })

    _, err := client.Ping(context.Background()).Result()
    if err != nil {
        return nil, err
    }

    return &RedisCache{
        client: client,
    }, nil
}

func (c *RedisCache) SetProduct(ctx context.Context, product *models.Product) error {
    json, err := json.Marshal(product)
    if err != nil {
        return err
    }

    return c.client.Set(ctx, 
        getProductKey(product.ID), 
        json, 
        time.Hour*24,
    ).Err()
}

func (c *RedisCache) GetProduct(ctx context.Context, id uint) (*models.Product, error) {
    val, err := c.client.Get(ctx, getProductKey(id)).Result()
    if err == redis.Nil {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

    var product models.Product
    if err := json.Unmarshal([]byte(val), &product); err != nil {
        return nil, err
    }

    return &product, nil
}

func (c *RedisCache) InvalidateProduct(ctx context.Context, id uint) error {
    return c.client.Del(ctx, getProductKey(id)).Err()
}

func getProductKey(id uint) string {
    return fmt.Sprintf("product:%d", id)
}