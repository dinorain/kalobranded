//go:generate mockgen -source redis_repository.go -destination mock/redis_repository.go -package mock
package product

import (
	"context"

	"github.com/dinorain/checkoutaja/internal/models"
)

// Product Redis repository interface
type ProductRedisRepository interface {
	GetByIdCtx(ctx context.Context, key string) (*models.Product, error)
	SetProductCtx(ctx context.Context, key string, seconds int, user *models.Product) error
	DeleteProductCtx(ctx context.Context, key string) error
}
