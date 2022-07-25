package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/dinorain/kalobranded/internal/brand"
	"github.com/dinorain/kalobranded/internal/models"
	"github.com/dinorain/kalobranded/pkg/utils"
)

// Brand repository
type BrandRepository struct {
	db *sqlx.DB
}

var _ brand.BrandPGRepository = (*BrandRepository)(nil)

// Brand repository constructor
func NewBrandPGRepository(db *sqlx.DB) *BrandRepository {
	return &BrandRepository{db: db}
}

// Create new brand
func (r *BrandRepository) Create(ctx context.Context, brand *models.Brand) (*models.Brand, error) {
	createdBrand := &models.Brand{}
	if err := r.db.QueryRowxContext(
		ctx,
		createBrandQuery,
		brand.BrandName,
		brand.Logo,
		brand.PickupAddress,
	).StructScan(createdBrand); err != nil {
		return nil, errors.Wrap(err, "BrandRepository.Create.QueryRowxContext")
	}

	return createdBrand, nil
}

// UpdateById update existing brand
func (r *BrandRepository) UpdateById(ctx context.Context, brand *models.Brand) (*models.Brand, error) {
	if res, err := r.db.ExecContext(
		ctx,
		updateByIdQuery,
		brand.BrandID,
		brand.BrandName,
		brand.Logo,
		brand.PickupAddress,
	); err != nil {
		return nil, errors.Wrap(err, "UpdateById.Update.ExecContext")
	} else {
		_, err := res.RowsAffected()
		if err != nil {
			return nil, errors.Wrap(err, "UpdateById.Update.RowsAffected")
		}
	}

	return brand, nil
}

// FindAll Find brands
func (r *BrandRepository) FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Brand, error) {
	var brands []models.Brand
	if err := r.db.SelectContext(ctx, &brands, findAllQuery, pagination.GetLimit(), pagination.GetOffset()); err != nil {
		return nil, errors.Wrap(err, "BrandRepository.FindById.SelectContext")
	}

	return brands, nil
}

// FindById Find brand by uuid
func (r *BrandRepository) FindById(ctx context.Context, brandID uuid.UUID) (*models.Brand, error) {
	brand := &models.Brand{}
	if err := r.db.GetContext(ctx, brand, findByIdQuery, brandID); err != nil {
		return nil, errors.Wrap(err, "BrandRepository.FindById.GetContext")
	}

	return brand, nil
}

// DeleteById Find brand by uuid
func (r *BrandRepository) DeleteById(ctx context.Context, brandID uuid.UUID) error {
	if res, err := r.db.ExecContext(ctx, deleteByIdQuery, brandID); err != nil {
		return errors.Wrap(err, "BrandRepository.DeleteById.ExecContext")
	} else {
		cnt, err := res.RowsAffected()
		if err != nil {
			return errors.Wrap(err, "BrandRepository.DeleteById.RowsAffected")
		} else if cnt == 0 {
			return sql.ErrNoRows
		}
	}

	return nil
}
