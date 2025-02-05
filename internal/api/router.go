package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/pervukhinpm/link-shortener.git/internal/middleware"
	"net/http/pprof"
)

func Router(
	databaseHealthHandler *DatabaseHealthHandler,
	shortenerHandler *ShortenerHandler,
) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Gzip)

	// Подключаем обработчики pprof
	// Подключаем pprof-эндпоинты
	r.Route("/debug/pprof", func(r chi.Router) {
		r.HandleFunc("/", pprof.Index)
		r.HandleFunc("/cmdline", pprof.Cmdline)
		r.HandleFunc("/profile", pprof.Profile)
		r.HandleFunc("/symbol", pprof.Symbol)
		r.HandleFunc("/trace", pprof.Trace)

		// Добавляем обработку всех стандартных pprof-метрик
		r.HandleFunc("/allocs", pprof.Handler("allocs").ServeHTTP)
		r.HandleFunc("/block", pprof.Handler("block").ServeHTTP)
		r.HandleFunc("/goroutine", pprof.Handler("goroutine").ServeHTTP)
		r.HandleFunc("/heap", pprof.Handler("heap").ServeHTTP)
		r.HandleFunc("/mutex", pprof.Handler("mutex").ServeHTTP)
		r.HandleFunc("/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
	})

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
		r.Delete("/api/user/urls", shortenerHandler.DeleteURLBatchByUser)
	})

	return r
}
