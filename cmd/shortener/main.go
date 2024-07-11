package main

import (
	"github.com/pervukhinpm/link-shortener.git/internal/api"
	"github.com/pervukhinpm/link-shortener.git/internal/repository"
	"github.com/pervukhinpm/link-shortener.git/internal/url"
	"log"
)

func main() {
	inMemoryRepository := repository.NewInMemoryRepository()
	urlService := url.NewURLService(inMemoryRepository)
	httpHandler := api.NewHandler(urlService)
	server := api.NewServer(":8080", httpHandler)
	log.Fatal(server.Start())
}
