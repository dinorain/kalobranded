package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/dinorain/kalobranded/config"
	mockBrandUC "github.com/dinorain/kalobranded/internal/brand/mock"
	"github.com/dinorain/kalobranded/internal/middlewares"
	"github.com/dinorain/kalobranded/internal/models"
	"github.com/dinorain/kalobranded/internal/order/delivery/http/dto"
	"github.com/dinorain/kalobranded/internal/order/mock"
	mockProductUC "github.com/dinorain/kalobranded/internal/product/mock"
	mockSessUC "github.com/dinorain/kalobranded/internal/session/mock"
	mockUserUC "github.com/dinorain/kalobranded/internal/user/mock"
	"github.com/dinorain/kalobranded/pkg/converter"
	"github.com/dinorain/kalobranded/pkg/logger"
)

func TestOrdersHandler_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderUC := mock.NewMockOrderUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)
	userUC := mockUserUC.NewMockUserUseCase(ctrl)
	brandUC := mockBrandUC.NewMockBrandUseCase(ctrl)
	productUC := mockProductUC.NewMockProductUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}, Server: config.ServerConfig{JwtSecretKey: "secret"}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, cfg)

	v := validator.New()

	mux := http.NewServeMux()
	handlers := NewOrderHandlersHTTP(mux, appLogger, cfg, mw, v, orderUC, userUC, brandUC, productUC, sessUC)

	userUUID := uuid.New()
	brandUUID := uuid.New()
	orderUUID := uuid.New()
	sessUUID := uuid.New()
	productUUID := uuid.New()

	reqDto := &dto.OrderCreateRequestDto{
		ProductID: productUUID,
		Quantity:  1,
	}

	buf := &bytes.Buffer{}
	_ = json.NewEncoder(buf).Encode(reqDto)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["session_id"] = sessUUID.String()
	claims["user_id"] = userUUID.String()
	claims["role"] = models.UserRoleUser
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	validToken, _ := token.SignedString([]byte(cfg.Server.JwtSecretKey))

	req := httptest.NewRequest(http.MethodPost, "/order/create", buf)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", validToken))
	w := httptest.NewRecorder()

	wDto := &dto.OrderCreateResponseDto{
		OrderID: orderUUID,
	}

	buf, _ = converter.AnyToBytesBuffer(wDto)

	sessUC.EXPECT().GetSessionById(gomock.Any(), sessUUID.String()).AnyTimes().Return(&models.Session{UserID: userUUID, SessionID: sessUUID.String()}, nil)
	userUC.EXPECT().CachedFindById(gomock.Any(), userUUID).AnyTimes().Return(&models.User{UserID: userUUID}, nil)
	productUC.EXPECT().CachedFindById(gomock.Any(), productUUID).AnyTimes().Return(&models.Product{ProductID: productUUID, BrandID: brandUUID}, nil)
	brandUC.EXPECT().CachedFindById(gomock.Any(), brandUUID).AnyTimes().Return(&models.Brand{BrandID: brandUUID}, nil)
	orderUC.EXPECT().Create(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Order{OrderID: orderUUID}, nil)

	handler := http.HandlerFunc(handlers.Create)
	handler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	require.NotNil(t, data)
	require.Equal(t, http.StatusCreated, w.Code)
	require.Equal(t, strings.Trim(buf.String(), "\n"), string(data))
}

func TestOrdersHandler_Find(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderUC := mock.NewMockOrderUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)
	userUC := mockUserUC.NewMockUserUseCase(ctrl)
	brandUC := mockBrandUC.NewMockBrandUseCase(ctrl)
	productUC := mockProductUC.NewMockProductUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}, Server: config.ServerConfig{JwtSecretKey: "secret"}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, cfg)

	v := validator.New()

	mux := http.NewServeMux()
	handlers := NewOrderHandlersHTTP(mux, appLogger, cfg, mw, v, orderUC, userUC, brandUC, productUC, sessUC)

	userUUID := uuid.New()
	brandUUID := uuid.New()
	orderUUID := uuid.New()
	sessUUID := uuid.New()
	productUUID := uuid.New()

	var orders []models.Order
	var oneOnly []models.Order

	m := models.Order{
		OrderID: orderUUID,
		UserID:  userUUID,
		BrandID: brandUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			BrandID:     brandUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}
	oneOnly = append(oneOnly, m)
	orders = append(orders, m, models.Order{UserID: uuid.New()})

	t.Run("ByUserRole", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["session_id"] = sessUUID.String()
		claims["user_id"] = userUUID.String()
		claims["role"] = models.UserRoleUser
		claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
		validToken, _ := token.SignedString([]byte(cfg.Server.JwtSecretKey))

		req := httptest.NewRequest(http.MethodGet, "/order", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", validToken))
		w := httptest.NewRecorder()

		orderUC.EXPECT().FindAllByUserId(gomock.Any(), userUUID, gomock.Any()).AnyTimes().Return(oneOnly, nil)
		sessUC.EXPECT().GetSessionById(gomock.Any(), sessUUID.String()).AnyTimes().Return(&models.Session{UserID: userUUID, SessionID: sessUUID.String()}, nil)
		userUC.EXPECT().CachedFindById(gomock.Any(), userUUID).AnyTimes().Return(&models.User{UserID: userUUID}, nil)

		handler := http.HandlerFunc(handlers.FindAll)
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}
		require.NotNil(t, data)
		require.Equal(t, http.StatusOK, w.Code)

		resDto := &dto.OrderFindResponseDto{}
		_ = json.Unmarshal(data, resDto)

		require.Equal(t, len(oneOnly), len(resDto.Data.([]interface{})))
	})

	t.Run("ByAdminRole", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["session_id"] = sessUUID.String()
		claims["user_id"] = userUUID.String()
		claims["role"] = models.UserRoleAdmin
		claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
		validToken, _ := token.SignedString([]byte(cfg.Server.JwtSecretKey))

		req := httptest.NewRequest(http.MethodGet, "/order", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", validToken))
		w := httptest.NewRecorder()

		orderUC.EXPECT().FindAll(gomock.Any(), gomock.Any()).AnyTimes().Return(orders, nil)
		sessUC.EXPECT().GetSessionById(gomock.Any(), sessUUID.String()).AnyTimes().Return(&models.Session{UserID: userUUID, SessionID: sessUUID.String()}, nil)
		userUC.EXPECT().CachedFindById(gomock.Any(), userUUID).AnyTimes().Return(&models.User{UserID: userUUID}, nil)

		handler := http.HandlerFunc(handlers.FindAll)
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}
		require.NotNil(t, data)
		require.Equal(t, http.StatusOK, w.Code)

		resDto := &dto.OrderFindResponseDto{}
		_ = json.Unmarshal(data, resDto)

		require.Equal(t, len(orders), len(resDto.Data.([]interface{})))
	})

	t.Run("FindById", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["session_id"] = sessUUID.String()
		claims["user_id"] = userUUID.String()
		claims["role"] = models.UserRoleUser
		claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
		validToken, _ := token.SignedString([]byte(cfg.Server.JwtSecretKey))

		req := httptest.NewRequest(http.MethodGet, "/order?id="+m.OrderID.String(), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", validToken))
		w := httptest.NewRecorder()

		orderUC.EXPECT().CachedFindById(gomock.Any(), gomock.Any()).AnyTimes().Return(&m, nil)

		handler := http.HandlerFunc(handlers.FindById)
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}
		require.NotNil(t, data)
		require.Equal(t, http.StatusOK, w.Code)

		resDto := &dto.OrderResponseDto{}
		err = json.NewDecoder(res.Body).Decode(resDto)
		if err != nil {
			fmt.Println(err)
		}

		require.Equal(t, m.OrderID.String(), resDto.OrderID.String())
	})
}
