package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/require"

	"github.com/dinorain/checkoutaja/config"
	"github.com/dinorain/checkoutaja/internal/middlewares"
	"github.com/dinorain/checkoutaja/internal/models"
	"github.com/dinorain/checkoutaja/internal/seller/delivery/http/dto"
	"github.com/dinorain/checkoutaja/internal/seller/mock"
	mockSessUC "github.com/dinorain/checkoutaja/internal/session/mock"
	"github.com/dinorain/checkoutaja/pkg/converter"
	"github.com/dinorain/checkoutaja/pkg/logger"
)

func TestSellersHandler_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sellerUC := mock.NewMockSellerUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	appLogger := logger.NewAppLogger(nil)
	mw := middlewares.NewMiddlewareManager(appLogger, nil)

	e := echo.New()
	v := validator.New()
	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	handlers := NewSellerHandlersHTTP(e.Group("seller"), appLogger, cfg, mw, v, sellerUC, sessUC)

	reqDto := &dto.SellerRegisterRequestDto{
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Password:      "123456",
		PickupAddress: "PickupAddress",
	}

	buf := &bytes.Buffer{}
	_ = json.NewEncoder(buf).Encode(reqDto)

	req := httptest.NewRequest(http.MethodPost, "/seller", buf)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	resDto := &dto.SellerRegisterResponseDto{
		SellerID: uuid.Nil,
	}

	buf, _ = converter.AnyToBytesBuffer(resDto)

	sellerUC.EXPECT().Register(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Seller{}, nil)
	require.NoError(t, handlers.Register()(ctx))
	require.Equal(t, http.StatusCreated, res.Code)
	require.Equal(t, buf.String(), res.Body.String())
}

func TestSellersHandler_Login(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sellerUC := mock.NewMockSellerUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	appLogger := logger.NewAppLogger(nil)
	mw := middlewares.NewMiddlewareManager(appLogger, nil)

	e := echo.New()
	v := validator.New()
	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	handlers := NewSellerHandlersHTTP(e.Group("seller"), appLogger, cfg, mw, v, sellerUC, sessUC)

	reqDto := &dto.SellerLoginRequestDto{
		Email:    "email@gmail.com",
		Password: "123456",
	}

	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(reqDto)

	req := httptest.NewRequest(http.MethodPost, "/seller/login", &buf)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	mockSeller := &models.Seller{
		SellerID:      uuid.New(),
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Password:      "123456",
		PickupAddress: "PickupAddress",
	}

	sellerUC.EXPECT().Login(gomock.Any(), reqDto.Email, reqDto.Password).AnyTimes().Return(mockSeller, nil)
	sessUC.EXPECT().CreateSession(gomock.Any(), &models.Session{UserID: mockSeller.SellerID}, cfg.Session.Expire).AnyTimes().Return("s", nil)
	sellerUC.EXPECT().GenerateTokenPair(gomock.Any(), gomock.Any()).AnyTimes().Return("rt", "at", nil)
	require.NoError(t, handlers.Login()(ctx))
	require.Equal(t, http.StatusCreated, res.Code)
}

func TestSellersHandler_FindAll(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sellerUC := mock.NewMockSellerUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	appLogger := logger.NewAppLogger(nil)
	mw := middlewares.NewMiddlewareManager(appLogger, nil)

	e := echo.New()
	v := validator.New()
	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	handlers := NewSellerHandlersHTTP(e.Group("seller"), appLogger, cfg, mw, v, sellerUC, sessUC)

	req := httptest.NewRequest(http.MethodGet, "/seller", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	var sellers []models.Seller
	sellers = append(sellers, models.Seller{
		SellerID:      uuid.New(),
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Password:      "123456",
		PickupAddress: "PickupAddress",
	})

	sellerUC.EXPECT().FindAll(gomock.Any(), gomock.Any()).AnyTimes().Return(sellers, nil)
	require.NoError(t, handlers.FindAll()(ctx))
	require.Equal(t, http.StatusOK, res.Code)
}

func TestSellersHandler_FindById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sellerUC := mock.NewMockSellerUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, nil)

	e := echo.New()
	v := validator.New()
	handlers := NewSellerHandlersHTTP(e.Group("seller"), appLogger, cfg, mw, v, sellerUC, sessUC)

	req := httptest.NewRequest(http.MethodGet, "/seller/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	ctx.SetParamNames("id")
	ctx.SetParamValues("2ceba62a-35f4-444b-a358-4b14834837e1")

	sellerUC.EXPECT().CachedFindById(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Seller{}, nil)
	require.NoError(t, handlers.FindById()(ctx))
	require.Equal(t, http.StatusOK, res.Code)
}

func TestSellersHandler_UpdateById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sellerUC := mock.NewMockSellerUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, cfg)

	e := echo.New()
	e.Use(middleware.JWT([]byte("secret")))
	v := validator.New()
	handlers := NewSellerHandlersHTTP(e.Group("seller"), appLogger, cfg, mw, v, sellerUC, sessUC)

	change := "changed"
	reqDto := &dto.SellerUpdateRequestDto{
		FirstName:     &change,
		LastName:      &change,
		Password:      &change,
		PickupAddress: &change,
		Avatar:        &change,
	}

	buf := &bytes.Buffer{}
	_ = json.NewEncoder(buf).Encode(reqDto)

	sellerUUID := uuid.New()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["session_id"] = uuid.New().String()
	claims["seller_id"] = sellerUUID.String()
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	validToken, _ := token.SignedString([]byte("secret"))

	req := httptest.NewRequest(http.MethodPost, "/seller/:id", buf)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("bearer %v", validToken))

	t.Run("Forbidden update by other seller", func(t *testing.T) {
		t.Parallel()

		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)

		handler := handlers.UpdateById()
		h := middleware.JWTWithConfig(middleware.JWTConfig{
			Claims:     claims,
			SigningKey: []byte("secret"),
		})(handler)

		ctx.SetParamNames("id")
		ctx.SetParamValues("2ceba62a-35f4-444b-a358-4b14834837e1")

		sellerUC.EXPECT().UpdateById(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Seller{SellerID: sellerUUID}, nil)

		require.NoError(t, h(ctx))
		require.Equal(t, http.StatusForbidden, res.Code)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)

		handler := handlers.UpdateById()
		h := middleware.JWTWithConfig(middleware.JWTConfig{
			Claims:     claims,
			SigningKey: []byte("secret"),
		})(handler)

		ctx.SetParamNames("id")
		ctx.SetParamValues(sellerUUID.String())

		sellerUC.EXPECT().UpdateById(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Seller{SellerID: sellerUUID}, nil)
		sellerUC.EXPECT().FindById(gomock.Any(), sellerUUID).AnyTimes().Return(&models.Seller{SellerID: sellerUUID}, nil)

		require.NoError(t, h(ctx))
		require.Equal(t, http.StatusOK, res.Code)
	})
}

func TestSellersHandler_DeleteById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sellerUC := mock.NewMockSellerUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, nil)

	e := echo.New()
	v := validator.New()
	handlers := NewSellerHandlersHTTP(e.Group("seller"), appLogger, cfg, mw, v, sellerUC, sessUC)

	req := httptest.NewRequest(http.MethodDelete, "/seller/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	sellerUUID := uuid.New()
	ctx.SetParamNames("id")
	ctx.SetParamValues(sellerUUID.String())

	sellerUC.EXPECT().DeleteById(gomock.Any(), sellerUUID).AnyTimes().Return(nil)
	require.NoError(t, handlers.DeleteById()(ctx))
	require.Equal(t, http.StatusOK, res.Code)
}

func TestSellersHandler_GetMe(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sellerUC := mock.NewMockSellerUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, cfg)

	e := echo.New()
	e.Use(middleware.JWT([]byte("secret")))
	v := validator.New()
	handlers := NewSellerHandlersHTTP(e.Group("seller"), appLogger, cfg, mw, v, sellerUC, sessUC)

	sellerUUID := uuid.New()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["session_id"] = uuid.New().String()
	claims["seller_id"] = sellerUUID.String()
	claims["role"] = "seller"
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	validToken, _ := token.SignedString([]byte("secret"))

	req := httptest.NewRequest(http.MethodPost, "/seller/logout", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("bearer %v", validToken))

	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	handler := handlers.GetMe()
	h := middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     claims,
		SigningKey: []byte("secret"),
	})(handler)

	sessUC.EXPECT().GetSessionById(gomock.Any(), claims["session_id"].(string)).AnyTimes().Return(&models.Session{}, nil)
	sellerUC.EXPECT().CachedFindById(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Seller{}, nil)

	require.NoError(t, h(ctx))
	require.Equal(t, http.StatusOK, res.Code)
}

func TestSellersHandler_Logout(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sellerUC := mock.NewMockSellerUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, cfg)

	e := echo.New()
	e.Use(middleware.JWT([]byte("secret")))
	v := validator.New()
	handlers := NewSellerHandlersHTTP(e.Group("seller"), appLogger, cfg, mw, v, sellerUC, sessUC)

	sellerUUID := uuid.New()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["session_id"] = uuid.New().String()
	claims["seller_id"] = sellerUUID.String()
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	validToken, _ := token.SignedString([]byte("secret"))

	req := httptest.NewRequest(http.MethodPost, "/seller/logout", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("bearer %v", validToken))

	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	handler := handlers.Logout()
	h := middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     claims,
		SigningKey: []byte("secret"),
	})(handler)

	sessUC.EXPECT().DeleteById(gomock.Any(), claims["session_id"].(string)).AnyTimes().Return(nil)

	require.NoError(t, h(ctx))
	require.Equal(t, http.StatusOK, res.Code)
}

func TestSellersHandler_RefreshToken(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sellerUC := mock.NewMockSellerUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}, Server: config.ServerConfig{JwtSecretKey: "secret"}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, cfg)

	e := echo.New()
	v := validator.New()
	handlers := NewSellerHandlersHTTP(e.Group("seller"), appLogger, cfg, mw, v, sellerUC, sessUC)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["session_id"] = uuid.New().String()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	validToken, _ := token.SignedString([]byte("secret"))

	reqDto := &dto.SellerRefreshTokenDto{
		RefreshToken: validToken,
	}

	buf := &bytes.Buffer{}
	_ = json.NewEncoder(buf).Encode(reqDto)

	req := httptest.NewRequest(http.MethodPost, "/seller/refresh", buf)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	sessUC.EXPECT().GetSessionById(gomock.Any(), claims["session_id"].(string)).AnyTimes().Return(&models.Session{}, nil)
	sellerUC.EXPECT().FindById(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Seller{}, nil)
	sellerUC.EXPECT().GenerateTokenPair(gomock.Any(), gomock.Any()).AnyTimes().Return("rt", "at", nil)

	require.NoError(t, handlers.RefreshToken()(ctx))
	require.Equal(t, http.StatusOK, res.Code)
}
