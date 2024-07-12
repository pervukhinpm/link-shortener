package api

import (
	"fmt"
	"net/http"
)

type Server struct {
	shortenerHandler *ShortenerHandler
	serverURL        *ServerURL
}

func NewServer(serverURL *ServerURL, handler *ShortenerHandler) *Server {
	return &Server{
		shortenerHandler: handler,
		serverURL:        serverURL,
	}
}

func (s *Server) Start() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", s.serverURL.Port), s.shortenerHandler.useRoutes())
}
