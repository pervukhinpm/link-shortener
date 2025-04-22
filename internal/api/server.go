package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Server представляет HTTP-сервер приложения.
// Обрабатывает входящие запросы с помощью маршрутизатора.
type Server struct {
	router    chi.Router
	serverURL *ServerURL
}

// NewServer создает новый экземпляр HTTP-сервера.
// Инициализирует маршрутизатор и базовый URL.
func NewServer(serverURL *ServerURL, router chi.Router) *Server {
	return &Server{
		router:    router,
		serverURL: serverURL,
	}
}

// Start запускает HTTP-сервер на указанном порту.
// Использует маршрутизатор для обработки запросов.
func (s *Server) Start() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", s.serverURL.Port), s.router)
}
