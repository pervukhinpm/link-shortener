package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/pervukhinpm/link-shortener.git/internal/middleware"
)

func (h *ShortenerHandler) useRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.HandleFunc("POST /", middleware.RequestLogger(h.CreateShortenerURL))
	r.HandleFunc("GET /{id}", h.GetShortenerURL)

	return r
}
