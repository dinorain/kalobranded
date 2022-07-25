package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/dinorain/checkoutaja/internal/models"
	"github.com/dinorain/checkoutaja/internal/seller"
	"github.com/dinorain/checkoutaja/pkg/logger"
)

// Seller redis repository
type sellerRedisRepo struct {
	redisClient *redis.Client
	basePrefix  string
	logger      logger.Logger
}

var _ seller.SellerRedisRepository = (*sellerRedisRepo)(nil)

// Seller redis repository constructor
func NewSellerRedisRepo(redisClient *redis.Client, logger logger.Logger) *sellerRedisRepo {
	return &sellerRedisRepo{redisClient: redisClient, basePrefix: "seller:", logger: logger}
}

// Get seller by id
func (r *sellerRedisRepo) GetByIdCtx(ctx context.Context, key string) (*models.Seller, error) {
	sellerBytes, err := r.redisClient.Get(ctx, r.createKey(key)).Bytes()
	if err != nil {
		return nil, err
	}
	seller := &models.Seller{}
	if err = json.Unmarshal(sellerBytes, seller); err != nil {
		return nil, err
	}

	return seller, nil
}

// Cache seller with duration in seconds
func (r *sellerRedisRepo) SetSellerCtx(ctx context.Context, key string, seconds int, seller *models.Seller) error {
	sellerBytes, err := json.Marshal(seller)
	if err != nil {
		return err
	}

	return r.redisClient.Set(ctx, r.createKey(key), sellerBytes, time.Second*time.Duration(seconds)).Err()
}

// Delete seller by key
func (r *sellerRedisRepo) DeleteSellerCtx(ctx context.Context, key string) error {
	return r.redisClient.Del(ctx, r.createKey(key)).Err()
}

func (r *sellerRedisRepo) createKey(value string) string {
	return fmt.Sprintf("%s: %s", r.basePrefix, value)
}
