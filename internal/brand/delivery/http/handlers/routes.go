package handlers

func (h *brandHandlersHTTP) BrandMapRoutes() {
	h.group.Use(h.mw.IsLoggedIn())

	h.group.GET("/:id", h.FindById())

	h.group.PUT("/:id", h.UpdateById(), h.mw.IsAdmin)
	h.group.POST("", h.Register(), h.mw.IsAdmin)
	h.group.GET("", h.FindAll(), h.mw.IsAdmin)
	h.group.DELETE("/:id", h.DeleteById(), h.mw.IsAdmin)
}
