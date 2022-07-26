package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"

	"github.com/dinorain/kalobranded/config"
	"github.com/dinorain/kalobranded/internal/brand"
	"github.com/dinorain/kalobranded/internal/middlewares"
	"github.com/dinorain/kalobranded/internal/models"
	"github.com/dinorain/kalobranded/internal/product"
	"github.com/dinorain/kalobranded/internal/product/delivery/http/dto"
	"github.com/dinorain/kalobranded/internal/session"
	"github.com/dinorain/kalobranded/pkg/constants"
	httpErrors "github.com/dinorain/kalobranded/pkg/http_errors"
	"github.com/dinorain/kalobranded/pkg/logger"
	"github.com/dinorain/kalobranded/pkg/utils"
)

type productHandlersHTTP struct {
	mux       *http.ServeMux
	logger    logger.Logger
	cfg       *config.Config
	mw        middlewares.MiddlewareManager
	v         *validator.Validate
	brandUC   brand.BrandUseCase
	productUC product.ProductUseCase
	sessUC    session.SessUseCase
}

var _ product.ProductHandlers = (*productHandlersHTTP)(nil)

func NewProductHandlersHTTP(
	mux *http.ServeMux,
	logger logger.Logger,
	cfg *config.Config,
	mw middlewares.MiddlewareManager,
	v *validator.Validate,
	brandUC brand.BrandUseCase,
	productUC product.ProductUseCase,
	sessUC session.SessUseCase,
) *productHandlersHTTP {
	return &productHandlersHTTP{mux: mux, logger: logger, cfg: cfg, mw: mw, v: v, brandUC: brandUC, productUC: productUC, sessUC: sessUC}
}

// Create
// @Tags Products
// @Summary Create product
// @Description Create product
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param payload body dto.ProductCreateRequestDto true "Payload"
// @Success 200 {object} dto.ProductCreateResponseDto
// @Router /product [post]
func (h *productHandlersHTTP) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	createDto := &dto.ProductCreateRequestDto{}
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

	if _, err := h.brandUC.CachedFindById(ctx, createDto.BrandID); err != nil {
		h.logger.Errorf("brandUC.CachedFindById: %v", err)
		httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	product, err := h.registerReqToProductModel(createDto)
	if err != nil {
		h.logger.Errorf("registerReqToProductModel: %v", err)
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	createdProduct, err := h.productUC.Create(ctx, product)
	if err != nil {
		h.logger.Errorf("productUC.Create: %v", err)
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	res, _ := json.Marshal(dto.ProductCreateResponseDto{ProductID: createdProduct.ProductID})
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
	return
}

// FindAll
// @Tags Products
// @Summary Find all products
// @Description Find all products
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param size query string false "pagination size"
// @Param page query string false "pagination page"
// @Success 200 {object} dto.ProductFindResponseDto
// @Router /product [get]
func (h *productHandlersHTTP) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queryParam := r.URL.Query()

	if queryParam.Get("id") != "" {
		h.FindById(w, r)
	}

	pq := utils.NewPaginationFromQueryParams(queryParam.Get(constants.Size), queryParam.Get(constants.Page))

	var products []models.Product
	if res, err := h.productUC.FindAll(ctx, pq); err != nil {
		h.logger.Errorf("productUC.FindAll: %v", err)
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	} else {
		products = res
	}

	res, _ := json.Marshal(dto.ProductFindResponseDto{
		Data: products,
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
// @Tags Products
// @Summary Find product by id
// @Description Find product by id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id query string true "product id"
// @Param page query string false "pagination page"
// @Success 200 {object} dto.ProductResponseDto
// @Router /product [get]
func (h *productHandlersHTTP) FindById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queryParam := r.URL.Query()

	if queryParam.Get("id") == "" {
		_ = httpErrors.NewBadRequestError(w, nil, h.cfg.Http.DebugErrorsResponse)
		return
	}
	productUUID, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	product, err := h.productUC.CachedFindById(ctx, productUUID)
	if err != nil {
		h.logger.Errorf("productUC.CachedFindById: %v", err)
		httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	res, _ := json.Marshal(dto.ProductResponseFromModel(product))
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return
}

// FindAllByBrandId
// @Tags Products
// @Summary Find all products by brand
// @Description Find all products by brand is
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id query string false "brand uuid"
// @Param size query string false "pagination size"
// @Param page query string false "pagination page"
// @Success 200 {object} dto.ProductFindResponseDto
// @Router /product [get]
func (h *productHandlersHTTP) FindAllByBrandId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queryParam := r.URL.Query()
	pq := utils.NewPaginationFromQueryParams(queryParam.Get(constants.Size), queryParam.Get(constants.Page))

	var products []models.Product
	id := queryParam.Get(constants.ID)
	brandUUID, err := uuid.Parse(id)
	if err != nil {
		_ = httpErrors.NewBadRequestError(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}
	if res, err := h.productUC.FindAllByBrandId(ctx, brandUUID, pq); err != nil {
		h.logger.Errorf("productUC.FindAllByBrandId: %v", err)
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	} else {
		products = res
	}

	res, _ := json.Marshal(dto.ProductFindResponseDto{
		Data: products,
		Meta: utils.PaginationMetaDto{
			Limit:  pq.GetLimit(),
			Offset: pq.GetOffset(),
			Page:   pq.GetPage(),
		}})
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return
}

func (h *productHandlersHTTP) registerReqToProductModel(r *dto.ProductCreateRequestDto) (*models.Product, error) {
	productCandidate := &models.Product{
		Name:        r.Name,
		Description: r.Description,
		Price:       r.Price,
		BrandID:     r.BrandID,
	}

	if err := productCandidate.PrepareCreate(); err != nil {
		return nil, err
	}

	return productCandidate, nil
}
