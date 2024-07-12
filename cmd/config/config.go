package config

import (
	"flag"
	"github.com/pervukhinpm/link-shortener.git/internal/api"
	"net/url"
	"strconv"
)

var ServerConfig struct {
	ServerAddress api.ServerURL
	BaseUrl       api.ServerURL
}

func ParseFlags() {
	var flagServerAddress string
	var flagBaseURL string

	flag.StringVar(&flagServerAddress, "a", "localhost:8080", "Host Port")
	flag.StringVar(&flagBaseURL, "b", "http://localhost:8080/", "Base URL")
	flag.Parse()

	ServerConfig.ServerAddress = *parseServerURL(flagServerAddress)
	ServerConfig.BaseUrl = *parseServerURL(flagBaseURL)
}

func parseServerURL(rawURL string) *api.ServerURL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}

	host := u.Hostname()
	if host == "" {
		host = "localhost"
	}

	port := u.Port()
	if port == "" {
		port = "8080"
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		portInt = 8080
	}

	scheme := u.Scheme
	if scheme == "" {
		scheme = ""
	}

	return api.NewServerURL(scheme, host, portInt)
}
