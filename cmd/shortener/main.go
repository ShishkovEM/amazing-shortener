package main

import (
	"log"
	"net/http"

	"github.com/ShishkovEM/amazing-shortener/internal/app/config"
	mw "github.com/ShishkovEM/amazing-shortener/internal/app/middleware"
	"github.com/ShishkovEM/amazing-shortener/internal/app/repository"
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

	// Создаём in-memory хранилище ссылок
	linkStorage := storage.NewLinkStore()

	// Запускаем сервис
	startLinkService(linkStorage)
}

func startLinkService(linkStorage *storage.LinkStore) {

	if lsc.FileStoragePath == "" {
		// Создаём сервис для обработки create- и read- операций
		linkService := service.NewLinkService(linkStorage, lsc.BaseURL+"/")

		// Запускаем маршрутизацию
		router := chi.NewRouter().With(mw.UnzipRequest, mw.ZipResponse)
		router.Mount("/", linkService.Routes())
		router.Mount("/api", linkService.RestRoutes())

		// Запускаем http-сервер
		err := http.ListenAndServe(lsc.Address, router)
		if err != nil {
			log.Printf("Error starting linkService: %s\n", err)
			return
		}
	} else {
		linkFileRepo, err := repository.NewLinkRepository(lsc.FileStoragePath, linkStorage)
		if err != nil {
			log.Printf("Error initializing file repository for linkStorage: %s\n", err)
			return
		}

		// Запускаем синхронизацию in-memory хранилища с файловым
		go linkFileRepo.Refresh(lsc.FileStoragePath)

		// Создаём сервис для обработки create- и read- операций
		linkService := service.NewLinkService(linkFileRepo.InMemory, lsc.BaseURL+"/")

		// Запускаем маршрутизацию
		router := chi.NewRouter().With(mw.UnzipRequest, mw.ZipResponse)
		router.Mount("/", linkService.Routes())
		router.Mount("/api", linkService.RestRoutes())

		// Запускаем http-сервер
		err = http.ListenAndServe(lsc.Address, router)
		if err != nil {
			log.Printf("Error starting linkService: %s\n", err)
			return
		}
	}
}
