package handlers

import "net/http"

func (h *userHandlersHTTP) UserMapRoutes() {
	h.mux.Handle("/user/create", http.HandlerFunc(h.Register))
	h.mux.Handle("/user", h.mw.GetHandler(http.HandlerFunc(h.FindAll)))
	h.mux.Handle("/user?id=", h.mw.GetHandler(http.HandlerFunc(h.FindById)))
	h.mux.Handle("/user/me", h.mw.GetHandler(http.HandlerFunc(h.GetMe)))
	h.mux.Handle("/user/login", h.mw.PostHandler(http.HandlerFunc(h.Login)))
	h.mux.Handle("/user/logout", h.mw.PostHandler(http.HandlerFunc(h.Logout)))
	h.mux.Handle("/user/refresh", h.mw.PostHandler(http.HandlerFunc(h.RefreshToken)))
}
