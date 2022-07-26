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

	"github.com/go-playground/validator"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
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

	cfg := &config.Config{Session: config.Session{Expire: 1234}, Server: config.ServerConfig{JwtSecretKey: "secret"}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, cfg)

	v := validator.New()

	mux := http.NewServeMux()
	handlers := NewProductHandlersHTTP(mux, appLogger, cfg, mw, v, productUC, sessUC)

	reqDto := &dto.ProductCreateRequestDto{
		Name:        "Name",
		Description: "Description",
		Price:       10000.0,
	}

	buf := &bytes.Buffer{}
	_ = json.NewEncoder(buf).Encode(reqDto)

	userUUID := uuid.New()
	brandUUID := uuid.New()
	sessUUID := uuid.New()
	productUUID := uuid.New()

	req := httptest.NewRequest(http.MethodPost, "/product/create", buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	wDto := &dto.ProductCreateResponseDto{
		ProductID: productUUID,
	}

	buf, _ = converter.AnyToBytesBuffer(wDto)

	sessUC.EXPECT().GetSessionById(gomock.Any(), sessUUID.String()).AnyTimes().Return(&models.Session{UserID: userUUID, SessionID: sessUUID.String()}, nil)
	productUC.EXPECT().Create(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Product{ProductID: productUUID, BrandID: brandUUID}, nil)

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

func TestProductsHandler_Find(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	productUC := mock.NewMockProductUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}, Server: config.ServerConfig{JwtSecretKey: "secret"}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, cfg)

	v := validator.New()

	mux := http.NewServeMux()
	handlers := NewProductHandlersHTTP(mux, appLogger, cfg, mw, v, productUC, sessUC)

	brandUUID := uuid.New()

	var products []models.Product
	var oneOnly []models.Product

	m := models.Product{
		ProductID:   uuid.New(),
		Name:        "Name",
		Description: "FirstName",
		Price:       10000.00,
		BrandID:     brandUUID,
	}
	oneOnly = append(oneOnly, m)
	products = append(products, m, models.Product{
		ProductID:   uuid.New(),
		Name:        "Name",
		Description: "FirstName",
		Price:       10000.00,
		BrandID:     uuid.New(),
	})

	t.Run("FindAll", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/product", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		productUC.EXPECT().FindAll(gomock.Any(), gomock.Any()).AnyTimes().Return(products, nil)

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

		resDto := &dto.ProductFindResponseDto{}
		_ = json.Unmarshal(data, resDto)

		require.Equal(t, len(products), len(resDto.Data.([]interface{})))
	})

	t.Run("FindAllByBrandId", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/product/brand?id="+brandUUID.String(), nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		productUC.EXPECT().FindAllByBrandId(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(oneOnly, nil)

		handler := http.HandlerFunc(handlers.FindAllByBrandId)
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}
		require.NotNil(t, data)
		require.Equal(t, http.StatusOK, w.Code)

		resDto := &dto.ProductFindResponseDto{}
		_ = json.Unmarshal(data, resDto)

		require.Equal(t, len(oneOnly), len(resDto.Data.([]interface{})))
	})

	t.Run("FindById", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/product?id="+m.ProductID.String(), nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		productUC.EXPECT().CachedFindById(gomock.Any(), gomock.Any()).AnyTimes().Return(&m, nil)

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

		resDto := &dto.ProductResponseDto{}
		err = json.NewDecoder(res.Body).Decode(resDto)
		if err != nil {
			fmt.Println(err)
		}

		require.Equal(t, m.ProductID.String(), resDto.ProductID.String())
	})
}
