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

	r.Use(middleware.Logger)
	r.Use(middleware.Gzip)

	// Публичные маршруты (без аутентификации)
	r.Group(func(r chi.Router) {
		r.Get("/ping", databaseHealthHandler.PingDatabase)
	})

	// Маршруты, требующие аутентификации
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)

		r.Post("/", shortenerHandler.CreateShortenerURL)
		r.Get("/{id}", shortenerHandler.GetShortenerURL)
		r.Post("/api/shorten", shortenerHandler.CreateJSONShortenerURL)
		r.Post("/api/shorten/batch", shortenerHandler.BatchCreateJSONShortenerURL)
		r.Get("/api/user/urls", shortenerHandler.getURLsByUser)
	})

	return r
}
