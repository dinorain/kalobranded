package order

import (
	"net/http"
)

// Order HTTP Handlers interface
type OrderHandlers interface {
	Create(w http.ResponseWriter, r *http.Request)
	FindAll(w http.ResponseWriter, r *http.Request)
}
