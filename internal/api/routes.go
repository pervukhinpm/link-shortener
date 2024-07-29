package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/pervukhinpm/link-shortener.git/internal/middleware"
)

func (h *ShortenerHandler) useRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.HandleFunc("POST /", middleware.GzipMiddleware(middleware.RequestLogger(h.CreateShortenerURL)))
	r.HandleFunc("GET /{id}", middleware.GzipMiddleware(middleware.RequestLogger(h.GetShortenerURL)))
	r.HandleFunc("POST /api/shorten", middleware.GzipMiddleware(middleware.RequestLogger(h.CreateJSONShortenerURL)))
	return r
}
