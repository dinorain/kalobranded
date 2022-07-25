package server

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"github.com/dinorain/checkoutaja/config"
	"github.com/dinorain/checkoutaja/internal/middlewares"
	"github.com/dinorain/checkoutaja/pkg/logger"

	orderDeliveryHTTP "github.com/dinorain/checkoutaja/internal/order/delivery/http/handlers"
	productDeliveryHTTP "github.com/dinorain/checkoutaja/internal/product/delivery/http/handlers"
	sellerDeliveryHTTP "github.com/dinorain/checkoutaja/internal/seller/delivery/http/handlers"
	userDeliveryHTTP "github.com/dinorain/checkoutaja/internal/user/delivery/http/handlers"

	orderUseCase "github.com/dinorain/checkoutaja/internal/order/usecase"
	productUseCase "github.com/dinorain/checkoutaja/internal/product/usecase"
	sellerUseCase "github.com/dinorain/checkoutaja/internal/seller/usecase"
	sessUseCase "github.com/dinorain/checkoutaja/internal/session/usecase"
	userUseCase "github.com/dinorain/checkoutaja/internal/user/usecase"

	orderRepository "github.com/dinorain/checkoutaja/internal/order/repository"
	productRepository "github.com/dinorain/checkoutaja/internal/product/repository"
	sellerRepository "github.com/dinorain/checkoutaja/internal/seller/repository"
	sessRepository "github.com/dinorain/checkoutaja/internal/session/repository"
	userRepository "github.com/dinorain/checkoutaja/internal/user/repository"
)

type Server struct {
	logger      logger.Logger
	cfg         *config.Config
	v           *validator.Validate
	echo        *echo.Echo
	mw          middlewares.MiddlewareManager
	db          *sqlx.DB
	redisClient *redis.Client
}

// Server constructor
func NewAppServer(logger logger.Logger, cfg *config.Config, db *sqlx.DB, redisClient *redis.Client) *Server {
	return &Server{
		logger:      logger,
		cfg:         cfg,
		v:           validator.New(),
		echo:        echo.New(),
		db:          db,
		redisClient: redisClient,
	}
}

// Run service
func (s *Server) Run() error {
	s.mw = middlewares.NewMiddlewareManager(s.logger, s.cfg)

	userRepo := userRepository.NewUserPGRepository(s.db)
	sellerRepo := sellerRepository.NewSellerPGRepository(s.db)
	productRepo := productRepository.NewProductPGRepository(s.db)
	orderRepo := orderRepository.NewOrderPGRepository(s.db)

	sessRepo := sessRepository.NewSessionRepository(s.redisClient, s.cfg)
	userRedisRepo := userRepository.NewUserRedisRepo(s.redisClient, s.logger)
	sellerRedisRepo := sellerRepository.NewSellerRedisRepo(s.redisClient, s.logger)
	productRedisRepo := productRepository.NewProductRedisRepo(s.redisClient, s.logger)
	orderRedisRepo := orderRepository.NewOrderRedisRepo(s.redisClient, s.logger)

	sessUC := sessUseCase.NewSessionUseCase(sessRepo, s.cfg)
	userUC := userUseCase.NewUserUseCase(s.cfg, s.logger, userRepo, userRedisRepo)
	sellerUC := sellerUseCase.NewSellerUseCase(s.cfg, s.logger, sellerRepo, sellerRedisRepo)
	productUC := productUseCase.NewProductUseCase(s.cfg, s.logger, productRepo, productRedisRepo)
	orderUC := orderUseCase.NewOrderUseCase(s.cfg, s.logger, orderRepo, orderRedisRepo)

	l, err := net.Listen("tcp", s.cfg.Server.Port)
	if err != nil {
		return err
	}
	defer l.Close()

	userHandlers := userDeliveryHTTP.NewUserHandlersHTTP(s.echo.Group("user"), s.logger, s.cfg, s.mw, s.v, userUC, sessUC)
	userHandlers.UserMapRoutes()

	sellerHandlers := sellerDeliveryHTTP.NewSellerHandlersHTTP(s.echo.Group("seller"), s.logger, s.cfg, s.mw, s.v, sellerUC, sessUC)
	sellerHandlers.SellerMapRoutes()

	productHandlers := productDeliveryHTTP.NewProductHandlersHTTP(s.echo.Group("product"), s.logger, s.cfg, s.mw, s.v, productUC, sessUC)
	productHandlers.ProductMapRoutes()

	orderHandlers := orderDeliveryHTTP.NewOrderHandlersHTTP(s.echo.Group("order"), s.logger, s.cfg, s.mw, s.v, orderUC, userUC, sellerUC, productUC, sessUC)
	orderHandlers.OrderMapRoutes()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	go func() {
		if err := s.runHttpServer(); err != nil {
			s.logger.Errorf("s.runHttpServer: %v", err)
			cancel()
		}
	}()

	<-ctx.Done()
	if err := s.echo.Server.Shutdown(ctx); err != nil {
		s.logger.WarnMsg("echo.Server.Shutdown", err)
	}

	return nil
}
