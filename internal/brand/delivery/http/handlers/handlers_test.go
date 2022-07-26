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
	"github.com/dinorain/kalobranded/internal/brand/delivery/http/dto"
	"github.com/dinorain/kalobranded/internal/brand/mock"
	"github.com/dinorain/kalobranded/internal/middlewares"
	"github.com/dinorain/kalobranded/internal/models"
	mockSessUC "github.com/dinorain/kalobranded/internal/session/mock"
	"github.com/dinorain/kalobranded/pkg/converter"
	"github.com/dinorain/kalobranded/pkg/logger"
)

func TestBrandsHandler_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	brandUC := mock.NewMockBrandUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}, Server: config.ServerConfig{JwtSecretKey: "secret"}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, cfg)

	v := validator.New()

	mux := http.NewServeMux()
	handlers := NewBrandHandlersHTTP(mux, appLogger, cfg, mw, v, brandUC, sessUC)

	reqDto := &dto.BrandRegisterRequestDto{
		BrandName:     "BrandName",
		PickupAddress: "PickupAddress",
	}

	buf := &bytes.Buffer{}
	_ = json.NewEncoder(buf).Encode(reqDto)

	userUUID := uuid.New()
	brandUUID := uuid.New()
	sessUUID := uuid.New()

	req := httptest.NewRequest(http.MethodPost, "/brand/create", buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	wDto := &dto.BrandRegisterResponseDto{
		BrandID: brandUUID,
	}

	buf, _ = converter.AnyToBytesBuffer(wDto)

	sessUC.EXPECT().GetSessionById(gomock.Any(), sessUUID.String()).AnyTimes().Return(&models.Session{UserID: userUUID, SessionID: sessUUID.String()}, nil)
	brandUC.EXPECT().Register(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Brand{BrandID: brandUUID}, nil)

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

func TestBrandsHandler_Find(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	brandUC := mock.NewMockBrandUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}, Server: config.ServerConfig{JwtSecretKey: "secret"}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, cfg)

	v := validator.New()

	mux := http.NewServeMux()
	handlers := NewBrandHandlersHTTP(mux, appLogger, cfg, mw, v, brandUC, sessUC)

	brandUUID := uuid.New()

	var brands []models.Brand
	var oneOnly []models.Brand

	m := models.Brand{
		BrandID:       brandUUID,
		BrandName:     "BrandName",
		PickupAddress: "PickupAddress",
	}
	oneOnly = append(oneOnly, m)
	brands = append(brands, m, models.Brand{
		BrandID:       uuid.New(),
		BrandName:     "BrandName",
		PickupAddress: "PickupAddress",
	})

	t.Run("FindAll", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/brand", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		brandUC.EXPECT().FindAll(gomock.Any(), gomock.Any()).AnyTimes().Return(brands, nil)

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

		resDto := &dto.BrandFindResponseDto{}
		_ = json.Unmarshal(data, resDto)

		require.Equal(t, len(brands), len(resDto.Data.([]interface{})))
	})

	t.Run("FindById", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/brand?id="+m.BrandID.String(), nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		brandUC.EXPECT().CachedFindById(gomock.Any(), gomock.Any()).AnyTimes().Return(&m, nil)

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

		resDto := &dto.BrandResponseDto{}
		err = json.NewDecoder(res.Body).Decode(resDto)
		if err != nil {
			fmt.Println(err)
		}

		require.Equal(t, m.BrandID.String(), resDto.BrandID.String())
	})
}
