package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"

	"github.com/dinorain/kalobranded/config"
	"github.com/dinorain/kalobranded/internal/brand"
	"github.com/dinorain/kalobranded/internal/brand/delivery/http/dto"
	"github.com/dinorain/kalobranded/internal/middlewares"
	"github.com/dinorain/kalobranded/internal/models"
	"github.com/dinorain/kalobranded/internal/session"
	"github.com/dinorain/kalobranded/pkg/constants"
	httpErrors "github.com/dinorain/kalobranded/pkg/http_errors"
	"github.com/dinorain/kalobranded/pkg/logger"
	"github.com/dinorain/kalobranded/pkg/utils"
)

type brandHandlersHTTP struct {
	mux     *http.ServeMux
	logger  logger.Logger
	cfg     *config.Config
	mw      middlewares.MiddlewareManager
	v       *validator.Validate
	brandUC brand.BrandUseCase
	sessUC  session.SessUseCase
}

var _ brand.BrandHandlers = (*brandHandlersHTTP)(nil)

func NewBrandHandlersHTTP(
	mux *http.ServeMux,
	logger logger.Logger,
	cfg *config.Config,
	mw middlewares.MiddlewareManager,
	v *validator.Validate,
	brandUC brand.BrandUseCase,
	sessUC session.SessUseCase,
) *brandHandlersHTTP {
	return &brandHandlersHTTP{mux: mux, logger: logger, cfg: cfg, mw: mw, v: v, brandUC: brandUC, sessUC: sessUC}
}

// Create
// @Tags Brands
// @Summary Create brand
// @Description Create brand
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param payload body dto.BrandRegisterRequestDto true "Payload"
// @Success 200 {object} dto.BrandRegisterResponseDto
// @Router /brand/create [post]
func (h *brandHandlersHTTP) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	createDto := &dto.BrandRegisterRequestDto{}
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

	brand, err := h.registerReqToBrandModel(createDto)
	if err != nil {
		h.logger.Errorf("registerReqToBrandModel: %v", err)
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	createdBrand, err := h.brandUC.Register(ctx, brand)
	if err != nil {
		h.logger.Errorf("brandUC.Create: %v", err)
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	res, _ := json.Marshal(dto.BrandRegisterResponseDto{BrandID: createdBrand.BrandID})
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
	return
}

// FindAll
// @Tags Brands
// @Summary Find all brands
// @Description Find all brands
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id query string false "brand"
// @Param size query string false "pagination size"
// @Param page query string false "pagination page"
// @Success 200 {object} dto.BrandFindResponseDto
// @Router /brand [get]
func (h *brandHandlersHTTP) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queryParam := r.URL.Query()

	if queryParam.Get("id") != "" {
		brandUUID, err := uuid.Parse(r.URL.Query().Get("id"))
		if err != nil {
			_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
			return
		}

		brand, err := h.brandUC.CachedFindById(ctx, brandUUID)
		if err != nil {
			h.logger.Errorf("brandUC.CachedFindById: %v", err)
			_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
			return
		}

		res, _ := json.Marshal(dto.BrandResponseFromModel(brand))
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	}

	pq := utils.NewPaginationFromQueryParams(queryParam.Get(constants.Size), queryParam.Get(constants.Page))

	var brands []models.Brand
	if res, err := h.brandUC.FindAll(ctx, pq); err != nil {
		h.logger.Errorf("brandUC.FindAll: %v", err)
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	} else {
		brands = res
	}

	res, _ := json.Marshal(dto.BrandFindResponseDto{
		Data: brands,
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
// @Tags Brands
// @Summary Find brand by id
// @Description Find brand by id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id query string true "brand uuid"
// @Success 200 {object} dto.BrandResponseDto
// @Router /brand [get]
func (h *brandHandlersHTTP) FindById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queryParam := r.URL.Query()

	if queryParam.Get("id") == "" {
		_ = httpErrors.NewBadRequestError(w, nil, h.cfg.Http.DebugErrorsResponse)
		return
	}
	brandUUID, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	brand, err := h.brandUC.CachedFindById(ctx, brandUUID)
	if err != nil {
		h.logger.Errorf("brandUC.CachedFindById: %v", err)
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	res, _ := json.Marshal(dto.BrandResponseFromModel(brand))
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return
}

func (h *brandHandlersHTTP) registerReqToBrandModel(r *dto.BrandRegisterRequestDto) (*models.Brand, error) {
	brandCandidate := &models.Brand{
		BrandName:     r.BrandName,
		PickupAddress: r.PickupAddress,
		Logo:          r.Logo,
	}

	if err := brandCandidate.PrepareCreate(); err != nil {
		return nil, err
	}

	return brandCandidate, nil
}
