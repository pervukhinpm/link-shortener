package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/pervukhinpm/link-shortener.git/internal/middleware"
)

func (h *ShortenerHandler) useRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.HandleFunc("POST /", middleware.RequestLogger(middleware.GzipMiddleware(h.CreateShortenerURL)))
	r.HandleFunc("GET /{id}", middleware.RequestLogger(middleware.GzipMiddleware(h.GetShortenerURL)))
	r.HandleFunc("POST /api/shorten", middleware.RequestLogger(middleware.GzipMiddleware(h.CreateJSONShortenerURL)))
	return r
}
