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

	"github.com/dinorain/kalobranded/config"
	"github.com/dinorain/kalobranded/internal/brand/delivery/http/dto"
	"github.com/dinorain/kalobranded/internal/brand/mock"
	"github.com/dinorain/kalobranded/internal/middlewares"
	"github.com/dinorain/kalobranded/internal/models"
	mockSessUC "github.com/dinorain/kalobranded/internal/session/mock"
	"github.com/dinorain/kalobranded/pkg/converter"
	"github.com/dinorain/kalobranded/pkg/logger"
)

func TestBrandsHandler_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	brandUC := mock.NewMockBrandUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	appLogger := logger.NewAppLogger(nil)
	mw := middlewares.NewMiddlewareManager(appLogger, nil)

	e := echo.New()
	v := validator.New()
	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	handlers := NewBrandHandlersHTTP(e.Group("brand"), appLogger, cfg, mw, v, brandUC, sessUC)

	reqDto := &dto.BrandRegisterRequestDto{
		BrandName:     "BrandName",
		PickupAddress: "PickupAddress",
	}

	buf := &bytes.Buffer{}
	_ = json.NewEncoder(buf).Encode(reqDto)

	req := httptest.NewRequest(http.MethodPost, "/brand", buf)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	resDto := &dto.BrandRegisterResponseDto{
		BrandID: uuid.Nil,
	}

	buf, _ = converter.AnyToBytesBuffer(resDto)

	brandUC.EXPECT().Register(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Brand{}, nil)
	require.NoError(t, handlers.Register()(ctx))
	require.Equal(t, http.StatusCreated, res.Code)
	require.Equal(t, buf.String(), res.Body.String())
}

func TestBrandsHandler_FindAll(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	brandUC := mock.NewMockBrandUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	appLogger := logger.NewAppLogger(nil)
	mw := middlewares.NewMiddlewareManager(appLogger, nil)

	e := echo.New()
	v := validator.New()
	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	handlers := NewBrandHandlersHTTP(e.Group("brand"), appLogger, cfg, mw, v, brandUC, sessUC)

	req := httptest.NewRequest(http.MethodGet, "/brand", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	var brands []models.Brand
	brands = append(brands, models.Brand{
		BrandID:       uuid.New(),
		BrandName:     "BrandName",
		PickupAddress: "PickupAddress",
	})

	brandUC.EXPECT().FindAll(gomock.Any(), gomock.Any()).AnyTimes().Return(brands, nil)
	require.NoError(t, handlers.FindAll()(ctx))
	require.Equal(t, http.StatusOK, res.Code)
}

func TestBrandsHandler_FindById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	brandUC := mock.NewMockBrandUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, nil)

	e := echo.New()
	v := validator.New()
	handlers := NewBrandHandlersHTTP(e.Group("brand"), appLogger, cfg, mw, v, brandUC, sessUC)

	req := httptest.NewRequest(http.MethodGet, "/brand/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	ctx.SetParamNames("id")
	ctx.SetParamValues("2ceba62a-35f4-444b-a358-4b14834837e1")

	brandUC.EXPECT().CachedFindById(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Brand{}, nil)
	require.NoError(t, handlers.FindById()(ctx))
	require.Equal(t, http.StatusOK, res.Code)
}

func TestBrandsHandler_UpdateById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	brandUC := mock.NewMockBrandUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, cfg)

	e := echo.New()
	e.Use(middleware.JWT([]byte("secret")))
	v := validator.New()
	handlers := NewBrandHandlersHTTP(e.Group("brand"), appLogger, cfg, mw, v, brandUC, sessUC)

	change := "changed"
	reqDto := &dto.BrandUpdateRequestDto{
		BrandName:     &change,
		PickupAddress: &change,
		Logo:          &change,
	}

	buf := &bytes.Buffer{}
	_ = json.NewEncoder(buf).Encode(reqDto)

	brandUUID := uuid.New()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["session_id"] = uuid.New().String()
	claims["brand_id"] = brandUUID.String()
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	validToken, _ := token.SignedString([]byte("secret"))

	req := httptest.NewRequest(http.MethodPost, "/brand/:id", buf)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("bearer %v", validToken))

	t.Run("Forbidden update by other brand", func(t *testing.T) {
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

		brandUC.EXPECT().UpdateById(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Brand{BrandID: brandUUID}, nil)

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
		ctx.SetParamValues(brandUUID.String())

		brandUC.EXPECT().UpdateById(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Brand{BrandID: brandUUID}, nil)
		brandUC.EXPECT().FindById(gomock.Any(), brandUUID).AnyTimes().Return(&models.Brand{BrandID: brandUUID}, nil)

		require.NoError(t, h(ctx))
		require.Equal(t, http.StatusOK, res.Code)
	})
}

func TestBrandsHandler_DeleteById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	brandUC := mock.NewMockBrandUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, nil)

	e := echo.New()
	v := validator.New()
	handlers := NewBrandHandlersHTTP(e.Group("brand"), appLogger, cfg, mw, v, brandUC, sessUC)

	req := httptest.NewRequest(http.MethodDelete, "/brand/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	brandUUID := uuid.New()
	ctx.SetParamNames("id")
	ctx.SetParamValues(brandUUID.String())

	brandUC.EXPECT().DeleteById(gomock.Any(), brandUUID).AnyTimes().Return(nil)
	require.NoError(t, handlers.DeleteById()(ctx))
	require.Equal(t, http.StatusOK, res.Code)
}
