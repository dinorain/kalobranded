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
	"github.com/dinorain/kalobranded/internal/middlewares"
	"github.com/dinorain/kalobranded/internal/models"
	"github.com/dinorain/kalobranded/internal/product/delivery/http/dto"
	"github.com/dinorain/kalobranded/internal/product/mock"
	mockSessUC "github.com/dinorain/kalobranded/internal/session/mock"
	"github.com/dinorain/kalobranded/pkg/converter"
	"github.com/dinorain/kalobranded/pkg/logger"
)

func TestProductsHandler_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	productUC := mock.NewMockProductUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	appLogger := logger.NewAppLogger(nil)
	mw := middlewares.NewMiddlewareManager(appLogger, nil)

	e := echo.New()
	e.Use(middleware.JWT([]byte("secret")))
	v := validator.New()
	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	handlers := NewProductHandlersHTTP(e.Group("product"), appLogger, cfg, mw, v, productUC, sessUC)

	reqDto := &dto.ProductCreateRequestDto{
		Name:        "Name",
		Description: "Description",
		Price:       10000.0,
	}

	buf := &bytes.Buffer{}
	_ = json.NewEncoder(buf).Encode(reqDto)

	brandUUID := uuid.New()
	sessUUID := uuid.New()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["session_id"] = sessUUID.String()
	claims["brand_id"] = brandUUID.String()
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	validToken, _ := token.SignedString([]byte("secret"))

	req := httptest.NewRequest(http.MethodPost, "/product", buf)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("bearer %v", validToken))
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	resDto := &dto.ProductCreateResponseDto{
		ProductID: uuid.Nil,
	}

	buf, _ = converter.AnyToBytesBuffer(resDto)

	handler := handlers.Create()
	h := middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     claims,
		SigningKey: []byte("secret"),
	})(handler)

	sessUC.EXPECT().GetSessionById(gomock.Any(), claims["session_id"]).AnyTimes().Return(&models.Session{UserID: brandUUID, SessionID: sessUUID.String()}, nil)
	productUC.EXPECT().Create(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Product{BrandID: brandUUID}, nil)
	require.NoError(t, h(ctx))
	require.Equal(t, http.StatusCreated, res.Code)
	require.Equal(t, buf.String(), res.Body.String())
}

func TestProductsHandler_FindAll(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	productUC := mock.NewMockProductUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	appLogger := logger.NewAppLogger(nil)
	mw := middlewares.NewMiddlewareManager(appLogger, nil)

	e := echo.New()
	e.Use(middleware.JWT([]byte("secret")))
	v := validator.New()
	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	handlers := NewProductHandlersHTTP(e.Group("product"), appLogger, cfg, mw, v, productUC, sessUC)

	brandUUID := uuid.New()

	var products []models.Product
	var oneOnly []models.Product

	m := models.Product{
		ProductID:   uuid.New(),
		Name:        "Name",
		Description: "FirstName",
		Price:       10000.00,
		BrandID:    brandUUID,
	}
	oneOnly = append(oneOnly, m)
	products = append(products, m, m)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	sessUUID := uuid.New()
	claims["session_id"] = sessUUID.String()
	claims["brand_id"] = brandUUID.String()
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	validToken, _ := token.SignedString([]byte("secret"))

	sessUC.EXPECT().GetSessionById(gomock.Any(), claims["session_id"]).AnyTimes().Return(&models.Session{UserID: brandUUID, SessionID: sessUUID.String()}, nil)
	productUC.EXPECT().FindAllByBrandId(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(oneOnly, nil)

	t.Run("Filtered view when accessed by brand", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/product", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("bearer %v", validToken))
		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)

		handler := handlers.FindAll()
		h := middleware.JWTWithConfig(middleware.JWTConfig{
			Claims:     claims,
			SigningKey: []byte("secret"),
		})(handler)

		require.NoError(t, h(ctx))
		require.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("All products view when accessed by user", func(t *testing.T) {

		claims["session_id"] = sessUUID.String()
		claims["user_id"] = brandUUID.String()
		claims["role"] = models.UserRoleUser
		claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
		validToken, _ := token.SignedString([]byte("secret"))

		req := httptest.NewRequest(http.MethodGet, "/product", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("bearer %v", validToken))
		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)

		handler := handlers.FindAll()
		h := middleware.JWTWithConfig(middleware.JWTConfig{
			Claims:     claims,
			SigningKey: []byte("secret"),
		})(handler)

		sessUC.EXPECT().GetSessionById(gomock.Any(), claims["session_id"]).AnyTimes().Return(&models.Session{UserID: brandUUID, SessionID: sessUUID.String()}, nil)
		productUC.EXPECT().FindAll(gomock.Any(), gomock.Any()).AnyTimes().Return(products, nil)

		require.NoError(t, h(ctx))
		require.Equal(t, http.StatusOK, res.Code)

		resDto := &dto.ProductFindResponseDto{}
		_ = json.Unmarshal(res.Body.Bytes(), resDto)

		require.Equal(t, len(resDto.Data.([]interface{})), len(products))
	})
}

func TestProductsHandler_FindById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	productUC := mock.NewMockProductUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, nil)

	e := echo.New()
	e.Use(middleware.JWT([]byte("secret")))
	v := validator.New()
	handlers := NewProductHandlersHTTP(e.Group("product"), appLogger, cfg, mw, v, productUC, sessUC)

	req := httptest.NewRequest(http.MethodGet, "/product/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	ctx.SetParamNames("id")
	ctx.SetParamValues("2ceba62a-35f4-444b-a358-4b14834837e1")

	productUC.EXPECT().CachedFindById(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Product{}, nil)
	require.NoError(t, handlers.FindById()(ctx))
	require.Equal(t, http.StatusOK, res.Code)
}

func TestProductsHandler_UpdateById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	productUC := mock.NewMockProductUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, cfg)

	e := echo.New()
	e.Use(middleware.JWT([]byte("secret")))
	v := validator.New()
	handlers := NewProductHandlersHTTP(e.Group("product"), appLogger, cfg, mw, v, productUC, sessUC)

	change := "changed"
	reqDto := &dto.ProductUpdateRequestDto{
		Name:        &change,
		Description: &change,
	}

	brandUUID := uuid.New()
	productUUID := uuid.New()
	sessUUID := uuid.New()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	sessUC.EXPECT().GetSessionById(gomock.Any(), claims["session_id"]).AnyTimes().Return(&models.Session{UserID: brandUUID, SessionID: sessUUID.String()}, nil)
	productUC.EXPECT().FindById(gomock.Any(), productUUID).AnyTimes().Return(&models.Product{BrandID: brandUUID}, nil)
	productUC.EXPECT().UpdateById(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Product{BrandID: brandUUID}, nil)

	t.Run("Success update by brand", func(t *testing.T) {

		claims["session_id"] = sessUUID.String()
		claims["brand_id"] = brandUUID.String()
		claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
		validToken, _ := token.SignedString([]byte("secret"))

		buf := &bytes.Buffer{}
		_ = json.NewEncoder(buf).Encode(reqDto)

		req := httptest.NewRequest(http.MethodPost, "/product/:id", buf)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("bearer %v", validToken))

		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)

		handler := handlers.UpdateById()
		h := middleware.JWTWithConfig(middleware.JWTConfig{
			Claims:     claims,
			SigningKey: []byte("secret"),
		})(handler)

		ctx.SetParamNames("id")
		ctx.SetParamValues(productUUID.String())

		require.NoError(t, h(ctx))
		require.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("Forbidden update by other brand", func(t *testing.T) {

		claims["session_id"] = sessUUID.String()
		claims["brand_id"] = sessUUID.String()
		claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
		validToken, _ := token.SignedString([]byte("secret"))

		buf := &bytes.Buffer{}
		_ = json.NewEncoder(buf).Encode(reqDto)

		req := httptest.NewRequest(http.MethodPost, "/product/:id", buf)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("bearer %v", validToken))

		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)

		handler := handlers.UpdateById()
		h := middleware.JWTWithConfig(middleware.JWTConfig{
			Claims:     claims,
			SigningKey: []byte("secret"),
		})(handler)

		ctx.SetParamNames("id")
		ctx.SetParamValues(productUUID.String())

		require.NoError(t, h(ctx))
		require.Equal(t, http.StatusForbidden, res.Code)
	})
}

func TestProductsHandler_DeleteById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	productUC := mock.NewMockProductUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, cfg)

	e := echo.New()
	e.Use(middleware.JWT([]byte("secret")))
	v := validator.New()
	handlers := NewProductHandlersHTTP(e.Group("product"), appLogger, cfg, mw, v, productUC, sessUC)

	brandUUID := uuid.New()
	productUUID := uuid.New()
	sessUUID := uuid.New()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	productUC.EXPECT().CachedFindById(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Product{BrandID: brandUUID}, nil)
	productUC.EXPECT().DeleteById(gomock.Any(), productUUID).AnyTimes().Return(nil)

	t.Run("Success update by brand", func(t *testing.T) {

		claims["session_id"] = sessUUID.String()
		claims["brand_id"] = brandUUID.String()
		claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
		validToken, _ := token.SignedString([]byte("secret"))

		req := httptest.NewRequest(http.MethodPost, "/product/:id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("bearer %v", validToken))

		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)

		handler := handlers.DeleteById()
		h := middleware.JWTWithConfig(middleware.JWTConfig{
			Claims:     claims,
			SigningKey: []byte("secret"),
		})(handler)

		ctx.SetParamNames("id")
		ctx.SetParamValues(productUUID.String())

		require.NoError(t, h(ctx))
		require.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("Forbidden update by other brand", func(t *testing.T) {

		claims["session_id"] = sessUUID.String()
		claims["brand_id"] = sessUUID.String()
		claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
		validToken, _ := token.SignedString([]byte("secret"))

		req := httptest.NewRequest(http.MethodPost, "/product/:id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("bearer %v", validToken))

		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)

		handler := handlers.DeleteById()
		h := middleware.JWTWithConfig(middleware.JWTConfig{
			Claims:     claims,
			SigningKey: []byte("secret"),
		})(handler)

		ctx.SetParamNames("id")
		ctx.SetParamValues(productUUID.String())

		require.NoError(t, h(ctx))
		require.Equal(t, http.StatusForbidden, res.Code)
	})
}
