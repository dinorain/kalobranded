package usecase

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/dinorain/checkoutaja/config"
	"github.com/dinorain/checkoutaja/internal/models"
	"github.com/dinorain/checkoutaja/internal/seller"
	"github.com/dinorain/checkoutaja/pkg/grpc_errors"
	"github.com/dinorain/checkoutaja/pkg/logger"
	"github.com/dinorain/checkoutaja/pkg/utils"
)

const (
	sellerByIdCacheDuration = 3600
)

// Seller UseCase
type sellerUseCase struct {
	cfg          *config.Config
	logger       logger.Logger
	sellerPgRepo seller.SellerPGRepository
	redisRepo    seller.SellerRedisRepository
}

var _ seller.SellerUseCase = (*sellerUseCase)(nil)

// New Seller UseCase
func NewSellerUseCase(cfg *config.Config, logger logger.Logger, sellerRepo seller.SellerPGRepository, redisRepo seller.SellerRedisRepository) *sellerUseCase {
	return &sellerUseCase{cfg: cfg, logger: logger, sellerPgRepo: sellerRepo, redisRepo: redisRepo}
}

// Register new seller
func (u *sellerUseCase) Register(ctx context.Context, seller *models.Seller) (*models.Seller, error) {
	existsSeller, err := u.sellerPgRepo.FindByEmail(ctx, seller.Email)
	if existsSeller != nil || err == nil {
		return nil, grpc_errors.ErrEmailExists
	}

	return u.sellerPgRepo.Create(ctx, seller)
}

// FindAll find sellers
func (u *sellerUseCase) FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Seller, error) {
	sellers, err := u.sellerPgRepo.FindAll(ctx, pagination)
	if err != nil {
		return nil, errors.Wrap(err, "sellerPgRepo.FindAll")
	}

	return sellers, nil
}

// FindByEmail find seller by email address
func (u *sellerUseCase) FindByEmail(ctx context.Context, email string) (*models.Seller, error) {
	findByEmail, err := u.sellerPgRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "sellerPgRepo.FindByEmail")
	}

	findByEmail.SanitizePassword()

	return findByEmail, nil
}

// FindById find seller by uuid
func (u *sellerUseCase) FindById(ctx context.Context, sellerID uuid.UUID) (*models.Seller, error) {
	foundSeller, err := u.sellerPgRepo.FindById(ctx, sellerID)
	if err != nil {
		return nil, errors.Wrap(err, "sellerPgRepo.FindById")
	}

	return foundSeller, nil
}

// CachedFindById find seller by uuid from cache
func (u *sellerUseCase) CachedFindById(ctx context.Context, sellerID uuid.UUID) (*models.Seller, error) {
	cachedSeller, err := u.redisRepo.GetByIdCtx(ctx, sellerID.String())
	if err != nil && !errors.Is(err, redis.Nil) {
		u.logger.Errorf("redisRepo.GetByIdCtx", err)
	}
	if cachedSeller != nil {
		return cachedSeller, nil
	}

	foundSeller, err := u.sellerPgRepo.FindById(ctx, sellerID)
	if err != nil {
		return nil, errors.Wrap(err, "sellerPgRepo.FindById")
	}

	if err := u.redisRepo.SetSellerCtx(ctx, foundSeller.SellerID.String(), sellerByIdCacheDuration, foundSeller); err != nil {
		u.logger.Errorf("redisRepo.SetSellerCtx", err)
	}

	return foundSeller, nil
}

// UpdateById update seller by uuid
func (u *sellerUseCase) UpdateById(ctx context.Context, seller *models.Seller) (*models.Seller, error) {
	updatedSeller, err := u.sellerPgRepo.UpdateById(ctx, seller)
	if err != nil {
		return nil, errors.Wrap(err, "sellerPgRepo.UpdateById")
	}

	if err := u.redisRepo.SetSellerCtx(ctx, updatedSeller.SellerID.String(), sellerByIdCacheDuration, updatedSeller); err != nil {
		u.logger.Errorf("redisRepo.SetSellerCtx", err)
	}

	updatedSeller.SanitizePassword()

	return updatedSeller, nil
}

// DeleteById delete seller by uuid
func (u *sellerUseCase) DeleteById(ctx context.Context, sellerID uuid.UUID) error {
	err := u.sellerPgRepo.DeleteById(ctx, sellerID)
	if err != nil {
		return errors.Wrap(err, "sellerPgRepo.DeleteById")
	}

	if err := u.redisRepo.DeleteSellerCtx(ctx, sellerID.String()); err != nil {
		u.logger.Errorf("redisRepo.DeleteSellerCtx", err)
	}

	return nil
}

// Login seller with email and password
func (u *sellerUseCase) Login(ctx context.Context, email string, password string) (*models.Seller, error) {
	foundSeller, err := u.sellerPgRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "sellerPgRepo.FindByEmail")
	}

	if err := foundSeller.ComparePasswords(password); err != nil {
		return nil, errors.Wrap(err, "seller.ComparePasswords")
	}

	return foundSeller, err
}

func (u *sellerUseCase) GenerateTokenPair(seller *models.Seller, sessionID string) (access string, refresh string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["session_id"] = sessionID
	claims["seller_id"] = seller.SellerID
	claims["email"] = seller.Email
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	access, err = token.SignedString([]byte(u.cfg.Server.JwtSecretKey))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["session_id"] = sessionID
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	refresh, err = refreshToken.SignedString([]byte(u.cfg.Server.JwtSecretKey))
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}
