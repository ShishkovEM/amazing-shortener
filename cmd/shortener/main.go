package main

import (
	"log"
	"net/http"

	"github.com/ShishkovEM/amazing-shortener/internal/app/service"
	"github.com/ShishkovEM/amazing-shortener/internal/app/storage"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
)

type LinkServiceConfig struct {
	Address string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL string `env:"BASE_URL" envDefault:"localhost:8080"`
}

func main() {
	var cfg LinkServiceConfig
	err := env.Parse(&cfg)
	if err != nil {
		log.Printf("Error parsing linservice config: %s\n", err)
		return
	}

	linkStorage := storage.NewLinkStore()
	linkService := service.NewLinkService(linkStorage, "http://"+cfg.BaseURL+"/")
	router := chi.NewRouter()
	router.Mount("/", linkService.Routes())
	router.Mount("/api", linkService.RestRoutes())
	err = http.ListenAndServe(cfg.Address, router)
	if err != nil {
		log.Printf("Error starting linkService: %s\n", err)
		return
	}
}
