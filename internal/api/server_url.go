package api

import (
	"fmt"
)

type ServerURL struct {
	Scheme string
	Host   string
	Port   int
}

func NewServerURL(scheme string, host string, port int) *ServerURL {
	return &ServerURL{
		Scheme: scheme,
		Host:   host,
		Port:   port,
	}
}

func (s *ServerURL) String() string {
	if s.Scheme != "" && s.Host != "" && s.Port != 0 {
		return fmt.Sprintf("%s://%s:%d", s.Scheme, s.Host, s.Port)
	}
	if s.Scheme == "" && s.Host != "" && s.Port != 0 {
		return fmt.Sprintf("%s:%d", s.Host, s.Port)
	}
	return ""
}
