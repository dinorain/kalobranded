package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/dinorain/checkoutaja/internal/models"
	"github.com/dinorain/checkoutaja/internal/product"
	"github.com/dinorain/checkoutaja/pkg/logger"
)

// Product redis repository
type productRedisRepo struct {
	redisClient *redis.Client
	basePrefix  string
	logger      logger.Logger
}

var _ product.ProductRedisRepository = (*productRedisRepo)(nil)

// Product redis repository constructor
func NewProductRedisRepo(redisClient *redis.Client, logger logger.Logger) *productRedisRepo {
	return &productRedisRepo{redisClient: redisClient, basePrefix: "product:", logger: logger}
}

// Get product by id
func (r *productRedisRepo) GetByIdCtx(ctx context.Context, key string) (*models.Product, error) {
	productBytes, err := r.redisClient.Get(ctx, r.createKey(key)).Bytes()
	if err != nil {
		return nil, err
	}
	product := &models.Product{}
	if err = json.Unmarshal(productBytes, product); err != nil {
		return nil, err
	}

	return product, nil
}

// Cache product with duration in seconds
func (r *productRedisRepo) SetProductCtx(ctx context.Context, key string, seconds int, product *models.Product) error {
	productBytes, err := json.Marshal(product)
	if err != nil {
		return err
	}

	return r.redisClient.Set(ctx, r.createKey(key), productBytes, time.Second*time.Duration(seconds)).Err()
}

// Delete product by key
func (r *productRedisRepo) DeleteProductCtx(ctx context.Context, key string) error {
	return r.redisClient.Del(ctx, r.createKey(key)).Err()
}

func (r *productRedisRepo) createKey(value string) string {
	return fmt.Sprintf("%s: %s", r.basePrefix, value)
}
