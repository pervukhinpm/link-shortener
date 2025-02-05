package main

import (
	"github.com/pervukhinpm/link-shortener.git/cmd/config"
	"github.com/pervukhinpm/link-shortener.git/internal/api"
	"github.com/pervukhinpm/link-shortener.git/internal/db"
	"github.com/pervukhinpm/link-shortener.git/internal/middleware"
	"github.com/pervukhinpm/link-shortener.git/internal/repository"
	"github.com/pervukhinpm/link-shortener.git/internal/service"
	"log"
)

import _ "net/http/pprof"

func main() {
	middleware.Initialize()
	config.ParseFlags()

	database, err := db.NewDB(config.ServerConfig.DatabaseDSN)
	if err != nil {
		middleware.Log.Error("Failed to create database: %v", err)
		return
	}

	appRepository, err := repository.NewRepository(
		config.ServerConfig.DatabaseDSN,
		config.ServerConfig.FileStoragePath,
		database,
	)

	if err != nil {
		middleware.Log.Error("Failed to initialize repository: %v", err)
		return
	}

	defer func(appRepository repository.Repository) {
		err := appRepository.Close()
		if err != nil {
			middleware.Log.Error("Failed to close repository: %v", err)
		}
	}(appRepository)

	urlService := service.NewURLService(appRepository)
	shortenerHandler := api.NewHandler(urlService, config.ServerConfig.BaseURL)
	ping := service.NewPingService(database)
	databaseHandler := api.NewDatabaseHealthHandler(ping)
	router := api.Router(databaseHandler, shortenerHandler)
	server := api.NewServer(&config.ServerConfig.ServerAddress, router)

	log.Fatal(server.Start())
}
