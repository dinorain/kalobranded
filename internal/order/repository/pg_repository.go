package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/dinorain/checkoutaja/internal/models"
	"github.com/dinorain/checkoutaja/internal/order"
	"github.com/dinorain/checkoutaja/pkg/utils"
)

// Order repository
type OrderRepository struct {
	db *sqlx.DB
}

var _ order.OrderPGRepository = (*OrderRepository)(nil)

// Order repository constructor
func NewOrderPGRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create new order
func (r *OrderRepository) Create(ctx context.Context, order *models.Order) (*models.Order, error) {
	createdOrder := &models.Order{}
	if err := r.db.QueryRowxContext(
		ctx,
		createOrderQuery,
		order.UserID,
		order.SellerID,
		order.Item,
		order.Quantity,
		order.TotalPrice,
		order.Status,
		order.DeliverySourceAddress,
		order.DeliveryDestinationAddress,
	).StructScan(createdOrder); err != nil {
		return nil, errors.Wrap(err, "OrderPGRepository.Create.QueryRowxContext")
	}

	return createdOrder, nil
}

// UpdateById update existing order
func (r *OrderRepository) UpdateById(ctx context.Context, order *models.Order) (*models.Order, error) {
	if res, err := r.db.ExecContext(
		ctx,
		updateByIdQuery,
		order.OrderID,
		order.UserID,
		order.SellerID,
		order.Item,
		order.Quantity,
		order.TotalPrice,
		order.Status,
		order.DeliverySourceAddress,
		order.DeliveryDestinationAddress,
	); err != nil {
		return nil, errors.Wrap(err, "OrderPGRepository.Update.ExecContext")
	} else {
		_, err := res.RowsAffected()
		if err != nil {
			return nil, errors.Wrap(err, "OrderPGRepository.Update.RowsAffected")
		}
	}

	return order, nil
}

// FindAll Find orders
func (r *OrderRepository) FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.SelectContext(ctx, &orders, findAllQuery, pagination.GetLimit(), pagination.GetOffset()); err != nil {
		return nil, errors.Wrap(err, "OrderPGRepository.FindById.SelectContext")
	}

	return orders, nil
}

// FindAllByUserId Find orders by user uuid
func (r *OrderRepository) FindAllByUserId(ctx context.Context, userID uuid.UUID, pagination *utils.Pagination) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.SelectContext(ctx, &orders, findByUserIdQuery, userID, pagination.GetLimit(), pagination.GetOffset()); err != nil {
		return nil, errors.Wrap(err, "OrderPGRepository.FindAllByUserId.SelectContext")
	}

	return orders, nil
}

// FindAllBySellerId Find orders by seller uuid
func (r *OrderRepository) FindAllBySellerId(ctx context.Context, sellerID uuid.UUID, pagination *utils.Pagination) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.SelectContext(ctx, &orders, findAllBySellerIdQuery, sellerID, pagination.GetLimit(), pagination.GetOffset()); err != nil {
		return nil, errors.Wrap(err, "OrderPGRepository.FindAllBySellerId.SelectContext")
	}

	return orders, nil
}

// FindAllByUserIdSellerId Find orders by user uuid and seller uuid
func (r *OrderRepository) FindAllByUserIdSellerId(ctx context.Context, userID uuid.UUID, sellerID uuid.UUID, pagination *utils.Pagination) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.SelectContext(ctx, &orders, findAllByUserIdSellerIDQuery, userID, sellerID, pagination.GetLimit(), pagination.GetOffset()); err != nil {
		return nil, errors.Wrap(err, "OrderPGRepository.FindAllByUserIdSellerId.SelectContext")
	}

	return orders, nil
}

// FindById Find order by uuid
func (r *OrderRepository) FindById(ctx context.Context, orderID uuid.UUID) (*models.Order, error) {
	order := &models.Order{}
	if err := r.db.GetContext(ctx, order, findByIdQuery, orderID); err != nil {
		return nil, errors.Wrap(err, "OrderPGRepository.FindById.GetContext")
	}

	return order, nil
}

// DeleteById Find order by uuid
func (r *OrderRepository) DeleteById(ctx context.Context, orderID uuid.UUID) error {
	if res, err := r.db.ExecContext(ctx, deleteByIdQuery, orderID); err != nil {
		return errors.Wrap(err, "OrderPGRepository.DeleteById.ExecContext")
	} else {
		cnt, err := res.RowsAffected()
		if err != nil {
			return errors.Wrap(err, "OrderPGRepository.DeleteById.RowsAffected")
		} else if cnt == 0 {
			return sql.ErrNoRows
		}
	}

	return nil
}
