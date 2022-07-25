package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/dinorain/checkoutaja/internal/models"
	"github.com/dinorain/checkoutaja/internal/order"
	"github.com/dinorain/checkoutaja/pkg/logger"
)

// Order redis repository
type orderRedisRepo struct {
	redisClient *redis.Client
	basePrefix  string
	logger      logger.Logger
}

var _ order.OrderRedisRepository = (*orderRedisRepo)(nil)

// Order redis repository constructor
func NewOrderRedisRepo(redisClient *redis.Client, logger logger.Logger) *orderRedisRepo {
	return &orderRedisRepo{redisClient: redisClient, basePrefix: "order:", logger: logger}
}

// Get order by id
func (r *orderRedisRepo) GetByIdCtx(ctx context.Context, key string) (*models.Order, error) {
	orderBytes, err := r.redisClient.Get(ctx, r.createKey(key)).Bytes()
	if err != nil {
		return nil, err
	}
	order := &models.Order{}
	if err = json.Unmarshal(orderBytes, order); err != nil {
		return nil, err
	}

	return order, nil
}

// Cache order with duration in seconds
func (r *orderRedisRepo) SetOrderCtx(ctx context.Context, key string, seconds int, order *models.Order) error {
	orderBytes, err := json.Marshal(order)
	if err != nil {
		return err
	}

	return r.redisClient.Set(ctx, r.createKey(key), orderBytes, time.Second*time.Duration(seconds)).Err()
}

// Delete order by key
func (r *orderRedisRepo) DeleteOrderCtx(ctx context.Context, key string) error {
	return r.redisClient.Del(ctx, r.createKey(key)).Err()
}

func (r *orderRedisRepo) createKey(value string) string {
	return fmt.Sprintf("%s: %s", r.basePrefix, value)
}
