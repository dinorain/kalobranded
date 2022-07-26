package handlers

import "net/http"

func (h *orderHandlersHTTP) OrderMapRoutes() {
	h.mux.Handle("/order/create", h.mw.IsAdmin(http.HandlerFunc(h.Create)))
	h.mux.Handle("/order", h.mw.GetHandler(http.HandlerFunc(h.FindAll)))
}
