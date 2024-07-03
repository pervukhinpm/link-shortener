package main

import (
	"github.com/pervukhinpm/link-shortener.git/internal/api"
	"github.com/pervukhinpm/link-shortener.git/internal/repository"
	"github.com/pervukhinpm/link-shortener.git/internal/url"
	"log"
	"net/http"
)

func main() {
	inMemoryRepository := repository.NewInMemoryRepository()
	urlService := url.NewURLService(inMemoryRepository)
	httpHandler := api.NewHandler(urlService)
	log.Fatal(http.ListenAndServe(":8080", httpHandler))
}
