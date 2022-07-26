package product

import (
	"net/http"
)

// Product HTTP Handlers interface
type ProductHandlers interface {
	Create(w http.ResponseWriter, r *http.Request)
	FindAll(w http.ResponseWriter, r *http.Request)
	FindAllByBrandId(w http.ResponseWriter, r *http.Request)
}
