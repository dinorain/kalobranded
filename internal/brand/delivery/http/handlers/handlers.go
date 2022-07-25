package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

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
	group   *echo.Group
	logger  logger.Logger
	cfg     *config.Config
	mw      middlewares.MiddlewareManager
	v       *validator.Validate
	brandUC brand.BrandUseCase
	sessUC  session.SessUseCase
}

var _ brand.BrandHandlers = (*brandHandlersHTTP)(nil)

func NewBrandHandlersHTTP(
	group *echo.Group,
	logger logger.Logger,
	cfg *config.Config,
	mw middlewares.MiddlewareManager,
	v *validator.Validate,
	brandUC brand.BrandUseCase,
	sessUC session.SessUseCase,
) *brandHandlersHTTP {
	return &brandHandlersHTTP{group: group, logger: logger, cfg: cfg, mw: mw, v: v, brandUC: brandUC, sessUC: sessUC}
}

// Register
// @Tags Brands
// @Summary To register brand
// @Description Admin create brand
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param payload body dto.BrandRegisterRequestDto true "Payload"
// @Success 200 {object} dto.BrandRegisterResponseDto
// @Router /brand [post]
func (h *brandHandlersHTTP) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		createDto := &dto.BrandRegisterRequestDto{}
		if err := c.Bind(createDto); err != nil {
			h.logger.WarnMsg("bind", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.v.StructCtx(ctx, createDto); err != nil {
			h.logger.WarnMsg("validate", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		brand, err := h.registerReqToBrandModel(createDto)
		if err != nil {
			h.logger.Errorf("registerReqToBrandModel: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		createdBrand, err := h.brandUC.Register(ctx, brand)
		if err != nil {
			h.logger.Errorf("brandUC.Register: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusCreated, dto.BrandRegisterResponseDto{BrandID: createdBrand.BrandID})
	}
}

// FindAll
// @Tags Brands
// @Summary Find all brands
// @Description Admin find all brands
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param size query string false "pagination size"
// @Param page query string false "pagination page"
// @Success 200 {object} dto.BrandFindResponseDto
// @Router /brand [get]
func (h *brandHandlersHTTP) FindAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		pq := utils.NewPaginationFromQueryParams(c.QueryParam(constants.Size), c.QueryParam(constants.Page))
		brands, err := h.brandUC.FindAll(ctx, pq)
		if err != nil {
			h.logger.Errorf("brandUC.FindAll: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, dto.BrandFindResponseDto{
			Data: brands,
			Meta: utils.PaginationMetaDto{
				Limit:  pq.GetLimit(),
				Offset: pq.GetOffset(),
				Page:   pq.GetPage(),
			},
		})
	}
}

// FindById
// @Tags Brands
// @Summary Find brand
// @Description Find existing brand by id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.BrandResponseDto
// @Router /brand/{id} [get]
func (h *brandHandlersHTTP) FindById() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		brandUUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			h.logger.WarnMsg("uuid.FromString", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		brand, err := h.brandUC.CachedFindById(ctx, brandUUID)
		if err != nil {
			h.logger.Errorf("brandUC.CachedFindById: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, dto.BrandResponseFromModel(brand))
	}
}

// UpdateById
// @Tags Brands
// @Summary Update brand
// @Description Update existing brand
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Brand ID"
// @Param payload body dto.BrandUpdateRequestDto true "Payload"
// @Success 200 {object} dto.BrandResponseDto
// @Router /brand/{id} [put]
func (h *brandHandlersHTTP) UpdateById() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		brandUUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			h.logger.WarnMsg("uuid.FromString", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		_, brandID, role, err := h.getSessionIDFromCtx(c)
		if err != nil {
			h.logger.Errorf("getSessionIDFromCtx: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if role != models.UserRoleAdmin && brandID != brandUUID.String() {
			return httpErrors.NewForbiddenError(c, nil, h.cfg.Http.DebugErrorsResponse)
		}

		updateDto := &dto.BrandUpdateRequestDto{}
		if err := c.Bind(updateDto); err != nil {
			h.logger.WarnMsg("bind", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.v.StructCtx(ctx, updateDto); err != nil {
			h.logger.WarnMsg("validate", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		brand, err := h.brandUC.FindById(ctx, brandUUID)
		if err != nil {
			h.logger.Errorf("brandUC.FindById: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		brand, err = h.updateReqToBrandModel(brand, updateDto)
		if err != nil {
			h.logger.Errorf("updateReqToBrandModel: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		brand, err = h.brandUC.UpdateById(ctx, brand)
		if err != nil {
			h.logger.Errorf("brandUC.UpdateById: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, dto.BrandResponseFromModel(brand))
	}
}

// DeleteById
// @Tags Brands
// @Summary Delete brand
// @Description Delete existing brand, admin only
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} nil
// @Param id path string true "Brand ID"
// @Router /brand/{id} [delete]
func (h *brandHandlersHTTP) DeleteById() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		brandUUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			h.logger.WarnMsg("uuid.FromString", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.brandUC.DeleteById(ctx, brandUUID); err != nil {
			h.logger.Errorf("brandUC.DeleteById: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, nil)
	}
}

func (h *brandHandlersHTTP) getSessionIDFromCtx(c echo.Context) (sessionID string, brandID string, role string, err error) {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		h.logger.Warnf("jwt.Token: %+v", c.Get("user"))
		return "", "", "", errors.New("invalid token header")
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		h.logger.Warnf("jwt.MapClaims: %+v", c.Get("user"))
		return "", "", "", errors.New("invalid token header")
	}

	sessionID, ok = claims["session_id"].(string)
	if !ok {
		h.logger.Warnf("session_id: %+v", claims)
		return "", "", "", errors.New("invalid token header")
	}

	brandID, _ = claims["brand_id"].(string)
	role, _ = claims["role"].(string)

	return sessionID, brandID, role, nil
}

func (h *brandHandlersHTTP) registerReqToBrandModel(r *dto.BrandRegisterRequestDto) (*models.Brand, error) {
	brandCandidate := &models.Brand{
		BrandName:     r.BrandName,
		Logo:          nil,
		PickupAddress: r.PickupAddress,
	}

	if err := brandCandidate.PrepareCreate(); err != nil {
		return nil, err
	}

	return brandCandidate, nil
}

func (h *brandHandlersHTTP) updateReqToBrandModel(updateCandidate *models.Brand, r *dto.BrandUpdateRequestDto) (*models.Brand, error) {

	if r.BrandName != nil {
		updateCandidate.BrandName = strings.TrimSpace(*r.BrandName)
	}
	if r.PickupAddress != nil {
		updateCandidate.PickupAddress = strings.TrimSpace(*r.PickupAddress)
	}
	if r.Logo != nil {
		logo := strings.TrimSpace(*r.Logo)
		updateCandidate.Logo = &logo
	}

	return updateCandidate, nil
}
