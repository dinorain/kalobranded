package usecase

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/dinorain/checkoutaja/config"
	"github.com/dinorain/checkoutaja/internal/models"
	"github.com/dinorain/checkoutaja/internal/product"
	"github.com/dinorain/checkoutaja/pkg/logger"
	"github.com/dinorain/checkoutaja/pkg/utils"
)

const (
	productByIdCacheDuration = 3600
)

// Product UseCase
type productUseCase struct {
	cfg           *config.Config
	logger        logger.Logger
	productPgRepo product.ProductPGRepository
	redisRepo     product.ProductRedisRepository
}

var _ product.ProductUseCase = (*productUseCase)(nil)

// New Product UseCase
func NewProductUseCase(cfg *config.Config, logger logger.Logger, productRepo product.ProductPGRepository, redisRepo product.ProductRedisRepository) *productUseCase {
	return &productUseCase{cfg: cfg, logger: logger, productPgRepo: productRepo, redisRepo: redisRepo}
}

// Create new product
func (u *productUseCase) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	return u.productPgRepo.Create(ctx, product)
}

// FindAll find products
func (u *productUseCase) FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Product, error) {
	products, err := u.productPgRepo.FindAll(ctx, pagination)
	if err != nil {
		return nil, errors.Wrap(err, "productPgRepo.FindAll")
	}

	return products, nil
}

// FindAllBySellerId find products by seller id
func (u *productUseCase) FindAllBySellerId(ctx context.Context, sellerID uuid.UUID, pagination *utils.Pagination) ([]models.Product, error) {
	products, err := u.productPgRepo.FindAllBySellerId(ctx, sellerID, pagination)
	if err != nil {
		return nil, errors.Wrap(err, "productPgRepo.FindAllBySellerId")
	}

	return products, nil
}

// FindById find product by uuid
func (u *productUseCase) FindById(ctx context.Context, productID uuid.UUID) (*models.Product, error) {
	foundProduct, err := u.productPgRepo.FindById(ctx, productID)
	if err != nil {
		return nil, errors.Wrap(err, "productPgRepo.FindById")
	}

	return foundProduct, nil
}

// CachedFindById find product by uuid from cache
func (u *productUseCase) CachedFindById(ctx context.Context, productID uuid.UUID) (*models.Product, error) {
	cachedProduct, err := u.redisRepo.GetByIdCtx(ctx, productID.String())
	if err != nil && !errors.Is(err, redis.Nil) {
		u.logger.Errorf("redisRepo.GetByIdCtx", err)
	}
	if cachedProduct != nil {
		return cachedProduct, nil
	}

	foundProduct, err := u.productPgRepo.FindById(ctx, productID)
	if err != nil {
		return nil, errors.Wrap(err, "productPgRepo.FindById")
	}

	if err := u.redisRepo.SetProductCtx(ctx, foundProduct.ProductID.String(), productByIdCacheDuration, foundProduct); err != nil {
		u.logger.Errorf("redisRepo.SetProductCtx", err)
	}

	return foundProduct, nil
}

// UpdateById update product by uuid
func (u *productUseCase) UpdateById(ctx context.Context, product *models.Product) (*models.Product, error) {
	updatedProduct, err := u.productPgRepo.UpdateById(ctx, product)
	if err != nil {
		return nil, errors.Wrap(err, "productPgRepo.UpdateById")
	}

	if err := u.redisRepo.SetProductCtx(ctx, updatedProduct.ProductID.String(), productByIdCacheDuration, updatedProduct); err != nil {
		u.logger.Errorf("redisRepo.SetProductCtx", err)
	}

	return updatedProduct, nil
}

// DeleteById delete product by uuid
func (u *productUseCase) DeleteById(ctx context.Context, productID uuid.UUID) error {
	err := u.productPgRepo.DeleteById(ctx, productID)
	if err != nil {
		return errors.Wrap(err, "productPgRepo.DeleteById")
	}

	if err := u.redisRepo.DeleteProductCtx(ctx, productID.String()); err != nil {
		u.logger.Errorf("redisRepo.DeleteProductCtx", err)
	}

	return nil
}
