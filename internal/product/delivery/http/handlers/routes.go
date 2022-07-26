package handlers

import "net/http"

func (h *productHandlersHTTP) ProductMapRoutes() {
	h.mux.Handle("/product/create", h.mw.IsAdmin(http.HandlerFunc(h.Create)))
	h.mux.Handle("/product/brand", h.mw.GetHandler(http.HandlerFunc(h.FindAllByBrandId)))
}
