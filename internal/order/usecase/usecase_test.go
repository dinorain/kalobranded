package usecase

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/dinorain/checkoutaja/config"
	"github.com/dinorain/checkoutaja/internal/models"
	"github.com/dinorain/checkoutaja/internal/order/mock"
	"github.com/dinorain/checkoutaja/pkg/logger"
)

func TestOrderUseCase_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderPGRepository := mock.NewMockOrderPGRepository(ctrl)
	orderRedisRepository := mock.NewMockOrderRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{Server: config.ServerConfig{JwtSecretKey: "secret123"}}
	orderUC := NewOrderUseCase(cfg, apiLogger, orderPGRepository, orderRedisRepository)

	orderUUID := uuid.New()
	userUUID := uuid.New()
	sellerUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		SellerID: sellerUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			SellerID:    sellerUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	ctx := context.Background()

	orderPGRepository.EXPECT().Create(gomock.Any(), mockOrder).Return(&models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		SellerID: sellerUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			SellerID:    sellerUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}, nil)

	createdOrder, err := orderUC.Create(ctx, mockOrder)
	require.NoError(t, err)
	require.NotNil(t, createdOrder)
	require.Equal(t, createdOrder.OrderID, orderUUID)
}

func TestOrderUseCase_FindAll(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sellerPGRepository := mock.NewMockOrderPGRepository(ctrl)
	sellerRedisRepository := mock.NewMockOrderRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	sellerUC := NewOrderUseCase(cfg, apiLogger, sellerPGRepository, sellerRedisRepository)

	orderUUID := uuid.New()
	userUUID := uuid.New()
	sellerUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		SellerID: sellerUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			SellerID:    sellerUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	ctx := context.Background()

	sellerPGRepository.EXPECT().FindAll(gomock.Any(), nil).AnyTimes().Return(append([]models.Order{}, *mockOrder), nil)

	sellers, err := sellerUC.FindAll(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, sellers)
	require.Equal(t, len(sellers), 1)
}

func TestOrderUseCase_FindAllBySellerId(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sellerPGRepository := mock.NewMockOrderPGRepository(ctrl)
	sellerRedisRepository := mock.NewMockOrderRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	sellerUC := NewOrderUseCase(cfg, apiLogger, sellerPGRepository, sellerRedisRepository)

	orderUUID := uuid.New()
	userUUID := uuid.New()
	sellerUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		SellerID: sellerUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			SellerID:    sellerUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	ctx := context.Background()

	sellerPGRepository.EXPECT().FindAllBySellerId(gomock.Any(), mockOrder.SellerID, nil).AnyTimes().Return(append([]models.Order{}, *mockOrder), nil)

	sellers, err := sellerUC.FindAllBySellerId(ctx, mockOrder.SellerID, nil)
	require.NoError(t, err)
	require.NotNil(t, sellers)
	require.Equal(t, len(sellers), 1)
}

func TestOrderUseCase_FindAllByUserId(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sellerPGRepository := mock.NewMockOrderPGRepository(ctrl)
	sellerRedisRepository := mock.NewMockOrderRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	sellerUC := NewOrderUseCase(cfg, apiLogger, sellerPGRepository, sellerRedisRepository)

	orderUUID := uuid.New()
	userUUID := uuid.New()
	sellerUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		SellerID: sellerUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			SellerID:    sellerUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	ctx := context.Background()

	sellerPGRepository.EXPECT().FindAllByUserId(gomock.Any(), mockOrder.UserID, nil).AnyTimes().Return(append([]models.Order{}, *mockOrder), nil)

	sellers, err := sellerUC.FindAllByUserId(ctx, mockOrder.UserID, nil)
	require.NoError(t, err)
	require.NotNil(t, sellers)
	require.Equal(t, len(sellers), 1)
}

func TestOrderUseCase_FindAllByUserSellerId(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sellerPGRepository := mock.NewMockOrderPGRepository(ctrl)
	sellerRedisRepository := mock.NewMockOrderRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	sellerUC := NewOrderUseCase(cfg, apiLogger, sellerPGRepository, sellerRedisRepository)

	orderUUID := uuid.New()
	userUUID := uuid.New()
	sellerUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		SellerID: sellerUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			SellerID:    sellerUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	ctx := context.Background()

	sellerPGRepository.EXPECT().FindAllByUserIdSellerId(gomock.Any(), mockOrder.UserID, mockOrder.SellerID, nil).AnyTimes().Return(append([]models.Order{}, *mockOrder), nil)

	sellers, err := sellerUC.FindAllByUserIdSellerId(ctx, mockOrder.UserID, mockOrder.SellerID, nil)
	require.NoError(t, err)
	require.NotNil(t, sellers)
	require.Equal(t, len(sellers), 1)
}

func TestOrderUseCase_FindById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderPGRepository := mock.NewMockOrderPGRepository(ctrl)
	orderRedisRepository := mock.NewMockOrderRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	orderUC := NewOrderUseCase(cfg, apiLogger, orderPGRepository, orderRedisRepository)

	orderUUID := uuid.New()
	userUUID := uuid.New()
	sellerUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		SellerID: sellerUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			SellerID:    sellerUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	ctx := context.Background()

	orderRedisRepository.EXPECT().GetByIdCtx(gomock.Any(), mockOrder.OrderID.String()).AnyTimes().Return(nil, redis.Nil)
	orderPGRepository.EXPECT().FindById(gomock.Any(), mockOrder.OrderID).Return(mockOrder, nil)

	order, err := orderUC.FindById(ctx, mockOrder.OrderID)
	require.NoError(t, err)
	require.NotNil(t, order)
	require.Equal(t, order.OrderID, mockOrder.OrderID)

	orderRedisRepository.EXPECT().GetByIdCtx(gomock.Any(), mockOrder.OrderID.String()).AnyTimes().Return(nil, redis.Nil)
}

func TestOrderUseCase_CachedFindById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderPGRepository := mock.NewMockOrderPGRepository(ctrl)
	orderRedisRepository := mock.NewMockOrderRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	orderUC := NewOrderUseCase(cfg, apiLogger, orderPGRepository, orderRedisRepository)

	orderUUID := uuid.New()
	userUUID := uuid.New()
	sellerUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		SellerID: sellerUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			SellerID:    sellerUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	ctx := context.Background()

	orderRedisRepository.EXPECT().GetByIdCtx(gomock.Any(), mockOrder.OrderID.String()).AnyTimes().Return(nil, redis.Nil)
	orderPGRepository.EXPECT().FindById(gomock.Any(), mockOrder.OrderID).Return(mockOrder, nil)
	orderRedisRepository.EXPECT().SetOrderCtx(gomock.Any(), mockOrder.OrderID.String(), 3600, mockOrder).AnyTimes().Return(nil)

	order, err := orderUC.CachedFindById(ctx, mockOrder.OrderID)
	require.NoError(t, err)
	require.NotNil(t, order)
	require.Equal(t, order.OrderID, mockOrder.OrderID)
}

func TestOrderUseCase_UpdateById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderPGRepository := mock.NewMockOrderPGRepository(ctrl)
	orderRedisRepository := mock.NewMockOrderRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	orderUC := NewOrderUseCase(cfg, apiLogger, orderPGRepository, orderRedisRepository)

	orderUUID := uuid.New()
	userUUID := uuid.New()
	sellerUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		SellerID: sellerUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			SellerID:    sellerUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	ctx := context.Background()

	orderPGRepository.EXPECT().UpdateById(gomock.Any(), mockOrder).Return(mockOrder, nil)
	orderRedisRepository.EXPECT().SetOrderCtx(gomock.Any(), mockOrder.OrderID.String(), 3600, mockOrder).AnyTimes().Return(nil)

	order, err := orderUC.UpdateById(ctx, mockOrder)
	require.NoError(t, err)
	require.NotNil(t, order)
	require.Equal(t, order.OrderID, mockOrder.OrderID)
}

func TestOrderUseCase_DeleteById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderPGRepository := mock.NewMockOrderPGRepository(ctrl)
	orderRedisRepository := mock.NewMockOrderRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	orderUC := NewOrderUseCase(cfg, apiLogger, orderPGRepository, orderRedisRepository)

	orderUUID := uuid.New()
	userUUID := uuid.New()
	sellerUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		SellerID: sellerUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			SellerID:    sellerUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	ctx := context.Background()

	orderPGRepository.EXPECT().DeleteById(gomock.Any(), mockOrder.OrderID).Return(nil)
	orderRedisRepository.EXPECT().DeleteOrderCtx(gomock.Any(), mockOrder.OrderID.String()).AnyTimes().Return(nil)

	err := orderUC.DeleteById(ctx, mockOrder.OrderID)
	require.NoError(t, err)

	orderPGRepository.EXPECT().FindById(gomock.Any(), mockOrder.OrderID).AnyTimes().Return(nil, nil)
	orderRedisRepository.EXPECT().GetByIdCtx(gomock.Any(), mockOrder.OrderID.String()).AnyTimes().Return(nil, redis.Nil)
}
