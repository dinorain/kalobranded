//go:generate mockgen -source pg_repository.go -destination mock/pg_repository.go -package mock
package brand

import (
	"context"

	"github.com/google/uuid"

	"github.com/dinorain/kalobranded/internal/models"
	"github.com/dinorain/kalobranded/pkg/utils"
)

// Brand pg repository
type BrandPGRepository interface {
	Create(ctx context.Context, user *models.Brand) (*models.Brand, error)
	FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Brand, error)
	FindById(ctx context.Context, userID uuid.UUID) (*models.Brand, error)
	UpdateById(ctx context.Context, user *models.Brand) (*models.Brand, error)
	DeleteById(ctx context.Context, userID uuid.UUID) error
}
