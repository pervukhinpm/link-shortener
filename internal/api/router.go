package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/pervukhinpm/link-shortener.git/internal/middleware"
)

func Router(
	databaseHealthHandler *DatabaseHealthHandler,
	shortenerHandler *ShortenerHandler,
) chi.Router {
	r := chi.NewRouter()
	r.Post("/", middleware.RequestLogger(middleware.GzipMiddleware(shortenerHandler.CreateShortenerURL)))
	r.Get("/{id}", middleware.RequestLogger(shortenerHandler.GetShortenerURL))
	r.Post("/api/shorten", middleware.RequestLogger(middleware.GzipMiddleware(shortenerHandler.CreateJSONShortenerURL)))
	r.Get("/ping", middleware.RequestLogger(databaseHealthHandler.PingDatabase))
	r.Post("/api/shorten/batch", middleware.RequestLogger(middleware.GzipMiddleware(shortenerHandler.BatchCreateJSONShortenerURL)))
	return r
}
