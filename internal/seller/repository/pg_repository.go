package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/dinorain/checkoutaja/internal/models"
	"github.com/dinorain/checkoutaja/internal/seller"
	"github.com/dinorain/checkoutaja/pkg/utils"
)

// Seller repository
type SellerRepository struct {
	db *sqlx.DB
}

var _ seller.SellerPGRepository = (*SellerRepository)(nil)

// Seller repository constructor
func NewSellerPGRepository(db *sqlx.DB) *SellerRepository {
	return &SellerRepository{db: db}
}

// Create new seller
func (r *SellerRepository) Create(ctx context.Context, seller *models.Seller) (*models.Seller, error) {
	createdSeller := &models.Seller{}
	if err := r.db.QueryRowxContext(
		ctx,
		createSellerQuery,
		seller.FirstName,
		seller.LastName,
		seller.Email,
		seller.Password,
		seller.Avatar,
		seller.PickupAddress,
	).StructScan(createdSeller); err != nil {
		return nil, errors.Wrap(err, "SellerRepository.Create.QueryRowxContext")
	}

	return createdSeller, nil
}

// UpdateById update existing seller
func (r *SellerRepository) UpdateById(ctx context.Context, seller *models.Seller) (*models.Seller, error) {
	if res, err := r.db.ExecContext(
		ctx,
		updateByIdQuery,
		seller.SellerID,
		seller.FirstName,
		seller.LastName,
		seller.Email,
		seller.Password,
		seller.Avatar,
		seller.PickupAddress,
	); err != nil {
		return nil, errors.Wrap(err, "UpdateById.Update.ExecContext")
	} else {
		_, err := res.RowsAffected()
		if err != nil {
			return nil, errors.Wrap(err, "UpdateById.Update.RowsAffected")
		}
	}

	return seller, nil
}

// FindAll Find sellers
func (r *SellerRepository) FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Seller, error) {
	var sellers []models.Seller
	if err := r.db.SelectContext(ctx, &sellers, findAllQuery, pagination.GetLimit(), pagination.GetOffset()); err != nil {
		return nil, errors.Wrap(err, "SellerRepository.FindById.SelectContext")
	}

	return sellers, nil
}

// FindByEmail Find by seller email address
func (r *SellerRepository) FindByEmail(ctx context.Context, email string) (*models.Seller, error) {
	seller := &models.Seller{}
	if err := r.db.GetContext(ctx, seller, findByEmailQuery, email); err != nil {
		return nil, errors.Wrap(err, "FindByEmail.GetContext")
	}

	return seller, nil
}

// FindById Find seller by uuid
func (r *SellerRepository) FindById(ctx context.Context, sellerID uuid.UUID) (*models.Seller, error) {
	seller := &models.Seller{}
	if err := r.db.GetContext(ctx, seller, findByIdQuery, sellerID); err != nil {
		return nil, errors.Wrap(err, "SellerRepository.FindById.GetContext")
	}

	return seller, nil
}

// DeleteById Find seller by uuid
func (r *SellerRepository) DeleteById(ctx context.Context, sellerID uuid.UUID) error {
	if res, err := r.db.ExecContext(ctx, deleteByIdQuery, sellerID); err != nil {
		return errors.Wrap(err, "SellerRepository.DeleteById.ExecContext")
	} else {
		cnt, err := res.RowsAffected()
		if err != nil {
			return errors.Wrap(err, "SellerRepository.DeleteById.RowsAffected")
		} else if cnt == 0 {
			return sql.ErrNoRows
		}
	}

	return nil
}
