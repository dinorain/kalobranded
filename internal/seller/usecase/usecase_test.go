package usecase

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/dinorain/checkoutaja/config"
	"github.com/dinorain/checkoutaja/internal/models"
	"github.com/dinorain/checkoutaja/internal/seller/mock"
	"github.com/dinorain/checkoutaja/pkg/logger"
)

func TestSellerUseCase_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sellerPGRepository := mock.NewMockSellerPGRepository(ctrl)
	sellerRedisRepository := mock.NewMockSellerRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{Server: config.ServerConfig{JwtSecretKey: "secret123"}}
	sellerUC := NewSellerUseCase(cfg, apiLogger, sellerPGRepository, sellerRedisRepository)

	sellerID := uuid.New()
	mockSeller := &models.Seller{
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
		PickupAddress: "PickupAddress",
	}

	ctx := context.Background()

	sellerPGRepository.EXPECT().FindByEmail(gomock.Any(), mockSeller.Email).Return(nil, sql.ErrNoRows)

	sellerPGRepository.EXPECT().Create(gomock.Any(), mockSeller).Return(&models.Seller{
		SellerID:      sellerID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
		PickupAddress: "PickupAddress",
	}, nil)

	createdSeller, err := sellerUC.Register(ctx, mockSeller)
	require.NoError(t, err)
	require.NotNil(t, createdSeller)
	require.Equal(t, createdSeller.SellerID, sellerID)
}

func TestSellerUseCase_FindByEmail(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sellerPGRepository := mock.NewMockSellerPGRepository(ctrl)
	sellerRedisRepository := mock.NewMockSellerRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{Server: config.ServerConfig{JwtSecretKey: "secret123"}}
	sellerUC := NewSellerUseCase(cfg, apiLogger, sellerPGRepository, sellerRedisRepository)

	sellerID := uuid.New()
	mockSeller := &models.Seller{
		SellerID:      sellerID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
		PickupAddress: "PickupAddress",
	}

	ctx := context.Background()

	sellerPGRepository.EXPECT().FindByEmail(gomock.Any(), mockSeller.Email).Return(mockSeller, nil)

	seller, err := sellerUC.FindByEmail(ctx, mockSeller.Email)
	require.NoError(t, err)
	require.NotNil(t, seller)
	require.Equal(t, seller.Email, mockSeller.Email)
}

func TestSellerUseCase_Login(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sellerPGRepository := mock.NewMockSellerPGRepository(ctrl)
	sellerRedisRepository := mock.NewMockSellerRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{Server: config.ServerConfig{JwtSecretKey: "secret123"}}
	sellerUC := NewSellerUseCase(cfg, apiLogger, sellerPGRepository, sellerRedisRepository)

	sellerID := uuid.New()
	mockSeller := &models.Seller{
		SellerID:      sellerID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
		PickupAddress: "PickupAddress",
	}

	ctx := context.Background()

	sellerPGRepository.EXPECT().FindByEmail(gomock.Any(), mockSeller.Email).Return(mockSeller, nil)
	_, err := sellerUC.Login(ctx, mockSeller.Email, mockSeller.Password)
	require.NotNil(t, err)
}

func TestSellerUseCase_FindAll(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sellerPGRepository := mock.NewMockSellerPGRepository(ctrl)
	sellerRedisRepository := mock.NewMockSellerRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	sellerUC := NewSellerUseCase(cfg, apiLogger, sellerPGRepository, sellerRedisRepository)

	sellerID := uuid.New()
	mockSeller := &models.Seller{
		SellerID:      sellerID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
		PickupAddress: "DeliveryAddress",
	}

	ctx := context.Background()

	sellerPGRepository.EXPECT().FindAll(gomock.Any(), nil).AnyTimes().Return(append([]models.Seller{}, *mockSeller), nil)

	sellers, err := sellerUC.FindAll(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, sellers)
	require.Equal(t, len(sellers), 1)
}

func TestSellerUseCase_FindById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sellerPGRepository := mock.NewMockSellerPGRepository(ctrl)
	sellerRedisRepository := mock.NewMockSellerRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	sellerUC := NewSellerUseCase(cfg, apiLogger, sellerPGRepository, sellerRedisRepository)

	sellerID := uuid.New()
	mockSeller := &models.Seller{
		SellerID:      sellerID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
		PickupAddress: "PickupAddress",
	}

	ctx := context.Background()

	sellerRedisRepository.EXPECT().GetByIdCtx(gomock.Any(), mockSeller.SellerID.String()).AnyTimes().Return(nil, redis.Nil)
	sellerPGRepository.EXPECT().FindById(gomock.Any(), mockSeller.SellerID).Return(mockSeller, nil)

	seller, err := sellerUC.FindById(ctx, mockSeller.SellerID)
	require.NoError(t, err)
	require.NotNil(t, seller)
	require.Equal(t, seller.SellerID, mockSeller.SellerID)

	sellerRedisRepository.EXPECT().GetByIdCtx(gomock.Any(), mockSeller.SellerID.String()).AnyTimes().Return(nil, redis.Nil)
}

func TestSellerUseCase_CachedFindById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sellerPGRepository := mock.NewMockSellerPGRepository(ctrl)
	sellerRedisRepository := mock.NewMockSellerRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	sellerUC := NewSellerUseCase(cfg, apiLogger, sellerPGRepository, sellerRedisRepository)

	sellerID := uuid.New()
	mockSeller := &models.Seller{
		SellerID:      sellerID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
		PickupAddress: "PickupAddress",
	}

	ctx := context.Background()

	sellerRedisRepository.EXPECT().GetByIdCtx(gomock.Any(), mockSeller.SellerID.String()).AnyTimes().Return(nil, redis.Nil)
	sellerPGRepository.EXPECT().FindById(gomock.Any(), mockSeller.SellerID).Return(mockSeller, nil)
	sellerRedisRepository.EXPECT().SetSellerCtx(gomock.Any(), mockSeller.SellerID.String(), 3600, mockSeller).AnyTimes().Return(nil)

	seller, err := sellerUC.CachedFindById(ctx, mockSeller.SellerID)
	require.NoError(t, err)
	require.NotNil(t, seller)
	require.Equal(t, seller.SellerID, mockSeller.SellerID)
}

func TestSellerUseCase_UpdateById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sellerPGRepository := mock.NewMockSellerPGRepository(ctrl)
	sellerRedisRepository := mock.NewMockSellerRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	sellerUC := NewSellerUseCase(cfg, apiLogger, sellerPGRepository, sellerRedisRepository)

	sellerID := uuid.New()
	mockSeller := &models.Seller{
		SellerID:      sellerID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
		PickupAddress: "PickupAddress",
	}

	ctx := context.Background()

	sellerPGRepository.EXPECT().UpdateById(gomock.Any(), mockSeller).Return(mockSeller, nil)
	sellerRedisRepository.EXPECT().SetSellerCtx(gomock.Any(), mockSeller.SellerID.String(), 3600, mockSeller).AnyTimes().Return(nil)

	seller, err := sellerUC.UpdateById(ctx, mockSeller)
	require.NoError(t, err)
	require.NotNil(t, seller)
	require.Equal(t, seller.SellerID, mockSeller.SellerID)
}

func TestSellerUseCase_DeleteById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sellerPGRepository := mock.NewMockSellerPGRepository(ctrl)
	sellerRedisRepository := mock.NewMockSellerRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	sellerUC := NewSellerUseCase(cfg, apiLogger, sellerPGRepository, sellerRedisRepository)

	sellerID := uuid.New()
	mockSeller := &models.Seller{
		SellerID:      sellerID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
		PickupAddress: "PickupAddress",
	}

	ctx := context.Background()

	sellerPGRepository.EXPECT().DeleteById(gomock.Any(), mockSeller.SellerID).Return(nil)
	sellerRedisRepository.EXPECT().DeleteSellerCtx(gomock.Any(), mockSeller.SellerID.String()).AnyTimes().Return(nil)

	err := sellerUC.DeleteById(ctx, mockSeller.SellerID)
	require.NoError(t, err)

	sellerPGRepository.EXPECT().FindById(gomock.Any(), mockSeller.SellerID).AnyTimes().Return(nil, nil)
	sellerRedisRepository.EXPECT().GetByIdCtx(gomock.Any(), mockSeller.SellerID.String()).AnyTimes().Return(nil, redis.Nil)
}

func TestSellerUseCase_GenerateTokenPair(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sellerPGRepository := mock.NewMockSellerPGRepository(ctrl)
	sellerRedisRepository := mock.NewMockSellerRedisRepository(ctrl)
	apiLogger := logger.NewAppLogger(nil)

	cfg := &config.Config{}
	sellerUC := NewSellerUseCase(cfg, apiLogger, sellerPGRepository, sellerRedisRepository)

	sellerID := uuid.New()
	mockSeller := &models.Seller{
		SellerID:      sellerID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
		PickupAddress: "PickupAddress",
	}

	at, rt, err := sellerUC.GenerateTokenPair(mockSeller, mockSeller.SellerID.String())
	require.NoError(t, err)
	require.NotEqual(t, at, "")
	require.NotEqual(t, rt, "")
}
