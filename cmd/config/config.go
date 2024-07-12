package config

import (
	"flag"
	"github.com/pervukhinpm/link-shortener.git/internal/api"
	"os"
	"strconv"
	"strings"
)

var ServerConfig struct {
	ServerAddress api.ServerURL
	BaseURL       api.ServerURL
}

func ParseFlags() {
	var flagServerAddress string
	var flagBaseURL string

	flag.StringVar(&flagServerAddress, "a", "localhost:8080", "Host Port")
	flag.StringVar(&flagBaseURL, "b", "http://localhost:8080/", "Base URL")
	flag.Parse()

	if envServerAddressEnv := os.Getenv("SERVER_ADDR"); envServerAddressEnv != "" {
		flagServerAddress = envServerAddressEnv
	}

	if baseURLEnv := os.Getenv("BASE_URL"); baseURLEnv != "" {
		flagBaseURL = baseURLEnv
	}

	ServerConfig.ServerAddress = *parseServerURL(flagServerAddress)
	ServerConfig.BaseURL = *parseServerURL(flagBaseURL)
}

func parseServerURL(rawURL string) *api.ServerURL {
	scheme := ""
	if strings.HasPrefix(rawURL, "http://") {
		scheme = "http"
	}
	if strings.HasPrefix(rawURL, "https://") {
		scheme = "https"
	}

	rawURL = strings.TrimPrefix(rawURL, "http://")
	rawURL = strings.TrimPrefix(rawURL, "https://")

	parts := strings.SplitN(rawURL, "/", 2)
	hostPort := parts[0]

	hostPortParts := strings.Split(hostPort, ":")
	host := hostPortParts[0]

	var port int
	if len(hostPortParts) == 2 {
		var err error
		port, err = strconv.Atoi(hostPortParts[1])
		if err != nil {
			port = 8080
		}
	} else {
		port = 8080
	}

	return api.NewServerURL(scheme, host, port)
}
