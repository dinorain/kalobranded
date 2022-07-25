package usecase

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/dinorain/kalobranded/config"
	"github.com/dinorain/kalobranded/internal/brand/mock"
	"github.com/dinorain/kalobranded/internal/models"
	"github.com/dinorain/kalobranded/pkg/logger"
)

func TestBrandUseCase_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	brandPGRepository := mock.NewMockBrandPGRepository(ctrl)
	brandRedisRepository := mock.NewMockBrandRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{Server: config.ServerConfig{JwtSecretKey: "secret123"}}
	brandUC := NewBrandUseCase(cfg, apiLogger, brandPGRepository, brandRedisRepository)

	brandID := uuid.New()
	mockBrand := &models.Brand{
		BrandName:     "BrandName",
		Logo:          nil,
		PickupAddress: "PickupAddress",
	}

	ctx := context.Background()

	brandPGRepository.EXPECT().Create(gomock.Any(), mockBrand).Return(&models.Brand{
		BrandID:       brandID,
		BrandName:     "BrandName",
		Logo:          nil,
		PickupAddress: "PickupAddress",
	}, nil)

	createdBrand, err := brandUC.Register(ctx, mockBrand)
	require.NoError(t, err)
	require.NotNil(t, createdBrand)
	require.Equal(t, createdBrand.BrandID, brandID)
}

func TestBrandUseCase_FindAll(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	brandPGRepository := mock.NewMockBrandPGRepository(ctrl)
	brandRedisRepository := mock.NewMockBrandRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	brandUC := NewBrandUseCase(cfg, apiLogger, brandPGRepository, brandRedisRepository)

	brandID := uuid.New()
	mockBrand := &models.Brand{
		BrandID:       brandID,
		BrandName:     "BrandName",
		Logo:          nil,
		PickupAddress: "DeliveryAddress",
	}

	ctx := context.Background()

	brandPGRepository.EXPECT().FindAll(gomock.Any(), nil).AnyTimes().Return(append([]models.Brand{}, *mockBrand), nil)

	brands, err := brandUC.FindAll(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, brands)
	require.Equal(t, len(brands), 1)
}

func TestBrandUseCase_FindById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	brandPGRepository := mock.NewMockBrandPGRepository(ctrl)
	brandRedisRepository := mock.NewMockBrandRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	brandUC := NewBrandUseCase(cfg, apiLogger, brandPGRepository, brandRedisRepository)

	brandID := uuid.New()
	mockBrand := &models.Brand{
		BrandID:       brandID,
		BrandName:     "BrandName",
		Logo:          nil,
		PickupAddress: "PickupAddress",
	}

	ctx := context.Background()

	brandRedisRepository.EXPECT().GetByIdCtx(gomock.Any(), mockBrand.BrandID.String()).AnyTimes().Return(nil, redis.Nil)
	brandPGRepository.EXPECT().FindById(gomock.Any(), mockBrand.BrandID).Return(mockBrand, nil)

	brand, err := brandUC.FindById(ctx, mockBrand.BrandID)
	require.NoError(t, err)
	require.NotNil(t, brand)
	require.Equal(t, brand.BrandID, mockBrand.BrandID)

	brandRedisRepository.EXPECT().GetByIdCtx(gomock.Any(), mockBrand.BrandID.String()).AnyTimes().Return(nil, redis.Nil)
}

func TestBrandUseCase_CachedFindById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	brandPGRepository := mock.NewMockBrandPGRepository(ctrl)
	brandRedisRepository := mock.NewMockBrandRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	brandUC := NewBrandUseCase(cfg, apiLogger, brandPGRepository, brandRedisRepository)

	brandID := uuid.New()
	mockBrand := &models.Brand{
		BrandID:       brandID,
		BrandName:     "BrandName",
		Logo:          nil,
		PickupAddress: "PickupAddress",
	}

	ctx := context.Background()

	brandRedisRepository.EXPECT().GetByIdCtx(gomock.Any(), mockBrand.BrandID.String()).AnyTimes().Return(nil, redis.Nil)
	brandPGRepository.EXPECT().FindById(gomock.Any(), mockBrand.BrandID).Return(mockBrand, nil)
	brandRedisRepository.EXPECT().SetBrandCtx(gomock.Any(), mockBrand.BrandID.String(), 3600, mockBrand).AnyTimes().Return(nil)

	brand, err := brandUC.CachedFindById(ctx, mockBrand.BrandID)
	require.NoError(t, err)
	require.NotNil(t, brand)
	require.Equal(t, brand.BrandID, mockBrand.BrandID)
}

func TestBrandUseCase_UpdateById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	brandPGRepository := mock.NewMockBrandPGRepository(ctrl)
	brandRedisRepository := mock.NewMockBrandRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	brandUC := NewBrandUseCase(cfg, apiLogger, brandPGRepository, brandRedisRepository)

	brandID := uuid.New()
	mockBrand := &models.Brand{
		BrandID:       brandID,
		BrandName:     "BrandName",
		Logo:          nil,
		PickupAddress: "PickupAddress",
	}

	ctx := context.Background()

	brandPGRepository.EXPECT().UpdateById(gomock.Any(), mockBrand).Return(mockBrand, nil)
	brandRedisRepository.EXPECT().SetBrandCtx(gomock.Any(), mockBrand.BrandID.String(), 3600, mockBrand).AnyTimes().Return(nil)

	brand, err := brandUC.UpdateById(ctx, mockBrand)
	require.NoError(t, err)
	require.NotNil(t, brand)
	require.Equal(t, brand.BrandID, mockBrand.BrandID)
}

func TestBrandUseCase_DeleteById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	brandPGRepository := mock.NewMockBrandPGRepository(ctrl)
	brandRedisRepository := mock.NewMockBrandRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	brandUC := NewBrandUseCase(cfg, apiLogger, brandPGRepository, brandRedisRepository)

	brandID := uuid.New()
	mockBrand := &models.Brand{
		BrandID:       brandID,
		BrandName:     "BrandName",
		Logo:          nil,
		PickupAddress: "PickupAddress",
	}

	ctx := context.Background()

	brandPGRepository.EXPECT().DeleteById(gomock.Any(), mockBrand.BrandID).Return(nil)
	brandRedisRepository.EXPECT().DeleteBrandCtx(gomock.Any(), mockBrand.BrandID.String()).AnyTimes().Return(nil)

	err := brandUC.DeleteById(ctx, mockBrand.BrandID)
	require.NoError(t, err)

	brandPGRepository.EXPECT().FindById(gomock.Any(), mockBrand.BrandID).AnyTimes().Return(nil, nil)
	brandRedisRepository.EXPECT().GetByIdCtx(gomock.Any(), mockBrand.BrandID.String()).AnyTimes().Return(nil, redis.Nil)
}
