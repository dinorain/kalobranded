package usecase

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/dinorain/kalobranded/config"
	"github.com/dinorain/kalobranded/internal/models"
	"github.com/dinorain/kalobranded/internal/order/mock"
	"github.com/dinorain/kalobranded/pkg/logger"
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
	brandUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		BrandID: brandUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			BrandID:    brandUUID,
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
		BrandID: brandUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			BrandID:    brandUUID,
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

	brandPGRepository := mock.NewMockOrderPGRepository(ctrl)
	brandRedisRepository := mock.NewMockOrderRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	brandUC := NewOrderUseCase(cfg, apiLogger, brandPGRepository, brandRedisRepository)

	orderUUID := uuid.New()
	userUUID := uuid.New()
	brandUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		BrandID: brandUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			BrandID:    brandUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	ctx := context.Background()

	brandPGRepository.EXPECT().FindAll(gomock.Any(), nil).AnyTimes().Return(append([]models.Order{}, *mockOrder), nil)

	brands, err := brandUC.FindAll(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, brands)
	require.Equal(t, len(brands), 1)
}

func TestOrderUseCase_FindAllByBrandId(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	brandPGRepository := mock.NewMockOrderPGRepository(ctrl)
	brandRedisRepository := mock.NewMockOrderRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	brandUC := NewOrderUseCase(cfg, apiLogger, brandPGRepository, brandRedisRepository)

	orderUUID := uuid.New()
	userUUID := uuid.New()
	brandUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		BrandID: brandUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			BrandID:    brandUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	ctx := context.Background()

	brandPGRepository.EXPECT().FindAllByBrandId(gomock.Any(), mockOrder.BrandID, nil).AnyTimes().Return(append([]models.Order{}, *mockOrder), nil)

	brands, err := brandUC.FindAllByBrandId(ctx, mockOrder.BrandID, nil)
	require.NoError(t, err)
	require.NotNil(t, brands)
	require.Equal(t, len(brands), 1)
}

func TestOrderUseCase_FindAllByUserId(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	brandPGRepository := mock.NewMockOrderPGRepository(ctrl)
	brandRedisRepository := mock.NewMockOrderRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	brandUC := NewOrderUseCase(cfg, apiLogger, brandPGRepository, brandRedisRepository)

	orderUUID := uuid.New()
	userUUID := uuid.New()
	brandUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		BrandID: brandUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			BrandID:    brandUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	ctx := context.Background()

	brandPGRepository.EXPECT().FindAllByUserId(gomock.Any(), mockOrder.UserID, nil).AnyTimes().Return(append([]models.Order{}, *mockOrder), nil)

	brands, err := brandUC.FindAllByUserId(ctx, mockOrder.UserID, nil)
	require.NoError(t, err)
	require.NotNil(t, brands)
	require.Equal(t, len(brands), 1)
}

func TestOrderUseCase_FindAllByUserBrandId(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	brandPGRepository := mock.NewMockOrderPGRepository(ctrl)
	brandRedisRepository := mock.NewMockOrderRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	brandUC := NewOrderUseCase(cfg, apiLogger, brandPGRepository, brandRedisRepository)

	orderUUID := uuid.New()
	userUUID := uuid.New()
	brandUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		BrandID: brandUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			BrandID:    brandUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	ctx := context.Background()

	brandPGRepository.EXPECT().FindAllByUserIdBrandId(gomock.Any(), mockOrder.UserID, mockOrder.BrandID, nil).AnyTimes().Return(append([]models.Order{}, *mockOrder), nil)

	brands, err := brandUC.FindAllByUserIdBrandId(ctx, mockOrder.UserID, mockOrder.BrandID, nil)
	require.NoError(t, err)
	require.NotNil(t, brands)
	require.Equal(t, len(brands), 1)
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
	brandUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		BrandID: brandUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			BrandID:    brandUUID,
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
	brandUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		BrandID: brandUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			BrandID:    brandUUID,
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
	brandUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		BrandID: brandUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			BrandID:    brandUUID,
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
	brandUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		BrandID: brandUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			BrandID:    brandUUID,
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
