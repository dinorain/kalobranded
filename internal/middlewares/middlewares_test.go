package middlewares

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/dinorain/kalobranded/config"
	"github.com/dinorain/kalobranded/internal/models"
	"github.com/dinorain/kalobranded/pkg/logger"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestMiddlewares_IsLoggedIn(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{Session: config.Session{Expire: 1234}, Server: config.ServerConfig{JwtSecretKey: "secret"}}
	appLogger := logger.NewAppLogger(cfg)
	mw := NewMiddlewareManager(appLogger, cfg)

	t.Run("Fail", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler := mw.IsLoggedIn(http.HandlerFunc(testHandler))
		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Success", func(t *testing.T) {
		userUUID := uuid.New()
		sessUUID := uuid.New()
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["session_id"] = sessUUID.String()
		claims["user_id"] = userUUID.String()
		claims["role"] = models.UserRoleUser
		claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
		validToken, _ := token.SignedString([]byte(cfg.Server.JwtSecretKey))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", validToken))
		w := httptest.NewRecorder()

		handler := mw.IsLoggedIn(http.HandlerFunc(testHandler))
		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})
}

func TestMiddlewares_IsAdmin(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{Session: config.Session{Expire: 1234}, Server: config.ServerConfig{JwtSecretKey: "secret"}}
	appLogger := logger.NewAppLogger(cfg)
	mw := NewMiddlewareManager(appLogger, cfg)

	userUUID := uuid.New()
	sessUUID := uuid.New()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["session_id"] = sessUUID.String()
	claims["user_id"] = userUUID.String()
	claims["role"] = models.UserRoleAdmin
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	validToken, _ := token.SignedString([]byte(cfg.Server.JwtSecretKey))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", validToken))
	w := httptest.NewRecorder()

	handler := mw.IsAdmin(http.HandlerFunc(testHandler))
	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestMiddlewares_IsUser(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{Session: config.Session{Expire: 1234}, Server: config.ServerConfig{JwtSecretKey: "secret"}}
	appLogger := logger.NewAppLogger(cfg)
	mw := NewMiddlewareManager(appLogger, cfg)

	userUUID := uuid.New()
	sessUUID := uuid.New()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["session_id"] = sessUUID.String()
	claims["user_id"] = userUUID.String()
	claims["role"] = models.UserRoleUser
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	validToken, _ := token.SignedString([]byte(cfg.Server.JwtSecretKey))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", validToken))
	w := httptest.NewRecorder()

	handler := mw.IsUser(http.HandlerFunc(testHandler))
	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestMiddlewares_PostHandler(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{Session: config.Session{Expire: 1234}, Server: config.ServerConfig{JwtSecretKey: "secret"}}
	appLogger := logger.NewAppLogger(cfg)
	mw := NewMiddlewareManager(appLogger, cfg)

	t.Run("Fail", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler := mw.PostHandler(http.HandlerFunc(testHandler))
		handler.ServeHTTP(w, req)
		require.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler := mw.PostHandler(http.HandlerFunc(testHandler))
		handler.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	})
}

func TestMiddlewares_GetHandler(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{Session: config.Session{Expire: 1234}, Server: config.ServerConfig{JwtSecretKey: "secret"}}
	appLogger := logger.NewAppLogger(cfg)
	mw := NewMiddlewareManager(appLogger, cfg)

	t.Run("Fail", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler := mw.GetHandler(http.HandlerFunc(testHandler))
		handler.ServeHTTP(w, req)
		require.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler := mw.GetHandler(http.HandlerFunc(testHandler))
		handler.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	})

}
