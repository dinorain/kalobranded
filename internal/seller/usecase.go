//go:generate mockgen -source usecase.go -destination mock/usecase.go -package mock
package seller

import (
	"context"

	"github.com/google/uuid"

	"github.com/dinorain/checkoutaja/internal/models"
	"github.com/dinorain/checkoutaja/pkg/utils"
)

//  Seller UseCase interface
type SellerUseCase interface {
	Register(ctx context.Context, seller *models.Seller) (*models.Seller, error)
	Login(ctx context.Context, email string, password string) (*models.Seller, error)
	FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Seller, error)
	FindByEmail(ctx context.Context, email string) (*models.Seller, error)
	FindById(ctx context.Context, sellerID uuid.UUID) (*models.Seller, error)
	CachedFindById(ctx context.Context, sellerID uuid.UUID) (*models.Seller, error)
	UpdateById(ctx context.Context, seller *models.Seller) (*models.Seller, error)
	DeleteById(ctx context.Context, sellerID uuid.UUID) error
	GenerateTokenPair(seller *models.Seller, sessionID string) (access string, refresh string, err error)
}
