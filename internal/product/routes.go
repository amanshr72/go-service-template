package product

import (
	"github.com/go-chi/chi/v5"
)

func NewRouter(svc Service) chi.Router {
	h := NewHandler(svc)
	r := chi.NewRouter()

	r.Post("/", h.Create)
	r.Get("/", h.GetAll)
	r.Get("/{id}", h.GetByID)

	r.Get("/slow", h.SlowGetAll)

	return r
}
