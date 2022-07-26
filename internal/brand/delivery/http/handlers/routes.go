package handlers

import "net/http"

func (h *brandHandlersHTTP) BrandMapRoutes() {
	h.mux.Handle("/brand/create", h.mw.IsAdmin(http.HandlerFunc(h.Create)))
	h.mux.Handle("/brand", h.mw.GetHandler(http.HandlerFunc(h.FindAll)))
	h.mux.Handle("/brand?id=", h.mw.GetHandler(http.HandlerFunc(h.FindById)))
}
