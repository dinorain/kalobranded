package httpUtils

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/dinorain/checkoutaja/internal/models"
	httpErrors "github.com/dinorain/checkoutaja/pkg/http_errors"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

// Get request id from echo context
func GetRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

// ReqIDCtxKey is a key used for the Request ID in context
type ReqIDCtxKey struct{}

// Get ctx with timeout and request id from echo context
func GetCtxWithReqID(c echo.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*15)
	ctx = context.WithValue(ctx, ReqIDCtxKey{}, GetRequestID(c))
	return ctx, cancel
}

// Get context  with request id
func GetRequestCtx(c echo.Context) context.Context {
	return context.WithValue(c.Request().Context(), ReqIDCtxKey{}, GetRequestID(c))
}

// UserCtxKey is a key used for the User object in the context
type UserCtxKey struct{}

// Get user from context
func GetUserFromCtx(ctx context.Context) (*models.User, error) {
	user, ok := ctx.Value(UserCtxKey{}).(*models.User)
	if !ok {
		return nil, httpErrors.Unauthorized
	}

	return user, nil
}

// NewWriterWrapper response writer wrapper
func NewWriterWrapper(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{ResponseWriter: w}
}

func (rw *responseWriterWrapper) Status() int {
	return rw.status
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}
