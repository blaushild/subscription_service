package server

import (
	"github.com/go-chi/chi/v5"
)

// setupRouter настраивает роутинг сервера
func setupRouter(r chi.Router, c Controller) {
	r.Group(func(r chi.Router) {
		r.Post("/api/v1/subscription", c.Create)
		r.Get("/api/v1/subscription/{id}", c.GetRecordByID)
		r.Put("/api/v1/subscription", c.Update)
		r.Delete("/api/v1/subscription/{id}", c.Delete)
		r.Get("/api/v1/subscription", c.List)
		r.Get("/api/v1/subscription/total", c.GetTotal)
	})

	r.Get("/swagger/*", c.Swagger)
}
