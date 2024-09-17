package main

import (
	"github.com/pervukhinpm/link-shortener.git/cmd/config"
	"github.com/pervukhinpm/link-shortener.git/internal/api"
	"github.com/pervukhinpm/link-shortener.git/internal/db"
	"github.com/pervukhinpm/link-shortener.git/internal/middleware"
	"github.com/pervukhinpm/link-shortener.git/internal/repository"
	"github.com/pervukhinpm/link-shortener.git/internal/service"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	middleware.Initialize()
	config.ParseFlags()

	database, err := db.NewDB(config.ServerConfig.DatabaseDSN)

	appRepository, err := repository.NewRepository(
		config.ServerConfig.DatabaseDSN,
		config.ServerConfig.FileStoragePath,
		database,
	)

	if err != nil {
		log.Fatal("Failed to initialize repository", err)
	}

	defer func(appRepository repository.Repository) {
		err := appRepository.Close()
		if err != nil {
			log.Fatal("Failed to close repository", err)
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
