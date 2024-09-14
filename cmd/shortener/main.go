package main

import (
	"github.com/pervukhinpm/link-shortener.git/cmd/config"
	"github.com/pervukhinpm/link-shortener.git/internal/api"
	"github.com/pervukhinpm/link-shortener.git/internal/middleware"
	"github.com/pervukhinpm/link-shortener.git/internal/repository"
	"github.com/pervukhinpm/link-shortener.git/internal/url"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	middleware.Initialize()
	config.ParseFlags()

	fileRepository, err := repository.NewFileRepository(config.ServerConfig.FileStoragePath)
	if err != nil {
		log.Fatal("Failed to initialize file repository", err)
	}

	urlService := url.NewURLService(fileRepository)
	httpHandler := api.NewHandler(urlService, config.ServerConfig.BaseURL, config.ServerConfig.DatabaseDSN)
	server := api.NewServer(&config.ServerConfig.ServerAddress, httpHandler)
	log.Fatal(server.Start())
}
