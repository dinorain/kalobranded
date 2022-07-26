package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"

	"github.com/dinorain/kalobranded/config"
	"github.com/dinorain/kalobranded/internal/models"
	httpErrors "github.com/dinorain/kalobranded/pkg/http_errors"
	"github.com/dinorain/kalobranded/pkg/logger"
)

type MiddlewareManager interface {
	RequestLoggerMiddleware(next http.Handler) http.Handler
	PostHandler(next http.Handler) http.Handler
	GetHandler(next http.Handler) http.Handler
	IsLoggedIn(next http.Handler) http.Handler
	IsUser(next http.Handler) http.Handler
	IsAdmin(next http.Handler) http.Handler
	GetJWTClaims(w http.ResponseWriter, r *http.Request) (*jwt.MapClaims, error)
}

type middlewareManager struct {
	logger logger.Logger
	cfg    *config.Config
}

var _ MiddlewareManager = (*middlewareManager)(nil)

func NewMiddlewareManager(logger logger.Logger, cfg *config.Config) *middlewareManager {
	return &middlewareManager{logger: logger, cfg: cfg}
}

func (mw *middlewareManager) PostHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httpErrors.NewStatusMethodNotAllowedError(w, nil, mw.cfg.Http.DebugErrorsResponse)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (mw *middlewareManager) GetHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			httpErrors.NewStatusMethodNotAllowedError(w, nil, mw.cfg.Http.DebugErrorsResponse)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (mw *middlewareManager) GetJWTClaims(w http.ResponseWriter, r *http.Request) (*jwt.MapClaims, error) {
	authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	if len(authHeader) != 2 {
		return nil, httpErrors.NewUnauthorizedError(w, nil, mw.cfg.Http.DebugErrorsResponse)
	} else {
		jwtToken := authHeader[1]
		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(mw.cfg.Server.JwtSecretKey), nil
		})

		if err != nil {
			return nil, httpErrors.NewUnauthorizedError(w, nil, mw.cfg.Http.DebugErrorsResponse)
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return &claims, nil
		} else {
			return nil, httpErrors.NewUnauthorizedError(w, nil, mw.cfg.Http.DebugErrorsResponse)
		}
	}

	return nil, httpErrors.NewUnauthorizedError(w, nil, mw.cfg.Http.DebugErrorsResponse)
}

func (mw *middlewareManager) IsLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := mw.GetJWTClaims(w, r); err == nil {
			next.ServeHTTP(w, r)
		}
		return
	})
}

func (mw *middlewareManager) IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtClaims, err := mw.GetJWTClaims(w, r)
		if err != nil {
			return
		}
		claims := *jwtClaims
		role, ok := claims["role"].(string)
		if !ok {
			mw.logger.Warnf("role: %+v", claims)
		}

		if role != models.UserRoleAdmin {
			httpErrors.NewForbiddenError(w, nil, mw.cfg.Http.DebugErrorsResponse)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (mw *middlewareManager) IsUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtClaims, err := mw.GetJWTClaims(w, r)
		if err != nil {
			return
		}
		claims := *jwtClaims
		role, ok := claims["role"].(string)
		if !ok {
			mw.logger.Warnf("role: %+v", claims)
		}

		if role != models.UserRoleUser {
			_ = httpErrors.NewForbiddenError(w, nil, mw.cfg.Http.DebugErrorsResponse)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (mw *middlewareManager) RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !mw.checkIgnoredURI(r.RequestURI, mw.cfg.Http.IgnoreLogUrls) {
			mw.logger.HttpMiddlewareAccessLogger(r.Method, r.URL.String())
		}

		next.ServeHTTP(w, r)
	})
}

func (mw *middlewareManager) checkIgnoredURI(requestURI string, uriList []string) bool {
	for _, s := range uriList {
		if strings.Contains(requestURI, s) {
			return true
		}
	}
	return false
}
