package api

import (
	"fmt"
	"net/url"
	"strconv"
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

func (s *ServerURL) Set(value string) error {
	parsedURL, err := url.Parse(value)
	if err != nil {
		return err
	}

	s.Scheme = parsedURL.Scheme
	s.Host = parsedURL.Hostname()
	portStr := parsedURL.Port()
	if portStr != "" {
		s.Port, err = strconv.Atoi(portStr)
		if err != nil {
			return err
		}
	} else {
		s.Port = 0
	}

	return nil
}
