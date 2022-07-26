package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/dinorain/kalobranded/internal/brand"
	"github.com/dinorain/kalobranded/internal/models"
	"github.com/dinorain/kalobranded/pkg/logger"
)

// Brand redis repository
type brandRedisRepo struct {
	redisClient *redis.Client
	basePrefix  string
	logger      logger.Logger
}

var _ brand.BrandRedisRepository = (*brandRedisRepo)(nil)

// Brand redis repository constructor
func NewBrandRedisRepo(redisClient *redis.Client, logger logger.Logger) *brandRedisRepo {
	return &brandRedisRepo{redisClient: redisClient, basePrefix: "brand:", logger: logger}
}

// Get brand by id
func (r *brandRedisRepo) GetByIdCtx(ctx context.Context, key string) (*models.Brand, error) {
	brandBytes, err := r.redisClient.Get(ctx, r.createKey(key)).Bytes()
	if err != nil {
		return nil, err
	}
	brand := &models.Brand{}
	if err = json.Unmarshal(brandBytes, brand); err != nil {
		return nil, err
	}

	return brand, nil
}

// Cache brand with duration in seconds
func (r *brandRedisRepo) SetBrandCtx(ctx context.Context, key string, seconds int, brand *models.Brand) error {
	brandBytes, err := json.Marshal(brand)
	if err != nil {
		return err
	}

	return r.redisClient.Set(ctx, r.createKey(key), brandBytes, time.Second*time.Duration(seconds)).Err()
}

// Delete brand by key
func (r *brandRedisRepo) DeleteBrandCtx(ctx context.Context, key string) error {
	return r.redisClient.Del(ctx, r.createKey(key)).Err()
}

func (r *brandRedisRepo) createKey(value string) string {
	return fmt.Sprintf("%s: %s", r.basePrefix, value)
}
