package httpx

import (
	"example.com/prc_notes_api/internal/http/handlers"
	"github.com/go-chi/chi/v5"
)

func NewRouter(h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/api/v1/notes", h.CreateNote)
	r.Get("/api/v1/notes/{id}", h.GetNote)
	r.Get("/api/v1/notes", h.ListNotes)
	r.Put("/api/v1/notes/{id}", h.Update)
	r.Patch("/api/v1/notes/{id}", h.Update)
	r.Delete("/api/v1/notes/{id}", h.Delete)
	return r
}
