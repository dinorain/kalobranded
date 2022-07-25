package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/dinorain/checkoutaja/config"
	"github.com/dinorain/checkoutaja/internal/middlewares"
	"github.com/dinorain/checkoutaja/internal/models"
	"github.com/dinorain/checkoutaja/internal/seller"
	"github.com/dinorain/checkoutaja/internal/seller/delivery/http/dto"
	"github.com/dinorain/checkoutaja/internal/session"
	"github.com/dinorain/checkoutaja/pkg/constants"
	httpErrors "github.com/dinorain/checkoutaja/pkg/http_errors"
	"github.com/dinorain/checkoutaja/pkg/logger"
	"github.com/dinorain/checkoutaja/pkg/utils"
)

type sellerHandlersHTTP struct {
	group    *echo.Group
	logger   logger.Logger
	cfg      *config.Config
	mw       middlewares.MiddlewareManager
	v        *validator.Validate
	sellerUC seller.SellerUseCase
	sessUC   session.SessUseCase
}

var _ seller.SellerHandlers = (*sellerHandlersHTTP)(nil)

func NewSellerHandlersHTTP(
	group *echo.Group,
	logger logger.Logger,
	cfg *config.Config,
	mw middlewares.MiddlewareManager,
	v *validator.Validate,
	sellerUC seller.SellerUseCase,
	sessUC session.SessUseCase,
) *sellerHandlersHTTP {
	return &sellerHandlersHTTP{group: group, logger: logger, cfg: cfg, mw: mw, v: v, sellerUC: sellerUC, sessUC: sessUC}
}

// Register
// @Tags Sellers
// @Summary To register seller
// @Description Admin create seller
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param payload body dto.SellerRegisterRequestDto true "Payload"
// @Success 200 {object} dto.SellerRegisterResponseDto
// @Router /seller [post]
func (h *sellerHandlersHTTP) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		createDto := &dto.SellerRegisterRequestDto{}
		if err := c.Bind(createDto); err != nil {
			h.logger.WarnMsg("bind", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.v.StructCtx(ctx, createDto); err != nil {
			h.logger.WarnMsg("validate", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		seller, err := h.registerReqToSellerModel(createDto)
		if err != nil {
			h.logger.Errorf("registerReqToSellerModel: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		createdSeller, err := h.sellerUC.Register(ctx, seller)
		if err != nil {
			h.logger.Errorf("sellerUC.Register: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusCreated, dto.SellerRegisterResponseDto{SellerID: createdSeller.SellerID})
	}
}

// Login
// @Tags Sellers
// @Summary Seller login
// @Description Seller login with email and password
// @Accept json
// @Produce json
// @Param payload body dto.SellerLoginRequestDto true "Payload"
// @Success 200 {object} dto.SellerLoginResponseDto
// @Router /seller/login [post]
func (h *sellerHandlersHTTP) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		loginDto := &dto.SellerLoginRequestDto{}
		if err := c.Bind(loginDto); err != nil {
			h.logger.WarnMsg("bind", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.v.StructCtx(ctx, loginDto); err != nil {
			h.logger.WarnMsg("validate", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		email := loginDto.Email
		if !utils.ValidateEmail(email) {
			h.logger.Errorf("ValidateEmail: %v", email)
			return httpErrors.ErrorCtxResponse(c, errors.New("invalid email"), h.cfg.Http.DebugErrorsResponse)
		}

		seller, err := h.sellerUC.Login(ctx, email, loginDto.Password)
		if err != nil {
			h.logger.Errorf("sellerUC.Login: %v", email)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		session, err := h.sessUC.CreateSession(ctx, &models.Session{
			UserID: seller.SellerID,
		}, h.cfg.Session.Expire)
		if err != nil {
			h.logger.Errorf("sessUC.CreateSession: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		accessToken, refreshToken, err := h.sellerUC.GenerateTokenPair(seller, session)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, dto.SellerLoginResponseDto{SellerID: seller.SellerID, Tokens: &dto.SellerRefreshTokenResponseDto{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}})
	}
}

// FindAll
// @Tags Sellers
// @Summary Find all sellers
// @Description Admin find all sellers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param size query string false "pagination size"
// @Param page query string false "pagination page"
// @Success 200 {object} dto.SellerFindResponseDto
// @Router /seller [get]
func (h *sellerHandlersHTTP) FindAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		pq := utils.NewPaginationFromQueryParams(c.QueryParam(constants.Size), c.QueryParam(constants.Page))
		sellers, err := h.sellerUC.FindAll(ctx, pq)
		if err != nil {
			h.logger.Errorf("sellerUC.FindAll: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, dto.SellerFindResponseDto{
			Data: sellers,
			Meta: utils.PaginationMetaDto{
				Limit:  pq.GetLimit(),
				Offset: pq.GetOffset(),
				Page:   pq.GetPage(),
			},
		})
	}
}

// FindById
// @Tags Sellers
// @Summary Find seller
// @Description Find existing seller by id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.SellerResponseDto
// @Router /seller/{id} [get]
func (h *sellerHandlersHTTP) FindById() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		sellerUUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			h.logger.WarnMsg("uuid.FromString", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		seller, err := h.sellerUC.CachedFindById(ctx, sellerUUID)
		if err != nil {
			h.logger.Errorf("sellerUC.CachedFindById: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, dto.SellerResponseFromModel(seller))
	}
}

// UpdateById
// @Tags Sellers
// @Summary Update seller
// @Description Update existing seller
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Seller ID"
// @Param payload body dto.SellerUpdateRequestDto true "Payload"
// @Success 200 {object} dto.SellerResponseDto
// @Router /seller/{id} [put]
func (h *sellerHandlersHTTP) UpdateById() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		sellerUUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			h.logger.WarnMsg("uuid.FromString", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		_, sellerID, role, err := h.getSessionIDFromCtx(c)
		if err != nil {
			h.logger.Errorf("getSessionIDFromCtx: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if role != models.UserRoleAdmin && sellerID != sellerUUID.String() {
			return httpErrors.NewForbiddenError(c, nil, h.cfg.Http.DebugErrorsResponse)
		}

		updateDto := &dto.SellerUpdateRequestDto{}
		if err := c.Bind(updateDto); err != nil {
			h.logger.WarnMsg("bind", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.v.StructCtx(ctx, updateDto); err != nil {
			h.logger.WarnMsg("validate", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		seller, err := h.sellerUC.FindById(ctx, sellerUUID)
		if err != nil {
			h.logger.Errorf("sellerUC.FindById: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		seller, err = h.updateReqToSellerModel(seller, updateDto)
		if err != nil {
			h.logger.Errorf("updateReqToSellerModel: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		seller, err = h.sellerUC.UpdateById(ctx, seller)
		if err != nil {
			h.logger.Errorf("sellerUC.UpdateById: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, dto.SellerResponseFromModel(seller))
	}
}

// DeleteById
// @Tags Sellers
// @Summary Delete seller
// @Description Delete existing seller, admin only
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} nil
// @Param id path string true "Seller ID"
// @Router /seller/{id} [delete]
func (h *sellerHandlersHTTP) DeleteById() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		sellerUUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			h.logger.WarnMsg("uuid.FromString", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.sellerUC.DeleteById(ctx, sellerUUID); err != nil {
			h.logger.Errorf("sellerUC.DeleteById: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, nil)
	}
}

// GetMe
// @Tags Sellers
// @Summary Find me
// @Description Get session id from token, find seller by uuid and returns it
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.SellerResponseDto
// @Router /seller/me [get]
func (h *sellerHandlersHTTP) GetMe() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		sessID, _, _, err := h.getSessionIDFromCtx(c)
		if err != nil {
			h.logger.Errorf("getSessionIDFromCtx: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		session, err := h.sessUC.GetSessionById(ctx, sessID)
		if err != nil {
			h.logger.Errorf("sessUC.GetSessionById: %v", err)
			if errors.Is(err, redis.Nil) {
				return httpErrors.NewUnauthorizedError(c, nil, h.cfg.Http.DebugErrorsResponse)
			}
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		seller, err := h.sellerUC.CachedFindById(ctx, session.UserID)
		if err != nil {
			h.logger.Errorf("sellerUC.CachedFindById: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, dto.SellerResponseFromModel(seller))
	}
}

// Logout
// @Tags Sellers
// @Summary Seller logout
// @Description Delete current session
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} nil
// @Router /seller/logout [post]
func (h *sellerHandlersHTTP) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		sessID, _, _, err := h.getSessionIDFromCtx(c)
		if err != nil {
			h.logger.Errorf("getSessionIDFromCtx: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.sessUC.DeleteById(ctx, sessID); err != nil {
			h.logger.Errorf("sessUC.DeleteById: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, nil)
	}
}

// RefreshToken
// @Tags Sellers
// @Summary Refresh access token
// @Description Refresh access token
// @Accept json
// @Produce json
// @Param payload body dto.SellerRefreshTokenDto true "Payload"
// @Success 200 {object} dto.SellerRefreshTokenResponseDto
// @Router /seller/refresh [post]
func (h *sellerHandlersHTTP) RefreshToken() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		refreshTokenDto := &dto.SellerRefreshTokenDto{}
		if err := c.Bind(refreshTokenDto); err != nil {
			h.logger.WarnMsg("bind", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		token, err := jwt.Parse(refreshTokenDto.RefreshToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				h.logger.Errorf("jwt.SigningMethodHMAC: %v", token.Header["alg"])
				return nil, fmt.Errorf("jwt.SigningMethodHMAC: %v", token.Header["alg"])
			}

			return []byte(h.cfg.Server.JwtSecretKey), nil
		})

		if err != nil {
			h.logger.Warnf("jwt.Parse")
			return httpErrors.ErrorCtxResponse(c, errors.New("invalid refresh token"), h.cfg.Http.DebugErrorsResponse)
		}

		if !token.Valid {
			h.logger.Warnf("token.Valid")
			return httpErrors.ErrorCtxResponse(c, errors.New("invalid refresh token"), h.cfg.Http.DebugErrorsResponse)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			h.logger.Warnf("jwt.MapClaims: %+v", token.Claims)
			return httpErrors.ErrorCtxResponse(c, errors.New("invalid refresh token"), h.cfg.Http.DebugErrorsResponse)
		}

		sessID, ok := claims["session_id"].(string)
		if !ok {
			h.logger.Warnf("session_id: %+v", claims)
			return httpErrors.ErrorCtxResponse(c, errors.New("invalid refresh token"), h.cfg.Http.DebugErrorsResponse)
		}

		session, err := h.sessUC.GetSessionById(ctx, sessID)
		if err != nil {
			h.logger.Errorf("sessUC.GetSessionById: %v", err)
			if errors.Is(err, redis.Nil) {
				return httpErrors.NewUnauthorizedError(c, nil, h.cfg.Http.DebugErrorsResponse)
			}
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		seller, err := h.sellerUC.FindById(ctx, session.UserID)
		if err != nil {
			h.logger.Errorf("sellerUC.FindById: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		accessToken, refreshToken, err := h.sellerUC.GenerateTokenPair(seller, sessID)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, dto.SellerRefreshTokenResponseDto{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	}
}

func (h *sellerHandlersHTTP) getSessionIDFromCtx(c echo.Context) (sessionID string, sellerID string, role string, err error) {
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

	sellerID, _ = claims["seller_id"].(string)
	role, _ = claims["role"].(string)

	return sessionID, sellerID, role, nil
}

func (h *sellerHandlersHTTP) registerReqToSellerModel(r *dto.SellerRegisterRequestDto) (*models.Seller, error) {
	sellerCandidate := &models.Seller{
		Email:         r.Email,
		FirstName:     r.FirstName,
		LastName:      r.LastName,
		Avatar:        nil,
		Password:      r.Password,
		PickupAddress: r.PickupAddress,
	}

	if err := sellerCandidate.PrepareCreate(); err != nil {
		return nil, err
	}

	return sellerCandidate, nil
}

func (h *sellerHandlersHTTP) updateReqToSellerModel(updateCandidate *models.Seller, r *dto.SellerUpdateRequestDto) (*models.Seller, error) {

	if r.FirstName != nil {
		updateCandidate.FirstName = strings.TrimSpace(*r.FirstName)
	}
	if r.LastName != nil {
		updateCandidate.LastName = strings.TrimSpace(*r.LastName)
	}
	if r.PickupAddress != nil {
		updateCandidate.PickupAddress = strings.TrimSpace(*r.PickupAddress)
	}
	if r.Avatar != nil {
		avatar := strings.TrimSpace(*r.Avatar)
		updateCandidate.Avatar = &avatar
	}
	if r.Password != nil {
		updateCandidate.Password = *r.Password
		if err := updateCandidate.HashPassword(); err != nil {
			return nil, err
		}
	}

	return updateCandidate, nil
}
