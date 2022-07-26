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
	"github.com/dinorain/kalobranded/internal/middlewares"
	"github.com/dinorain/kalobranded/internal/models"
	mockSessUC "github.com/dinorain/kalobranded/internal/session/mock"
	"github.com/dinorain/kalobranded/internal/user/delivery/http/dto"
	"github.com/dinorain/kalobranded/internal/user/mock"
	"github.com/dinorain/kalobranded/pkg/converter"
	"github.com/dinorain/kalobranded/pkg/logger"
)

func TestUsersHandler_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mock.NewMockUserUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}, Server: config.ServerConfig{JwtSecretKey: "secret"}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, cfg)

	v := validator.New()

	mux := http.NewServeMux()
	handlers := NewUserHandlersHTTP(mux, appLogger, cfg, mw, v, userUC, sessUC)

	reqDto := &dto.UserRegisterRequestDto{
		Email:           "email@gmail.com",
		FirstName:       "FirstName",
		LastName:        "LastName",
		Password:        "123456",
		Role:            models.UserRoleUser,
		DeliveryAddress: "DeliveryAddress",
	}

	buf := &bytes.Buffer{}
	_ = json.NewEncoder(buf).Encode(reqDto)

	req := httptest.NewRequest(http.MethodPost, "/user", buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	resDto := &dto.UserRegisterResponseDto{
		UserID: uuid.Nil,
	}

	buf, _ = converter.AnyToBytesBuffer(resDto)

	userUC.EXPECT().Register(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.User{}, nil)

	handler := http.HandlerFunc(handlers.Register)
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

func TestUsersHandler_Login(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mock.NewMockUserUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}, Server: config.ServerConfig{JwtSecretKey: "secret"}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, cfg)

	v := validator.New()

	mux := http.NewServeMux()
	handlers := NewUserHandlersHTTP(mux, appLogger, cfg, mw, v, userUC, sessUC)

	reqDto := &dto.UserLoginRequestDto{
		Email:    "email@gmail.com",
		Password: "123456",
	}

	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(reqDto)

	req := httptest.NewRequest(http.MethodPost, "/user/login", &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockUser := &models.User{
		UserID:          uuid.New(),
		Email:           "email@gmail.com",
		FirstName:       "FirstName",
		LastName:        "LastName",
		Password:        "123456",
		Role:            models.UserRoleUser,
		DeliveryAddress: "DeliveryAddress",
	}

	userUC.EXPECT().Login(gomock.Any(), reqDto.Email, reqDto.Password).AnyTimes().Return(mockUser, nil)
	sessUC.EXPECT().CreateSession(gomock.Any(), &models.Session{UserID: mockUser.UserID}, cfg.Session.Expire).AnyTimes().Return("s", nil)
	userUC.EXPECT().GenerateTokenPair(gomock.Any(), gomock.Any()).AnyTimes().Return("rt", "at", nil)

	handler := http.HandlerFunc(handlers.Login)
	handler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	require.NotNil(t, data)
	require.Equal(t, http.StatusOK, w.Code)

	resDto := &dto.UserLoginResponseDto{}
	err = json.NewDecoder(res.Body).Decode(resDto)
	if err != nil {
		fmt.Println(err)
	}

	require.Equal(t, mockUser.UserID.String(), resDto.UserID.String())
}

func TestUsersHandler_FindAll(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mock.NewMockUserUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}, Server: config.ServerConfig{JwtSecretKey: "secret"}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, cfg)

	v := validator.New()

	mux := http.NewServeMux()
	handlers := NewUserHandlersHTTP(mux, appLogger, cfg, mw, v, userUC, sessUC)

	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	var users []models.User
	users = append(users, models.User{
		UserID:          uuid.New(),
		Email:           "email@gmail.com",
		FirstName:       "FirstName",
		LastName:        "LastName",
		Password:        "123456",
		Role:            models.UserRoleUser,
		DeliveryAddress: "DeliveryAddress",
	})

	userUC.EXPECT().FindAll(gomock.Any(), gomock.Any()).AnyTimes().Return(users, nil)

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

	resDto := &dto.UserFindResponseDto{}
	_ = json.Unmarshal(data, resDto)

	require.Equal(t, len(users), len(resDto.Data.([]interface{})))
}

func TestUsersHandler_FindById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mock.NewMockUserUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}, Server: config.ServerConfig{JwtSecretKey: "secret"}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, cfg)

	v := validator.New()

	mux := http.NewServeMux()
	handlers := NewUserHandlersHTTP(mux, appLogger, cfg, mw, v, userUC, sessUC)

	userUUID := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/user?id="+userUUID.String(), nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	userUC.EXPECT().CachedFindById(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.User{UserID: userUUID}, nil)

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

	resDto := &dto.UserResponseDto{}
	err = json.NewDecoder(res.Body).Decode(resDto)
	if err != nil {
		fmt.Println(err)
	}

	require.Equal(t, userUUID.String(), resDto.UserID.String())
}

func TestUsersHandler_GetMe(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mock.NewMockUserUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}, Server: config.ServerConfig{JwtSecretKey: "secret"}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, cfg)

	v := validator.New()

	mux := http.NewServeMux()
	handlers := NewUserHandlersHTTP(mux, appLogger, cfg, mw, v, userUC, sessUC)

	userUUID := uuid.New()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["session_id"] = uuid.New().String()
	claims["user_id"] = userUUID.String()
	claims["role"] = models.UserRoleUser
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	validToken, _ := token.SignedString([]byte("secret"))

	req := httptest.NewRequest(http.MethodGet, "/user/me", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", validToken))

	w := httptest.NewRecorder()

	sessUC.EXPECT().GetSessionById(gomock.Any(), claims["session_id"].(string)).AnyTimes().Return(&models.Session{}, nil)
	userUC.EXPECT().CachedFindById(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.User{UserID: userUUID}, nil)

	handler := http.HandlerFunc(handlers.GetMe)
	handler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	require.NotNil(t, data)
	require.Equal(t, http.StatusOK, w.Code)

	resDto := &dto.UserResponseDto{}
	err = json.NewDecoder(res.Body).Decode(resDto)
	if err != nil {
		fmt.Println(err)
	}

	require.Equal(t, userUUID.String(), resDto.UserID.String())
}

func TestUsersHandler_Logout(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mock.NewMockUserUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}, Server: config.ServerConfig{JwtSecretKey: "secret"}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, cfg)

	v := validator.New()

	mux := http.NewServeMux()
	handlers := NewUserHandlersHTTP(mux, appLogger, cfg, mw, v, userUC, sessUC)

	userUUID := uuid.New()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["session_id"] = uuid.New().String()
	claims["user_id"] = userUUID.String()
	claims["role"] = models.UserRoleUser
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	validToken, _ := token.SignedString([]byte("secret"))

	req := httptest.NewRequest(http.MethodGet, "/user/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", validToken))

	w := httptest.NewRecorder()

	sessUC.EXPECT().DeleteById(gomock.Any(), claims["session_id"].(string)).AnyTimes().Return(nil)

	handler := http.HandlerFunc(handlers.Logout)
	handler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	require.NotNil(t, data)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestUsersHandler_RefreshToken(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUC := mock.NewMockUserUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}, Server: config.ServerConfig{JwtSecretKey: "secret"}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, cfg)

	v := validator.New()

	mux := http.NewServeMux()
	handlers := NewUserHandlersHTTP(mux, appLogger, cfg, mw, v, userUC, sessUC)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["session_id"] = uuid.New().String()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	validToken, _ := token.SignedString([]byte("secret"))

	reqDto := &dto.UserRefreshTokenDto{
		RefreshToken: validToken,
	}

	buf := &bytes.Buffer{}
	_ = json.NewEncoder(buf).Encode(reqDto)

	req := httptest.NewRequest(http.MethodPost, "/user/refresh", buf)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	sessUC.EXPECT().GetSessionById(gomock.Any(), claims["session_id"].(string)).AnyTimes().Return(&models.Session{}, nil)
	userUC.EXPECT().FindById(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.User{}, nil)
	userUC.EXPECT().GenerateTokenPair(gomock.Any(), gomock.Any()).AnyTimes().Return("rt", "at", nil)

	handler := http.HandlerFunc(handlers.RefreshToken)
	handler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	require.NotNil(t, data)
	require.Equal(t, http.StatusOK, w.Code)
}
