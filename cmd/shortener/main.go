package main

import (
	"log"
	"net/http"

	"github.com/ShishkovEM/amazing-shortener/internal/app/config"
	"github.com/ShishkovEM/amazing-shortener/internal/app/service"
	"github.com/ShishkovEM/amazing-shortener/internal/app/storage"

	"github.com/go-chi/chi/v5"
)

var (
	lsc config.LinkServiceConfig
)

func main() {

	// Считываем конфигурацию для LinkService
	lsc.Parse()

	// Создаём сервис обработки ссылок
	linkStorage := storage.NewLinkStore(lsc.FileStoragePath)
	linkService := service.NewLinkService(linkStorage, lsc.BaseURL+"/")

	// Запускаем маршрутизацию
	router := chi.NewRouter()
	router.Mount("/", linkService.Routes())
	router.Mount("/api", linkService.RestRoutes())

	// Запускаем http-сервер
	err := http.ListenAndServe(lsc.Address, service.GzipHandle(router))
	if err != nil {
		log.Printf("Error starting linkService: %s\n", err)
		return
	}
}
