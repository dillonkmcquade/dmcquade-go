package app

func (h *App) loadRoutes() {
	h.Router.Get("/", index(h))
	h.Router.Get("/project/{id}", projectsPage(h))
	h.Router.NotFound(h.NotFound)
}
