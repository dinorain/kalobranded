package server

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"

	"github.com/dinorain/kalobranded/config"
	"github.com/dinorain/kalobranded/internal/middlewares"
	"github.com/dinorain/kalobranded/pkg/logger"

	brandDeliveryHTTP "github.com/dinorain/kalobranded/internal/brand/delivery/http/handlers"
	orderDeliveryHTTP "github.com/dinorain/kalobranded/internal/order/delivery/http/handlers"
	productDeliveryHTTP "github.com/dinorain/kalobranded/internal/product/delivery/http/handlers"
	userDeliveryHTTP "github.com/dinorain/kalobranded/internal/user/delivery/http/handlers"

	brandUseCase "github.com/dinorain/kalobranded/internal/brand/usecase"
	orderUseCase "github.com/dinorain/kalobranded/internal/order/usecase"
	productUseCase "github.com/dinorain/kalobranded/internal/product/usecase"
	sessUseCase "github.com/dinorain/kalobranded/internal/session/usecase"
	userUseCase "github.com/dinorain/kalobranded/internal/user/usecase"

	brandRepository "github.com/dinorain/kalobranded/internal/brand/repository"
	orderRepository "github.com/dinorain/kalobranded/internal/order/repository"
	productRepository "github.com/dinorain/kalobranded/internal/product/repository"
	sessRepository "github.com/dinorain/kalobranded/internal/session/repository"
	userRepository "github.com/dinorain/kalobranded/internal/user/repository"
)

type Server struct {
	mux         *http.ServeMux
	httpS       *http.Server
	logger      logger.Logger
	cfg         *config.Config
	v           *validator.Validate
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
		mux:         http.NewServeMux(),
		db:          db,
		redisClient: redisClient,
	}
}

// Run service
func (s *Server) Run() error {
	s.mw = middlewares.NewMiddlewareManager(s.logger, s.cfg)

	userRepo := userRepository.NewUserPGRepository(s.db)
	brandRepo := brandRepository.NewBrandPGRepository(s.db)
	productRepo := productRepository.NewProductPGRepository(s.db)
	orderRepo := orderRepository.NewOrderPGRepository(s.db)

	sessRepo := sessRepository.NewSessionRepository(s.redisClient, s.cfg)
	userRedisRepo := userRepository.NewUserRedisRepo(s.redisClient, s.logger)
	brandRedisRepo := brandRepository.NewBrandRedisRepo(s.redisClient, s.logger)
	productRedisRepo := productRepository.NewProductRedisRepo(s.redisClient, s.logger)
	orderRedisRepo := orderRepository.NewOrderRedisRepo(s.redisClient, s.logger)

	sessUC := sessUseCase.NewSessionUseCase(sessRepo, s.cfg)
	userUC := userUseCase.NewUserUseCase(s.cfg, s.logger, userRepo, userRedisRepo)
	brandUC := brandUseCase.NewBrandUseCase(s.cfg, s.logger, brandRepo, brandRedisRepo)
	productUC := productUseCase.NewProductUseCase(s.cfg, s.logger, productRepo, productRedisRepo)
	orderUC := orderUseCase.NewOrderUseCase(s.cfg, s.logger, orderRepo, orderRedisRepo)

	l, err := net.Listen("tcp", s.cfg.Server.Port)
	if err != nil {
		return err
	}
	defer l.Close()

	userHandlers := userDeliveryHTTP.NewUserHandlersHTTP(s.mux, s.logger, s.cfg, s.mw, s.v, userUC, sessUC)
	userHandlers.UserMapRoutes()

	brandHandlers := brandDeliveryHTTP.NewBrandHandlersHTTP(s.mux, s.logger, s.cfg, s.mw, s.v, brandUC, sessUC)
	brandHandlers.BrandMapRoutes()

	productHandlers := productDeliveryHTTP.NewProductHandlersHTTP(s.mux, s.logger, s.cfg, s.mw, s.v, brandUC, productUC, sessUC)
	productHandlers.ProductMapRoutes()

	orderHandlers := orderDeliveryHTTP.NewOrderHandlersHTTP(s.mux, s.logger, s.cfg, s.mw, s.v, orderUC, userUC, brandUC, productUC, sessUC)
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
	if err := s.httpS.Shutdown(ctx); err != nil {
		s.logger.WarnMsg("httpS.Server.Shutdown", err)
	}

	return nil
}
