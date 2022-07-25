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
	"github.com/dinorain/checkoutaja/internal/order/delivery/http/dto"
	"github.com/dinorain/checkoutaja/internal/order/mock"
	mockProductUC "github.com/dinorain/checkoutaja/internal/product/mock"
	mockSellerUC "github.com/dinorain/checkoutaja/internal/seller/mock"
	mockSessUC "github.com/dinorain/checkoutaja/internal/session/mock"
	mockUserUC "github.com/dinorain/checkoutaja/internal/user/mock"
	"github.com/dinorain/checkoutaja/pkg/converter"
	"github.com/dinorain/checkoutaja/pkg/logger"
)

func TestOrdersHandler_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderUC := mock.NewMockOrderUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)
	userUC := mockUserUC.NewMockUserUseCase(ctrl)
	sellerUC := mockSellerUC.NewMockSellerUseCase(ctrl)
	productUC := mockProductUC.NewMockProductUseCase(ctrl)

	appLogger := logger.NewAppLogger(nil)
	mw := middlewares.NewMiddlewareManager(appLogger, nil)

	e := echo.New()
	e.Use(middleware.JWT([]byte("secret")))
	v := validator.New()
	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	handlers := NewOrderHandlersHTTP(e.Group("order"), appLogger, cfg, mw, v, orderUC, userUC, sellerUC, productUC, sessUC)

	productUUID := uuid.New()
	reqDto := &dto.OrderCreateRequestDto{
		ProductID: productUUID,
		Quantity:  1,
	}

	buf := &bytes.Buffer{}
	_ = json.NewEncoder(buf).Encode(reqDto)

	userUUID := uuid.New()
	sellerUUID := uuid.New()
	sessUUID := uuid.New()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["session_id"] = sessUUID.String()
	claims["user_id"] = userUUID.String()
	claims["role"] = models.UserRoleUser
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	validToken, _ := token.SignedString([]byte("secret"))

	req := httptest.NewRequest(http.MethodPost, "/order", buf)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("bearer %v", validToken))
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	resDto := &dto.OrderCreateResponseDto{
		OrderID: uuid.Nil,
	}

	buf, _ = converter.AnyToBytesBuffer(resDto)

	handler := handlers.Create()
	h := middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     claims,
		SigningKey: []byte("secret"),
	})(handler)

	sessUC.EXPECT().GetSessionById(gomock.Any(), claims["session_id"]).AnyTimes().Return(&models.Session{UserID: userUUID, SessionID: sessUUID.String()}, nil)
	userUC.EXPECT().CachedFindById(gomock.Any(), userUUID).AnyTimes().Return(&models.User{UserID: userUUID}, nil)
	productUC.EXPECT().CachedFindById(gomock.Any(), productUUID).AnyTimes().Return(&models.Product{ProductID: productUUID, SellerID: sellerUUID}, nil)
	sellerUC.EXPECT().CachedFindById(gomock.Any(), sellerUUID).AnyTimes().Return(&models.Seller{SellerID: sellerUUID}, nil)
	orderUC.EXPECT().Create(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Order{Item: models.OrderItem{ProductID: productUUID}}, nil)
	require.NoError(t, h(ctx))
	require.Equal(t, http.StatusCreated, res.Code)
	require.Equal(t, buf.String(), res.Body.String())
}

func TestOrdersHandler_FindAll(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderUC := mock.NewMockOrderUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)
	userUC := mockUserUC.NewMockUserUseCase(ctrl)
	sellerUC := mockSellerUC.NewMockSellerUseCase(ctrl)
	productUC := mockProductUC.NewMockProductUseCase(ctrl)

	appLogger := logger.NewAppLogger(nil)
	mw := middlewares.NewMiddlewareManager(appLogger, nil)

	e := echo.New()
	e.Use(middleware.JWT([]byte("secret")))
	v := validator.New()
	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	handlers := NewOrderHandlersHTTP(e.Group("order"), appLogger, cfg, mw, v, orderUC, userUC, sellerUC, productUC, sessUC)

	userUUID := uuid.New()
	sellerUUID := uuid.New()
	productUUID := uuid.New()

	var orders []models.Order

	m := models.Order{
		OrderID:  uuid.New(),
		UserID:   userUUID,
		SellerID: sellerUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			SellerID:    sellerUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}
	orders = append(orders, m, m)

	userOrders := make([]models.Order, len(orders)+1)
	copy(userOrders, orders)

	m.SellerID = uuid.New()
	userOrders = append(userOrders, m)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	sessUUID := uuid.New()
	claims["session_id"] = sessUUID.String()
	claims["seller_id"] = sellerUUID.String()
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	validToken, _ := token.SignedString([]byte("secret"))

	sessUC.EXPECT().GetSessionById(gomock.Any(), claims["session_id"]).AnyTimes().Return(&models.Session{UserID: sellerUUID, SessionID: sessUUID.String()}, nil)
	orderUC.EXPECT().FindAllBySellerId(gomock.Any(), sellerUUID, gomock.Any()).AnyTimes().Return(orders, nil)
	orderUC.EXPECT().FindAllByUserId(gomock.Any(), userUUID, gomock.Any()).AnyTimes().Return(userOrders, nil)
	orderUC.EXPECT().FindAll(gomock.Any(), gomock.Any()).AnyTimes().Return(userOrders, nil)

	t.Run("Filtered view when accessed by seller", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/order", nil)
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

		resDto := &dto.OrderFindResponseDto{}
		_ = json.Unmarshal(res.Body.Bytes(), resDto)

		require.Equal(t, len(resDto.Data.([]interface{})), len(orders))
	})

	t.Run("All orders view when accessed by user", func(t *testing.T) {

		claims["session_id"] = sessUUID.String()
		claims["user_id"] = userUUID.String()
		claims["role"] = models.UserRoleUser
		claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
		validToken, _ := token.SignedString([]byte("secret"))

		req := httptest.NewRequest(http.MethodGet, "/order", nil)
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

		resDto := &dto.OrderFindResponseDto{}
		_ = json.Unmarshal(res.Body.Bytes(), resDto)

		require.Equal(t, len(resDto.Data.([]interface{})), len(userOrders))
	})

	t.Run("All orders view when accessed by admin", func(t *testing.T) {

		claims["session_id"] = sessUUID.String()
		claims["user_id"] = userUUID.String()
		claims["role"] = models.UserRoleAdmin
		claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
		validToken, _ := token.SignedString([]byte("secret"))

		req := httptest.NewRequest(http.MethodGet, "/order", nil)
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

		resDto := &dto.OrderFindResponseDto{}
		_ = json.Unmarshal(res.Body.Bytes(), resDto)

		require.Equal(t, len(resDto.Data.([]interface{})), len(userOrders))
	})
}

func TestOrdersHandler_FindById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderUC := mock.NewMockOrderUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)
	userUC := mockUserUC.NewMockUserUseCase(ctrl)
	sellerUC := mockSellerUC.NewMockSellerUseCase(ctrl)
	productUC := mockProductUC.NewMockProductUseCase(ctrl)

	appLogger := logger.NewAppLogger(nil)
	mw := middlewares.NewMiddlewareManager(appLogger, nil)

	e := echo.New()
	e.Use(middleware.JWT([]byte("secret")))
	v := validator.New()
	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	handlers := NewOrderHandlersHTTP(e.Group("order"), appLogger, cfg, mw, v, orderUC, userUC, sellerUC, productUC, sessUC)

	req := httptest.NewRequest(http.MethodGet, "/order/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	ctx.SetParamNames("id")
	ctx.SetParamValues("2ceba62a-35f4-444b-a358-4b14834837e1")

	orderUC.EXPECT().CachedFindById(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Order{}, nil)
	require.NoError(t, handlers.FindById()(ctx))
	require.Equal(t, http.StatusOK, res.Code)
}

func TestOrdersHandler_AcceptById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderUC := mock.NewMockOrderUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)
	userUC := mockUserUC.NewMockUserUseCase(ctrl)
	sellerUC := mockSellerUC.NewMockSellerUseCase(ctrl)
	productUC := mockProductUC.NewMockProductUseCase(ctrl)

	appLogger := logger.NewAppLogger(nil)
	mw := middlewares.NewMiddlewareManager(appLogger, nil)

	e := echo.New()
	e.Use(middleware.JWT([]byte("secret")))
	v := validator.New()
	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	handlers := NewOrderHandlersHTTP(e.Group("order"), appLogger, cfg, mw, v, orderUC, userUC, sellerUC, productUC, sessUC)

	sellerUUID := uuid.New()
	orderUUID := uuid.New()
	sessUUID := uuid.New()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	sessUC.EXPECT().GetSessionById(gomock.Any(), claims["session_id"]).AnyTimes().Return(&models.Session{UserID: sellerUUID, SessionID: sessUUID.String()}, nil)
	orderUC.EXPECT().FindById(gomock.Any(), orderUUID).AnyTimes().Return(&models.Order{SellerID: sellerUUID}, nil)
	orderUC.EXPECT().UpdateById(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Order{SellerID: sellerUUID}, nil)

	t.Run("Success update by seller", func(t *testing.T) {

		claims["session_id"] = sessUUID.String()
		claims["seller_id"] = sellerUUID.String()
		claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
		validToken, _ := token.SignedString([]byte("secret"))

		req := httptest.NewRequest(http.MethodPost, "/order/:id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("bearer %v", validToken))

		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)

		handler := handlers.AcceptById()
		h := middleware.JWTWithConfig(middleware.JWTConfig{
			Claims:     claims,
			SigningKey: []byte("secret"),
		})(handler)

		ctx.SetParamNames("id")
		ctx.SetParamValues(orderUUID.String())

		require.NoError(t, h(ctx))
		require.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("Forbidden update by other seller", func(t *testing.T) {

		claims["session_id"] = sessUUID.String()
		claims["seller_id"] = sessUUID.String()
		claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
		validToken, _ := token.SignedString([]byte("secret"))

		req := httptest.NewRequest(http.MethodPost, "/order/:id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("bearer %v", validToken))

		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)

		handler := handlers.AcceptById()
		h := middleware.JWTWithConfig(middleware.JWTConfig{
			Claims:     claims,
			SigningKey: []byte("secret"),
		})(handler)

		ctx.SetParamNames("id")
		ctx.SetParamValues(orderUUID.String())

		require.NoError(t, h(ctx))
		require.Equal(t, http.StatusForbidden, res.Code)
	})
}
