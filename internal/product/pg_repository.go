//go:generate mockgen -source pg_repository.go -destination mock/pg_repository.go -package mock
package product

import (
	"context"

	"github.com/google/uuid"

	"github.com/dinorain/kalobranded/internal/models"
	"github.com/dinorain/kalobranded/pkg/utils"
)

// Product pg repository
type ProductPGRepository interface {
	Create(ctx context.Context, user *models.Product) (*models.Product, error)
	FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Product, error)
	FindAllByBrandId(ctx context.Context, brandID uuid.UUID, pagination *utils.Pagination) ([]models.Product, error)
	FindById(ctx context.Context, userID uuid.UUID) (*models.Product, error)
	UpdateById(ctx context.Context, user *models.Product) (*models.Product, error)
	DeleteById(ctx context.Context, userID uuid.UUID) error
}
