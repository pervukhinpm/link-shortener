package api

import (
	"fmt"
)

// ServerURL представляет URL сервера с разобранными компонентами.
// Используется для формирования полных URL в ответах API.
type ServerURL struct {
	Scheme string
	Host   string
	Port   int
}

// NewServerURL создает новый экземпляр ServerURL с указанными параметрами.
// Инициализирует все компоненты URL: схему, хост и порт.
func NewServerURL(scheme string, host string, port int) *ServerURL {
	return &ServerURL{
		Scheme: scheme,
		Host:   host,
		Port:   port,
	}
}

// String возвращает строковое представление URL сервера.
// Формирует полный URL в формате "scheme://host:port" или "host:port".
func (s *ServerURL) String() string {
	if s.Scheme != "" && s.Host != "" && s.Port != 0 {
		return fmt.Sprintf("%s://%s:%d", s.Scheme, s.Host, s.Port)
	}
	if s.Scheme == "" && s.Host != "" && s.Port != 0 {
		return fmt.Sprintf("%s:%d", s.Host, s.Port)
	}
	return ""
}
