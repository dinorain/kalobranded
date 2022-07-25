//go:generate mockgen -source usecase.go -destination mock/usecase.go -package mock
package brand

import (
	"context"

	"github.com/google/uuid"

	"github.com/dinorain/kalobranded/internal/models"
	"github.com/dinorain/kalobranded/pkg/utils"
)

//  Brand UseCase interface
type BrandUseCase interface {
	Register(ctx context.Context, brand *models.Brand) (*models.Brand, error)
	FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Brand, error)
	FindById(ctx context.Context, brandID uuid.UUID) (*models.Brand, error)
	CachedFindById(ctx context.Context, brandID uuid.UUID) (*models.Brand, error)
	UpdateById(ctx context.Context, brand *models.Brand) (*models.Brand, error)
	DeleteById(ctx context.Context, brandID uuid.UUID) error
}
