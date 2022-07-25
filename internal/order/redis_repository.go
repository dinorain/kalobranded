//go:generate mockgen -source redis_repository.go -destination mock/redis_repository.go -package mock
package order

import (
	"context"

	"github.com/dinorain/checkoutaja/internal/models"
)

// Order Redis repository interface
type OrderRedisRepository interface {
	GetByIdCtx(ctx context.Context, key string) (*models.Order, error)
	SetOrderCtx(ctx context.Context, key string, seconds int, user *models.Order) error
	DeleteOrderCtx(ctx context.Context, key string) error
}
