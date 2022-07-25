package handlers

func (h *productHandlersHTTP) ProductMapRoutes() {
	h.group.Use(h.mw.IsLoggedIn())
	h.group.GET("", h.FindAll())
	h.group.GET("/:id", h.FindById())

	h.group.POST("", h.Create(), h.mw.IsAdmin)
	h.group.PUT("/:id", h.UpdateById(), h.mw.IsAdmin)
	h.group.DELETE("/:id", h.DeleteById(), h.mw.IsAdmin)
}
