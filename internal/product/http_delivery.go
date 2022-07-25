package product

import "github.com/labstack/echo/v4"

// Product HTTP Handlers interface
type ProductHandlers interface {
	Create() echo.HandlerFunc
	FindAll() echo.HandlerFunc
	FindById() echo.HandlerFunc
	UpdateById() echo.HandlerFunc
	DeleteById() echo.HandlerFunc
}
