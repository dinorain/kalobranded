package usecase

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/dinorain/checkoutaja/config"
	"github.com/dinorain/checkoutaja/internal/models"
	"github.com/dinorain/checkoutaja/internal/order"
	"github.com/dinorain/checkoutaja/pkg/logger"
	"github.com/dinorain/checkoutaja/pkg/utils"
)

const (
	orderByIdCacheDuration = 3600
)

// Order UseCase
type orderUseCase struct {
	cfg         *config.Config
	logger      logger.Logger
	orderPgRepo order.OrderPGRepository
	redisRepo   order.OrderRedisRepository
}

var _ order.OrderUseCase = (*orderUseCase)(nil)

// New Order UseCase
func NewOrderUseCase(cfg *config.Config, logger logger.Logger, orderRepo order.OrderPGRepository, redisRepo order.OrderRedisRepository) *orderUseCase {
	return &orderUseCase{cfg: cfg, logger: logger, orderPgRepo: orderRepo, redisRepo: redisRepo}
}

// Create new order
func (u *orderUseCase) Create(ctx context.Context, order *models.Order) (*models.Order, error) {
	return u.orderPgRepo.Create(ctx, order)
}

// FindAll find orders
func (u *orderUseCase) FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Order, error) {
	orders, err := u.orderPgRepo.FindAll(ctx, pagination)
	if err != nil {
		return nil, errors.Wrap(err, "orderPgRepo.FindAll")
	}

	return orders, nil
}

// FindAllByUserId find orders by user id
func (u *orderUseCase) FindAllByUserId(ctx context.Context, userID uuid.UUID, pagination *utils.Pagination) ([]models.Order, error) {
	orders, err := u.orderPgRepo.FindAllByUserId(ctx, userID, pagination)
	if err != nil {
		return nil, errors.Wrap(err, "orderPgRepo.FindAllByUserId")
	}

	return orders, nil
}

// FindAllBySellerId find orders by seller id
func (u *orderUseCase) FindAllBySellerId(ctx context.Context, sellerID uuid.UUID, pagination *utils.Pagination) ([]models.Order, error) {
	orders, err := u.orderPgRepo.FindAllBySellerId(ctx, sellerID, pagination)
	if err != nil {
		return nil, errors.Wrap(err, "orderPgRepo.FindAllByUserId")
	}

	return orders, nil
}

// FindAllByUserIdSellerId find orders by seller id
func (u *orderUseCase) FindAllByUserIdSellerId(ctx context.Context, userID uuid.UUID, sellerID uuid.UUID, pagination *utils.Pagination) ([]models.Order, error) {
	orders, err := u.orderPgRepo.FindAllByUserIdSellerId(ctx, userID, sellerID, pagination)
	if err != nil {
		return nil, errors.Wrap(err, "orderPgRepo.FindAllByUserIdSellerId")
	}

	return orders, nil
}

// FindById find order by uuid
func (u *orderUseCase) FindById(ctx context.Context, orderID uuid.UUID) (*models.Order, error) {
	foundOrder, err := u.orderPgRepo.FindById(ctx, orderID)
	if err != nil {
		return nil, errors.Wrap(err, "orderPgRepo.FindById")
	}

	return foundOrder, nil
}

// CachedFindById find order by uuid from cache
func (u *orderUseCase) CachedFindById(ctx context.Context, orderID uuid.UUID) (*models.Order, error) {
	cachedOrder, err := u.redisRepo.GetByIdCtx(ctx, orderID.String())
	if err != nil && !errors.Is(err, redis.Nil) {
		u.logger.Errorf("redisRepo.GetByIdCtx", err)
	}
	if cachedOrder != nil {
		return cachedOrder, nil
	}

	foundOrder, err := u.orderPgRepo.FindById(ctx, orderID)
	if err != nil {
		return nil, errors.Wrap(err, "orderPgRepo.FindById")
	}

	if err := u.redisRepo.SetOrderCtx(ctx, foundOrder.OrderID.String(), orderByIdCacheDuration, foundOrder); err != nil {
		u.logger.Errorf("redisRepo.SetOrderCtx", err)
	}

	return foundOrder, nil
}

// UpdateById update order by uuid
func (u *orderUseCase) UpdateById(ctx context.Context, order *models.Order) (*models.Order, error) {
	updatedOrder, err := u.orderPgRepo.UpdateById(ctx, order)
	if err != nil {
		return nil, errors.Wrap(err, "orderPgRepo.UpdateById")
	}

	if err := u.redisRepo.SetOrderCtx(ctx, updatedOrder.OrderID.String(), orderByIdCacheDuration, updatedOrder); err != nil {
		u.logger.Errorf("redisRepo.SetOrderCtx", err)
	}

	return updatedOrder, nil
}

// DeleteById delete order by uuid
func (u *orderUseCase) DeleteById(ctx context.Context, orderID uuid.UUID) error {
	err := u.orderPgRepo.DeleteById(ctx, orderID)
	if err != nil {
		return errors.Wrap(err, "orderPgRepo.DeleteById")
	}

	if err := u.redisRepo.DeleteOrderCtx(ctx, orderID.String()); err != nil {
		u.logger.Errorf("redisRepo.DeleteOrderCtx", err)
	}

	return nil
}
