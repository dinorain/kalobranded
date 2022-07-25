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
	"github.com/dinorain/checkoutaja/internal/product/mock"
	"github.com/dinorain/checkoutaja/pkg/logger"
)

func TestProductUseCase_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	productPGRepository := mock.NewMockProductPGRepository(ctrl)
	productRedisRepository := mock.NewMockProductRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{Server: config.ServerConfig{JwtSecretKey: "secret123"}}
	productUC := NewProductUseCase(cfg, apiLogger, productPGRepository, productRedisRepository)

	productID := uuid.New()
	sellerUUID := uuid.New()
	mockProduct := &models.Product{
		ProductID:   productID,
		Name:        "Name",
		Description: "Description",
		Price:       10000.00,
		SellerID:    sellerUUID,
	}

	ctx := context.Background()

	productPGRepository.EXPECT().Create(gomock.Any(), mockProduct).Return(&models.Product{
		ProductID:   productID,
		Name:        "Name",
		Description: "Description",
		Price:       10000.00,
	}, nil)

	createdProduct, err := productUC.Create(ctx, mockProduct)
	require.NoError(t, err)
	require.NotNil(t, createdProduct)
	require.Equal(t, createdProduct.ProductID, productID)
}

func TestProductUseCase_FindAll(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sellerPGRepository := mock.NewMockProductPGRepository(ctrl)
	sellerRedisRepository := mock.NewMockProductRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	sellerUC := NewProductUseCase(cfg, apiLogger, sellerPGRepository, sellerRedisRepository)

	productID := uuid.New()
	sellerUUID := uuid.New()
	mockProduct := &models.Product{
		ProductID:   productID,
		Name:        "Name",
		Description: "Description",
		Price:       10000.00,
		SellerID:    sellerUUID,
	}

	ctx := context.Background()

	sellerPGRepository.EXPECT().FindAll(gomock.Any(), nil).AnyTimes().Return(append([]models.Product{}, *mockProduct), nil)

	sellers, err := sellerUC.FindAll(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, sellers)
	require.Equal(t, len(sellers), 1)
}

func TestProductUseCase_FindAllBySellerId(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sellerPGRepository := mock.NewMockProductPGRepository(ctrl)
	sellerRedisRepository := mock.NewMockProductRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	sellerUC := NewProductUseCase(cfg, apiLogger, sellerPGRepository, sellerRedisRepository)

	productID := uuid.New()
	sellerUUID := uuid.New()
	mockProduct := &models.Product{
		ProductID:   productID,
		Name:        "Name",
		Description: "Description",
		Price:       10000.00,
		SellerID:    sellerUUID,
	}

	ctx := context.Background()

	sellerPGRepository.EXPECT().FindAllBySellerId(gomock.Any(), sellerUUID, nil).AnyTimes().Return(append([]models.Product{}, *mockProduct), nil)

	sellers, err := sellerUC.FindAllBySellerId(ctx, sellerUUID, nil)
	require.NoError(t, err)
	require.NotNil(t, sellers)
	require.Equal(t, len(sellers), 1)
}

func TestProductUseCase_FindById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	productPGRepository := mock.NewMockProductPGRepository(ctrl)
	productRedisRepository := mock.NewMockProductRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	productUC := NewProductUseCase(cfg, apiLogger, productPGRepository, productRedisRepository)

	productID := uuid.New()
	sellerUUID := uuid.New()
	mockProduct := &models.Product{
		ProductID:   productID,
		Name:        "Name",
		Description: "Description",
		Price:       10000.00,
		SellerID:    sellerUUID,
	}

	ctx := context.Background()

	productRedisRepository.EXPECT().GetByIdCtx(gomock.Any(), mockProduct.ProductID.String()).AnyTimes().Return(nil, redis.Nil)
	productPGRepository.EXPECT().FindById(gomock.Any(), mockProduct.ProductID).Return(mockProduct, nil)

	product, err := productUC.FindById(ctx, mockProduct.ProductID)
	require.NoError(t, err)
	require.NotNil(t, product)
	require.Equal(t, product.ProductID, mockProduct.ProductID)

	productRedisRepository.EXPECT().GetByIdCtx(gomock.Any(), mockProduct.ProductID.String()).AnyTimes().Return(nil, redis.Nil)
}

func TestProductUseCase_CachedFindById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	productPGRepository := mock.NewMockProductPGRepository(ctrl)
	productRedisRepository := mock.NewMockProductRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	productUC := NewProductUseCase(cfg, apiLogger, productPGRepository, productRedisRepository)

	productID := uuid.New()
	sellerUUID := uuid.New()
	mockProduct := &models.Product{
		ProductID:   productID,
		Name:        "Name",
		Description: "Description",
		Price:       10000.00,
		SellerID:    sellerUUID,
	}

	ctx := context.Background()

	productRedisRepository.EXPECT().GetByIdCtx(gomock.Any(), mockProduct.ProductID.String()).AnyTimes().Return(nil, redis.Nil)
	productPGRepository.EXPECT().FindById(gomock.Any(), mockProduct.ProductID).Return(mockProduct, nil)
	productRedisRepository.EXPECT().SetProductCtx(gomock.Any(), mockProduct.ProductID.String(), 3600, mockProduct).AnyTimes().Return(nil)

	product, err := productUC.CachedFindById(ctx, mockProduct.ProductID)
	require.NoError(t, err)
	require.NotNil(t, product)
	require.Equal(t, product.ProductID, mockProduct.ProductID)
}

func TestProductUseCase_UpdateById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	productPGRepository := mock.NewMockProductPGRepository(ctrl)
	productRedisRepository := mock.NewMockProductRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	productUC := NewProductUseCase(cfg, apiLogger, productPGRepository, productRedisRepository)

	productID := uuid.New()
	sellerUUID := uuid.New()
	mockProduct := &models.Product{
		ProductID:   productID,
		Name:        "Name",
		Description: "Description",
		Price:       10000.00,
		SellerID:    sellerUUID,
	}

	ctx := context.Background()

	productPGRepository.EXPECT().UpdateById(gomock.Any(), mockProduct).Return(mockProduct, nil)
	productRedisRepository.EXPECT().SetProductCtx(gomock.Any(), mockProduct.ProductID.String(), 3600, mockProduct).AnyTimes().Return(nil)

	product, err := productUC.UpdateById(ctx, mockProduct)
	require.NoError(t, err)
	require.NotNil(t, product)
	require.Equal(t, product.ProductID, mockProduct.ProductID)
}

func TestProductUseCase_DeleteById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	productPGRepository := mock.NewMockProductPGRepository(ctrl)
	productRedisRepository := mock.NewMockProductRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	productUC := NewProductUseCase(cfg, apiLogger, productPGRepository, productRedisRepository)

	productID := uuid.New()
	sellerUUID := uuid.New()
	mockProduct := &models.Product{
		ProductID:   productID,
		Name:        "Name",
		Description: "Description",
		Price:       10000.00,
		SellerID:    sellerUUID,
	}

	ctx := context.Background()

	productPGRepository.EXPECT().DeleteById(gomock.Any(), mockProduct.ProductID).Return(nil)
	productRedisRepository.EXPECT().DeleteProductCtx(gomock.Any(), mockProduct.ProductID.String()).AnyTimes().Return(nil)

	err := productUC.DeleteById(ctx, mockProduct.ProductID)
	require.NoError(t, err)

	productPGRepository.EXPECT().FindById(gomock.Any(), mockProduct.ProductID).AnyTimes().Return(nil, nil)
	productRedisRepository.EXPECT().GetByIdCtx(gomock.Any(), mockProduct.ProductID.String()).AnyTimes().Return(nil, redis.Nil)
}
