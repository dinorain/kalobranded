//go:generate mockgen -source pg_repository.go -destination mock/pg_repository.go -package mock
package seller

import (
	"context"

	"github.com/google/uuid"

	"github.com/dinorain/checkoutaja/internal/models"
	"github.com/dinorain/checkoutaja/pkg/utils"
)

// Seller pg repository
type SellerPGRepository interface {
	Create(ctx context.Context, user *models.Seller) (*models.Seller, error)
	FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Seller, error)
	FindByEmail(ctx context.Context, email string) (*models.Seller, error)
	FindById(ctx context.Context, userID uuid.UUID) (*models.Seller, error)
	UpdateById(ctx context.Context, user *models.Seller) (*models.Seller, error)
	DeleteById(ctx context.Context, userID uuid.UUID) error
}
