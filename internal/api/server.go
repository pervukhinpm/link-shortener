package api

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Server struct {
	router    chi.Router
	serverURL *ServerURL
}

func NewServer(serverURL *ServerURL, router chi.Router) *Server {
	return &Server{
		router:    router,
		serverURL: serverURL,
	}
}

func (s *Server) Start() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", s.serverURL.Port), s.router)
}
