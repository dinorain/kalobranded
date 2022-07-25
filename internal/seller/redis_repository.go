//go:generate mockgen -source redis_repository.go -destination mock/redis_repository.go -package mock
package seller

import (
	"context"

	"github.com/dinorain/checkoutaja/internal/models"
)

// Seller Redis repository interface
type SellerRedisRepository interface {
	GetByIdCtx(ctx context.Context, key string) (*models.Seller, error)
	SetSellerCtx(ctx context.Context, key string, seconds int, user *models.Seller) error
	DeleteSellerCtx(ctx context.Context, key string) error
}
