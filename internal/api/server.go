package api

import (
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
	return http.ListenAndServe(s.serverURL.String(), s.shortenerHandler.useRoutes())
}
