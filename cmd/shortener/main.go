package main

import (
	"github.com/pervukhinpm/link-shortener.git/cmd/config"
	"github.com/pervukhinpm/link-shortener.git/internal/api"
	"github.com/pervukhinpm/link-shortener.git/internal/repository"
	"github.com/pervukhinpm/link-shortener.git/internal/url"
	"log"
)

func main() {
	config.ParseFlags()
	inMemoryRepository := repository.NewInMemoryRepository()
	urlService := url.NewURLService(inMemoryRepository)
	httpHandler := api.NewHandler(urlService, &config.ServerConfig.BaseUrl)
	server := api.NewServer(&config.ServerConfig.ServerAddress, httpHandler)
	log.Fatal(server.Start())
}
