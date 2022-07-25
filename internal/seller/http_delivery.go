package seller

import "github.com/labstack/echo/v4"

// Seller HTTP Handlers interface
type SellerHandlers interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	GetMe() echo.HandlerFunc
	FindAll() echo.HandlerFunc
	FindById() echo.HandlerFunc
	UpdateById() echo.HandlerFunc
	DeleteById() echo.HandlerFunc
	Logout() echo.HandlerFunc
	RefreshToken() echo.HandlerFunc
}
