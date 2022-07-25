package usecase

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/dinorain/kalobranded/config"
	"github.com/dinorain/kalobranded/internal/brand"
	"github.com/dinorain/kalobranded/internal/models"
	"github.com/dinorain/kalobranded/pkg/logger"
	"github.com/dinorain/kalobranded/pkg/utils"
)

const (
	brandByIdCacheDuration = 3600
)

// Brand UseCase
type brandUseCase struct {
	cfg         *config.Config
	logger      logger.Logger
	brandPgRepo brand.BrandPGRepository
	redisRepo   brand.BrandRedisRepository
}

var _ brand.BrandUseCase = (*brandUseCase)(nil)

// New Brand UseCase
func NewBrandUseCase(cfg *config.Config, logger logger.Logger, brandRepo brand.BrandPGRepository, redisRepo brand.BrandRedisRepository) *brandUseCase {
	return &brandUseCase{cfg: cfg, logger: logger, brandPgRepo: brandRepo, redisRepo: redisRepo}
}

// Register new brand
func (u *brandUseCase) Register(ctx context.Context, brand *models.Brand) (*models.Brand, error) {
	return u.brandPgRepo.Create(ctx, brand)
}

// FindAll find brands
func (u *brandUseCase) FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Brand, error) {
	brands, err := u.brandPgRepo.FindAll(ctx, pagination)
	if err != nil {
		return nil, errors.Wrap(err, "brandPgRepo.FindAll")
	}

	return brands, nil
}

// FindById find brand by uuid
func (u *brandUseCase) FindById(ctx context.Context, brandID uuid.UUID) (*models.Brand, error) {
	foundBrand, err := u.brandPgRepo.FindById(ctx, brandID)
	if err != nil {
		return nil, errors.Wrap(err, "brandPgRepo.FindById")
	}

	return foundBrand, nil
}

// CachedFindById find brand by uuid from cache
func (u *brandUseCase) CachedFindById(ctx context.Context, brandID uuid.UUID) (*models.Brand, error) {
	cachedBrand, err := u.redisRepo.GetByIdCtx(ctx, brandID.String())
	if err != nil && !errors.Is(err, redis.Nil) {
		u.logger.Errorf("redisRepo.GetByIdCtx", err)
	}
	if cachedBrand != nil {
		return cachedBrand, nil
	}

	foundBrand, err := u.brandPgRepo.FindById(ctx, brandID)
	if err != nil {
		return nil, errors.Wrap(err, "brandPgRepo.FindById")
	}

	if err := u.redisRepo.SetBrandCtx(ctx, foundBrand.BrandID.String(), brandByIdCacheDuration, foundBrand); err != nil {
		u.logger.Errorf("redisRepo.SetBrandCtx", err)
	}

	return foundBrand, nil
}

// UpdateById update brand by uuid
func (u *brandUseCase) UpdateById(ctx context.Context, brand *models.Brand) (*models.Brand, error) {
	updatedBrand, err := u.brandPgRepo.UpdateById(ctx, brand)
	if err != nil {
		return nil, errors.Wrap(err, "brandPgRepo.UpdateById")
	}

	if err := u.redisRepo.SetBrandCtx(ctx, updatedBrand.BrandID.String(), brandByIdCacheDuration, updatedBrand); err != nil {
		u.logger.Errorf("redisRepo.SetBrandCtx", err)
	}

	return updatedBrand, nil
}

// DeleteById delete brand by uuid
func (u *brandUseCase) DeleteById(ctx context.Context, brandID uuid.UUID) error {
	err := u.brandPgRepo.DeleteById(ctx, brandID)
	if err != nil {
		return errors.Wrap(err, "brandPgRepo.DeleteById")
	}

	if err := u.redisRepo.DeleteBrandCtx(ctx, brandID.String()); err != nil {
		u.logger.Errorf("redisRepo.DeleteBrandCtx", err)
	}

	return nil
}
