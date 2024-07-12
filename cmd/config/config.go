package config

import (
	"flag"
	"github.com/pervukhinpm/link-shortener.git/internal/api"
)

var ServerConfig struct {
	ServerAddress api.ServerURL
	BaseUrl       api.ServerURL
}

func ParseFlags() {
	flag.Var(
		&ServerConfig.ServerAddress,
		"a",
		"Host Port",
	)
	flag.Var(
		&ServerConfig.BaseUrl,
		"b",
		"Base URL",
	)
	flag.Parse()

	setDefaultServerAddress()
	setDefaultBaseURL()
}

func setDefaultServerAddress() {
	if ServerConfig.ServerAddress.String() == "" {
		serverURL := *api.NewServerURL("", "localhost", 8080)
		ServerConfig.ServerAddress = serverURL
	}
}

func setDefaultBaseURL() {
	if ServerConfig.ServerAddress.String() == "" {
		serverURL := *api.NewServerURL("http", "localhost", 8080)
		ServerConfig.ServerAddress = serverURL
	}
}
