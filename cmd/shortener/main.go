package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"net/http"

	"github.com/ShishkovEM/amazing-shortener/internal/app/service"
	"github.com/ShishkovEM/amazing-shortener/internal/app/storage"

	"github.com/go-chi/chi/v5"
)

type LinkServiceConfig struct {
	Address         string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func main() {

	// Считываем конфигурацию с помощью флагов
	var cfg LinkServiceConfig
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "server address")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "base url")
	flag.StringVar(&cfg.FileStoragePath, "f", "", "file storage path")
	flag.Parse()

	// Считываем конфигурацию с помощтю переменных окружения
	var envCfg LinkServiceConfig
	err := env.Parse(&envCfg)
	if err != nil {
		log.Printf("Error parsing linservice config: %s\n", err)
		return
	}

	// Если в переменных окружения переданы какие-то значения, перезапишем конфигурацию с их помощью
	if envCfg.FileStoragePath != "" {
		cfg.FileStoragePath = envCfg.FileStoragePath
	}
	if envCfg.Address != "" {
		cfg.Address = envCfg.Address
	}
	if envCfg.BaseURL != "" {
		cfg.BaseURL = envCfg.BaseURL
	}

	// Создаём сервис обработки ссылок
	linkStorage := storage.NewLinkStore(cfg.FileStoragePath)
	linkService := service.NewLinkService(linkStorage, cfg.BaseURL+"/")

	// Запускаем маршрутизацию
	router := chi.NewRouter()
	router.Mount("/", linkService.Routes())
	router.Mount("/api", linkService.RestRoutes())

	// Запускаем http-сервер
	err = http.ListenAndServe(cfg.Address, router)
	if err != nil {
		log.Printf("Error starting linkService: %s\n", err)
		return
	}
}
