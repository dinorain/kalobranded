package brand

import "github.com/labstack/echo/v4"

// Brand HTTP Handlers interface
type BrandHandlers interface {
	Register() echo.HandlerFunc
	FindAll() echo.HandlerFunc
	FindById() echo.HandlerFunc
	UpdateById() echo.HandlerFunc
	DeleteById() echo.HandlerFunc
}
