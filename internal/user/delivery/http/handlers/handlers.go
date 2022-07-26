package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"github.com/dinorain/kalobranded/config"
	"github.com/dinorain/kalobranded/internal/middlewares"
	"github.com/dinorain/kalobranded/internal/models"
	"github.com/dinorain/kalobranded/internal/session"
	"github.com/dinorain/kalobranded/internal/user"
	"github.com/dinorain/kalobranded/internal/user/delivery/http/dto"
	"github.com/dinorain/kalobranded/pkg/constants"
	httpErrors "github.com/dinorain/kalobranded/pkg/http_errors"
	"github.com/dinorain/kalobranded/pkg/logger"
	"github.com/dinorain/kalobranded/pkg/utils"
)

type userHandlersHTTP struct {
	mux    *http.ServeMux
	logger logger.Logger
	cfg    *config.Config
	mw     middlewares.MiddlewareManager
	v      *validator.Validate
	userUC user.UserUseCase
	sessUC session.SessUseCase
}

var _ user.UserHandlers = (*userHandlersHTTP)(nil)

func NewUserHandlersHTTP(
	mux *http.ServeMux,
	logger logger.Logger,
	cfg *config.Config,
	mw middlewares.MiddlewareManager,
	v *validator.Validate,
	userUC user.UserUseCase,
	sessUC session.SessUseCase,
) *userHandlersHTTP {
	return &userHandlersHTTP{mux: mux, logger: logger, cfg: cfg, mw: mw, v: v, userUC: userUC, sessUC: sessUC}
}

// Register
// @Tags Users
// @Summary Register user
// @Description Create user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param payload body dto.UserRegisterRequestDto true "Payload"
// @Success 200 {object} dto.UserRegisterResponseDto
// @Router /user/create [post]
func (h *userHandlersHTTP) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	createDto := &dto.UserRegisterRequestDto{}
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

	user, err := h.registerReqToUserModel(createDto)

	if err != nil {
		h.logger.Errorf("registerReqToUserModel: %v", err)
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	createdUser, err := h.userUC.Register(ctx, user)
	if err != nil {
		h.logger.Errorf("userUC.Register: %v", err)
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	res, _ := json.Marshal(dto.UserRegisterResponseDto{UserID: createdUser.UserID})
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
	return
}

// Login
// @Tags Users
// @Summary User login
// @Description User login with email and password
// @Accept json
// @Produce json
// @Param payload body dto.UserLoginRequestDto true "Payload"
// @Success 200 {object} dto.UserLoginResponseDto
// @Router /user/login [post]
func (h *userHandlersHTTP) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	loginDto := &dto.UserLoginRequestDto{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginDto); err != nil {
		h.logger.Errorf("decoder.Decode: %v", err)
		_ = httpErrors.NewBadRequestError(w, err.Error(), h.cfg.Http.DebugErrorsResponse)
		return
	}

	if err := h.v.Struct(loginDto); err != nil {
		h.logger.Errorf("h.v.Struct: %v", err)
		_ = httpErrors.NewBadRequestError(w, err.Error(), h.cfg.Http.DebugErrorsResponse)
		return
	}

	email := loginDto.Email
	if !utils.ValidateEmail(email) {
		h.logger.Errorf("ValidateEmail: %v", email)
		_ = httpErrors.ErrorCtxResponse(w, errors.New("invalid email"), h.cfg.Http.DebugErrorsResponse)
		return
	}

	user, err := h.userUC.Login(ctx, email, loginDto.Password)
	if err != nil {
		h.logger.Errorf("userUC.Login: %v", email)
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	session, err := h.sessUC.CreateSession(ctx, &models.Session{
		UserID: user.UserID,
	}, h.cfg.Session.Expire)
	if err != nil {
		h.logger.Errorf("sessUC.CreateSession: %v", err)
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	accessToken, refreshToken, err := h.userUC.GenerateTokenPair(user, session)
	if err != nil {
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	res, _ := json.Marshal(dto.UserLoginResponseDto{UserID: user.UserID, Tokens: &dto.UserRefreshTokenResponseDto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}})
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return
}

// FindAll
// @Tags Users
// @Summary Find all users
// @Description Admin find all users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param size query string false "pagination size"
// @Param page query string false "pagination page"
// @Success 200 {object} dto.UserFindResponseDto
// @Router /user [get]
func (h *userHandlersHTTP) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	queryParam := r.URL.Query()

	if queryParam.Get("id") != "" {
		h.FindById(w, r)
	}

	pq := utils.NewPaginationFromQueryParams(queryParam.Get(constants.Size), queryParam.Get(constants.Page))
	users, err := h.userUC.FindAll(ctx, pq)
	if err != nil {
		h.logger.Errorf("userUC.FindAll: %v", err)
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	res, _ := json.Marshal(dto.UserFindResponseDto{
		Data: users,
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
// @Tags Users
// @Summary Find user by id
// @Description Find user by id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id query string true "user uuid"
// @Success 200 {object} dto.UserResponseDto
// @Router /user [get]
func (h *userHandlersHTTP) FindById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	queryParam := r.URL.Query()

	if queryParam.Get("id") == "" {
		_ = httpErrors.NewBadRequestError(w, nil, h.cfg.Http.DebugErrorsResponse)
		return
	}
	userUUID, err := uuid.Parse(queryParam.Get("id"))
	if err != nil {
		h.logger.WarnMsg("uuid.FromString", err)
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	user, err := h.userUC.CachedFindById(ctx, userUUID)
	if err != nil {
		h.logger.Errorf("userUC.CachedFindById: %v", err)
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	res, _ := json.Marshal(dto.UserResponseFromModel(user))
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return
}

// GetMe
// @Tags Users
// @Summary Find me
// @Description Get session id from token, find user by uuid and returns it
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.UserResponseDto
// @Router /user/me [get]
func (h *userHandlersHTTP) GetMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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

	res, _ := json.Marshal(dto.UserResponseFromModel(user))
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return
}

// Logout
// @Tags Users
// @Summary User logout
// @Description Delete current session
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} nil
// @Router /user/logout [post]
func (h *userHandlersHTTP) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessID, _, _, err := h.getSessionIDFromCtx(w, r)
	if err != nil {
		h.logger.Errorf("getSessionIDFromCtx: %v", err)
		return
	}

	if err := h.sessUC.DeleteById(ctx, sessID); err != nil {
		h.logger.Errorf("sessUC.DeleteById: %v", err)
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// RefreshToken
// @Tags Users
// @Summary Refresh access token
// @Description Refresh access token
// @Accept json
// @Produce json
// @Param payload body dto.UserRefreshTokenDto true "Payload"
// @Success 200 {object} dto.UserRefreshTokenResponseDto
// @Router /user/refresh [post]
func (h *userHandlersHTTP) RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	refreshTokenDto := &dto.UserRefreshTokenDto{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&refreshTokenDto); err != nil {
		h.logger.Errorf("decoder.Decode: %v", err)
		_ = httpErrors.NewBadRequestError(w, err.Error(), h.cfg.Http.DebugErrorsResponse)
		return
	}

	if err := h.v.Struct(refreshTokenDto); err != nil {
		h.logger.Errorf("h.v.Struct: %v", err)
		_ = httpErrors.NewBadRequestError(w, err.Error(), h.cfg.Http.DebugErrorsResponse)
		return
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
		_ = httpErrors.ErrorCtxResponse(w, errors.New("invalid refresh token"), h.cfg.Http.DebugErrorsResponse)
	}

	if !token.Valid {
		h.logger.Warnf("token.Valid")
		_ = httpErrors.ErrorCtxResponse(w, errors.New("invalid refresh token"), h.cfg.Http.DebugErrorsResponse)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		h.logger.Warnf("jwt.MapClaims: %+v", token.Claims)
		_ = httpErrors.ErrorCtxResponse(w, errors.New("invalid refresh token"), h.cfg.Http.DebugErrorsResponse)
		return
	}

	sessID, ok := claims["session_id"].(string)
	if !ok {
		h.logger.Warnf("session_id: %+v", claims)
		_ = httpErrors.ErrorCtxResponse(w, errors.New("invalid refresh token"), h.cfg.Http.DebugErrorsResponse)
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

	user, err := h.userUC.FindById(ctx, session.UserID)
	if err != nil {
		h.logger.Errorf("userUC.FindById: %v", err)
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	accessToken, refreshToken, err := h.userUC.GenerateTokenPair(user, sessID)
	if err != nil {
		_ = httpErrors.ErrorCtxResponse(w, err, h.cfg.Http.DebugErrorsResponse)
		return
	}

	res, _ := json.Marshal(dto.UserRefreshTokenResponseDto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return
}

func (h *userHandlersHTTP) getSessionIDFromCtx(w http.ResponseWriter, r *http.Request) (sessionID string, userID string, role string, err error) {
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

	userID, ok = claims["user_id"].(string)
	if !ok {
		h.logger.Warnf("user_id: %+v", claims)
		return "", "", "", errors.New("invalid token header")
	}

	role, ok = claims["role"].(string)
	if !ok {
		h.logger.Warnf("role: %+v", claims)
		return "", "", "", errors.New("invalid token header")
	}

	return sessionID, userID, role, nil
}

func (h *userHandlersHTTP) registerReqToUserModel(r *dto.UserRegisterRequestDto) (*models.User, error) {
	userCandidate := &models.User{
		Email:           r.Email,
		FirstName:       r.FirstName,
		LastName:        r.LastName,
		Role:            r.Role,
		Avatar:          nil,
		Password:        r.Password,
		DeliveryAddress: r.DeliveryAddress,
	}

	if err := userCandidate.PrepareCreate(); err != nil {
		return nil, err
	}

	return userCandidate, nil
}
