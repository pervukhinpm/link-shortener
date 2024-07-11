package api

import "github.com/go-chi/chi/v5"

func (h *ShortenerHandler) useRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.HandleFunc("POST /", h.CreateShortenerURL)
	r.HandleFunc("GET /{id}", h.GetShortenerURL)

	return r
}
