//go:generate mockgen -source pg_repository.go -destination mock/pg_repository.go -package mock
package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/dinorain/kalobranded/internal/models"
	"github.com/dinorain/kalobranded/pkg/utils"
)

// Order pg repository
type OrderPGRepository interface {
	Create(ctx context.Context, user *models.Order) (*models.Order, error)
	FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Order, error)
	FindAllByUserId(ctx context.Context, userID uuid.UUID, pagination *utils.Pagination) ([]models.Order, error)
	FindAllByBrandId(ctx context.Context, brandID uuid.UUID, pagination *utils.Pagination) ([]models.Order, error)
	FindAllByUserIdBrandId(ctx context.Context, userID uuid.UUID, brandID uuid.UUID, pagination *utils.Pagination) ([]models.Order, error)
	FindById(ctx context.Context, userID uuid.UUID) (*models.Order, error)
	UpdateById(ctx context.Context, user *models.Order) (*models.Order, error)
	DeleteById(ctx context.Context, userID uuid.UUID) error
}
