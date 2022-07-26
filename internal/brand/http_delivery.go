package brand

import (
	"net/http"
)

// Brand HTTP Handlers interface
type BrandHandlers interface {
	Create(w http.ResponseWriter, r *http.Request)
	FindAll(w http.ResponseWriter, r *http.Request)
	FindById(w http.ResponseWriter, r *http.Request)
}
