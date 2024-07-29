package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/pervukhinpm/link-shortener.git/internal/middleware"
)

func (h *ShortenerHandler) useRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", middleware.RequestLogger(middleware.GzipMiddleware(h.CreateShortenerURL)))
	r.Get("/{id}", middleware.RequestLogger(h.GetShortenerURL))
	r.Post("/api/shorten", middleware.RequestLogger(middleware.GzipMiddleware(h.CreateJSONShortenerURL)))
	return r
}
