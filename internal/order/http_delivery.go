package order

import "github.com/labstack/echo/v4"

// Order HTTP Handlers interface
type OrderHandlers interface {
	Create() echo.HandlerFunc
	FindAll() echo.HandlerFunc
	FindById() echo.HandlerFunc
}
