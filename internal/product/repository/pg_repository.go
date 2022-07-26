package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/dinorain/kalobranded/internal/models"
	"github.com/dinorain/kalobranded/internal/product"
	"github.com/dinorain/kalobranded/pkg/utils"
)

// Product repository
type ProductRepository struct {
	db *sqlx.DB
}

var _ product.ProductPGRepository = (*ProductRepository)(nil)

// Product repository constructor
func NewProductPGRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create new product
func (r *ProductRepository) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	createdProduct := &models.Product{}
	if err := r.db.QueryRowxContext(
		ctx,
		createProductQuery,
		product.Name,
		product.Description,
		product.Price,
		product.BrandID,
	).StructScan(createdProduct); err != nil {
		return nil, errors.Wrap(err, "ProductRepository.Create.QueryRowxContext")
	}

	return createdProduct, nil
}

// UpdateById update existing product
func (r *ProductRepository) UpdateById(ctx context.Context, product *models.Product) (*models.Product, error) {
	if res, err := r.db.ExecContext(
		ctx,
		updateByIdQuery,
		product.ProductID,
		product.Name,
		product.Description,
		product.Price,
		product.BrandID,
	); err != nil {
		return nil, errors.Wrap(err, "ProductRepository.Update.ExecContext")
	} else {
		_, err := res.RowsAffected()
		if err != nil {
			return nil, errors.Wrap(err, "Update.RowsAffected")
		}
	}

	return product, nil
}

// FindAll Find products
func (r *ProductRepository) FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Product, error) {
	var products []models.Product
	if err := r.db.SelectContext(ctx, &products, findAllQuery, pagination.GetLimit(), pagination.GetOffset()); err != nil {
		return nil, errors.Wrap(err, "ProductRepository.FindById.SelectContext")
	}

	return products, nil
}

// FindAllByBrandId Find products by brand uuid
func (r *ProductRepository) FindAllByBrandId(ctx context.Context, brandID uuid.UUID, pagination *utils.Pagination) ([]models.Product, error) {
	var products []models.Product
	if err := r.db.SelectContext(ctx, &products, findAllByBrandIdQuery, brandID, pagination.GetLimit(), pagination.GetOffset()); err != nil {
		return nil, errors.Wrap(err, "ProductPGRepository.FindAllByBrandId.SelectContext")
	}

	return products, nil
}

// FindById Find product by uuid
func (r *ProductRepository) FindById(ctx context.Context, productID uuid.UUID) (*models.Product, error) {
	product := &models.Product{}
	if err := r.db.GetContext(ctx, product, findByIdQuery, productID); err != nil {
		return nil, errors.Wrap(err, "ProductRepository.FindById.GetContext")
	}

	return product, nil
}

// DeleteById Find product by uuid
func (r *ProductRepository) DeleteById(ctx context.Context, productID uuid.UUID) error {
	if res, err := r.db.ExecContext(ctx, deleteByIdQuery, productID); err != nil {
		return errors.Wrap(err, "ProductRepository.DeleteById.ExecContext")
	} else {
		cnt, err := res.RowsAffected()
		if err != nil {
			return errors.Wrap(err, "ProductRepository.DeleteById.RowsAffected")
		} else if cnt == 0 {
			return sql.ErrNoRows
		}
	}

	return nil
}
