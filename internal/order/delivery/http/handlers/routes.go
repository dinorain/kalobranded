package handlers

func (h *orderHandlersHTTP) OrderMapRoutes() {
	h.group.Use(h.mw.IsLoggedIn())
	h.group.GET("", h.FindAll())
	h.group.POST("", h.Create(), h.mw.IsUser)

	h.group.GET("/:id", h.FindById())
	h.group.POST("/:id", h.AcceptById(), h.mw.IsSeller)
}
