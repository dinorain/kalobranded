//go:generate mockgen -source redis_repository.go -destination mock/redis_repository.go -package mock
package brand

import (
	"context"

	"github.com/dinorain/kalobranded/internal/models"
)

// Brand Redis repository interface
type BrandRedisRepository interface {
	GetByIdCtx(ctx context.Context, key string) (*models.Brand, error)
	SetBrandCtx(ctx context.Context, key string, seconds int, user *models.Brand) error
	DeleteBrandCtx(ctx context.Context, key string) error
}
