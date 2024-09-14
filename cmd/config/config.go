package config

import (
	"flag"
	"github.com/pervukhinpm/link-shortener.git/internal/api"
	"os"
	"strconv"
	"strings"
)

var ServerConfig struct {
	ServerAddress   api.ServerURL
	BaseURL         api.ServerURL
	FileStoragePath string
	DatabaseDSN     string
}

func ParseFlags() {
	var flagServerAddress string
	var flagBaseURL string
	var flagFileStoragePath string
	var flagDatabaseDSN string

	flag.StringVar(&flagServerAddress, "a", "localhost:8080", "Host Port")
	flag.StringVar(&flagBaseURL, "b", "http://localhost:8080/", "Base URL")
	flag.StringVar(&flagFileStoragePath, "f", "/tmp/url-db.json", "File storage path")
	flag.StringVar(&flagDatabaseDSN, "d", "", "Database DSN")

	flag.Parse()

	if envServerAddressEnv := os.Getenv("SERVER_ADDR"); envServerAddressEnv != "" {
		flagServerAddress = envServerAddressEnv
	}

	if baseURLEnv := os.Getenv("BASE_URL"); baseURLEnv != "" {
		flagBaseURL = baseURLEnv
	}

	if fileStoragePathEnv := os.Getenv("FILE_STORAGE_PATH"); fileStoragePathEnv != "" {
		flagFileStoragePath = fileStoragePathEnv
	}

	if databaseDSNEnv := os.Getenv("DATABASE_DSN"); databaseDSNEnv != "" {
		flagDatabaseDSN = databaseDSNEnv
	}

	ServerConfig.ServerAddress = *parseServerURL(flagServerAddress)
	ServerConfig.BaseURL = *parseServerURL(flagBaseURL)
	ServerConfig.FileStoragePath = flagFileStoragePath
	ServerConfig.DatabaseDSN = flagDatabaseDSN
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

	if len(hostPortParts) != 2 {
		return api.NewServerURL(scheme, host, 8080)
	}

	var port int
	var err error
	port, err = strconv.Atoi(hostPortParts[1])
	if err != nil {
		port = 8080
	}

	return api.NewServerURL(scheme, host, port)
}
