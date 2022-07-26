package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"github.com/dinorain/kalobranded/config"
	"github.com/dinorain/kalobranded/internal/brand"
	"github.com/dinorain/kalobranded/internal/middlewares"
	"github.com/dinorain/kalobranded/internal/models"
	"github.com/dinorain/kalobranded/internal/order"
	"github.com/dinorain/kalobranded/internal/order/delivery/http/dto"
	"github.com/dinorain/kalobranded/internal/product"
	"github.com/dinorain/kalobranded/internal/session"
	"github.com/dinorain/kalobranded/internal/user"
	"github.com/dinorain/kalobranded/pkg/constants"
	httpErrors "github.com/dinorain/kalobranded/pkg/http_errors"
	"github.com/dinorain/kalobranded/pkg/logger"
	"github.com/dinorain/kalobranded/pkg/utils"
)

type orderHandlersHTTP struct {
	mux       *http.ServeMux
	logger    logger.Logger
	cfg       *config.Config
	mw        middlewares.MiddlewareManager
	v         *validator.Validate
	orderUC   order.OrderUseCase
	userUC    user.UserUseCase
	brandUC   brand.BrandUseCase
	productUC product.ProductUseCase
	sessUC    session.SessUseCase
}

var _ order.OrderHandlers = (*orderHandlersHTTP)(nil)

func NewOrderHandlersHTTP(
	mux *http.ServeMux,
	logger logger.Logger,
	cfg *config.Config,
	mw middlewares.MiddlewareManager,
	v *validator.Validate,
	orderUC order.OrderUseCase,
	userUC user.UserUseCase,
	brandUC brand.BrandUseCase,
	productUC product.ProductUseCase,
	sessUC session.SessUseCase,
) *orderHandlersHTTP {
	return &orderHandlersHTTP{mux: mux, logger: logger, cfg: cfg, mw: mw, v: v, orderUC: orderUC, userUC: userUC, brandUC: brandUC, productUC: productUC, sessUC: sessUC}
}

// Create
// @Tags Orders
// @Summary To create order
// @Description Order create order
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param payload body dto.OrderCreateRequestDto true "Payload"
// @Success 200 {object} dto.OrderCreateResponseDto
// @Router /order [post]
func (h *orderHandlersHTTP) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	createDto := &dto.OrderCreateRequestDto{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&createDto); err != nil {
		h.logger.Errorf("decoder.Decode: %v", err)
		_ = httpErrors.NewBadRequestError(w, err.Error(), h.cfg.Http.DebugErrorsResponse)
		return
	}

	if err := h.v.Struct(createDto); err != nil {
		h.logger.Errorf("h.v.Struct: %v", err)
		_ = httpErrors.NewBadRequestError(w, err.Error(), h.cfg.Http.DebugErrorsResponse)
		return
	}

	sessID, _, _, err := h.getSessionIDFromCtx(w, r)
	if err != nil {
		h.logger.Errorf("getSessionIDFromCtx: %v", err)
		return
	}

	session, err := h.sessUC.GetSessionById(ctx, sessID)
	if err != nil {
		h.logger.Errorf("sessUC.GetSessionById: %v", err)
		if errors.Is(err, redis.Nil) {
			_ = httpErrors.NewUnauthorizedError(w, nil, h.cfg.Http.DebugErrorsResponse)
			return
		}
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	user, err := h.userUC.CachedFindById(ctx, session.UserID)
	if err != nil {
		h.logger.Errorf("userUC.CachedFindById: %v", err)
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	product, err := h.productUC.CachedFindById(ctx, createDto.ProductID)
	if err != nil {
		h.logger.Errorf("productUC.CachedFindById: %v", err)
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	brand, err := h.brandUC.CachedFindById(ctx, product.BrandID)
	if err != nil {
		h.logger.Errorf("brandUC.CachedFindById: %v", err)
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	order, err := h.registerReqToOrderModel(createDto, user, brand, product)
	if err != nil {
		h.logger.Errorf("orderHandlersHTTP.registerReqToOrderModel: %v", err)
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	createdOrder, err := h.orderUC.Create(ctx, order)
	if err != nil {
		h.logger.Errorf("orderUC.Create: %v", err)
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	res, _ := json.Marshal(dto.OrderCreateResponseDto{OrderID: createdOrder.OrderID})
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
	return
}

// FindAll
// @Tags Orders
// @Summary Find all orders
// @Description Find all orders, will be filtered by user_id when accessed by user role
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param size query string false "pagination size"
// @Param page query string false "pagination page"
// @Success 200 {object} dto.OrderFindResponseDto
// @Router /order [get]
func (h *orderHandlersHTTP) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queryParam := r.URL.Query()
	pq := utils.NewPaginationFromQueryParams(queryParam.Get(constants.Size), queryParam.Get(constants.Page))

	var orders []models.Order

	sessID, _, role, err := h.getSessionIDFromCtx(w, r)
	if err != nil {
		h.logger.Errorf("getSessionIDFromCtx: %v", err)
		return
	}

	session, err := h.sessUC.GetSessionById(ctx, sessID)
	if err != nil {
		h.logger.Errorf("sessUC.GetSessionById: %v", err)
		if errors.Is(err, redis.Nil) {
			_ = httpErrors.NewUnauthorizedError(w, nil, h.cfg.Http.DebugErrorsResponse)
			return
		}
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	if role == models.UserRoleAdmin {
		if res, err := h.orderUC.FindAll(ctx, pq); err != nil {
			h.logger.Errorf("orderUC.FindAll: %v", err)
			_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
			return
		} else {
			orders = res
		}

	} else {
		userUUID := session.UserID
		if res, err := h.orderUC.FindAllByUserId(ctx, userUUID, pq); err != nil {
			h.logger.Errorf("orderUC.FindAll: %v", err)
			_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
			return
		} else {
			orders = res
		}

	}

	res, _ := json.Marshal(dto.OrderFindResponseDto{
		Data: orders,
		Meta: utils.PaginationMetaDto{
			Limit:  pq.GetLimit(),
			Offset: pq.GetOffset(),
			Page:   pq.GetPage(),
		}})
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return
}

// FindById
// @Tags Orders
// @Summary Find order by id
// @Description Find order by id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id query string false "order uuid"
// @Success 200 {object} dto.OrderResponseDto
// @Router /order [get]
func (h *orderHandlersHTTP) FindById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queryParam := r.URL.Query()

	if queryParam.Get("id") == "" {
		_ = httpErrors.NewBadRequestError(w, nil, h.cfg.Http.DebugErrorsResponse)
		return
	}
	orderUUID, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	order, err := h.orderUC.CachedFindById(ctx, orderUUID)
	if err != nil {
		h.logger.Errorf("orderUC.CachedFindById: %v", err)
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	res, _ := json.Marshal(dto.OrderResponseFromModel(order))
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return
}

func (h *orderHandlersHTTP) registerReqToOrderModel(r *dto.OrderCreateRequestDto, user *models.User, brand *models.Brand, product *models.Product) (*models.Order, error) {
	orderCandidate := &models.Order{
		UserID:  user.UserID,
		BrandID: brand.BrandID,
		Item: models.OrderItem{
			ProductID:   product.ProductID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			BrandID:     product.BrandID,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		},
		Quantity:                   r.Quantity,
		TotalPrice:                 float64(r.Quantity) * product.Price,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      brand.PickupAddress,
		DeliveryDestinationAddress: user.DeliveryAddress,
	}

	return orderCandidate, nil
}

func (h *orderHandlersHTTP) getSessionIDFromCtx(w http.ResponseWriter, r *http.Request) (sessionID string, userID string, role string, err error) {
	jwtClaims, err := h.mw.GetJWTClaims(w, r)
	if err != nil {
		return
	}
	claims := *jwtClaims
	sessionID, ok := claims["session_id"].(string)
	if !ok {
		h.logger.Warnf("session_id: %+v", claims)
		return "", "", "", errors.New("invalid token header")
	}

	role, _ = claims["role"].(string)
	if role != "" {
		userID, _ = claims["user_id"].(string)
	} else {
		userID, _ = claims["brand_id"].(string)
	}
	return sessionID, userID, role, nil
}
