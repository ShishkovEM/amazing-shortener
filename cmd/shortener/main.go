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

	// Если файловое хранилище не задано, запустим сервис только с in-memory хранилишем
	if lsc.FileStoragePath == "" {
		startLinkServiceWithInMemory(linkStorage)
	} else {
		// Если в конфигурации передано имя файла, инициализруем файловый репозиторий
		startLinkServiceWithFileStorage(linkStorage)
	}
}

func startLinkServiceWithInMemory(linkStorage *storage.LinkStore) {
	// Создаём сервис для обработки create- и read- операций
	linkService := service.NewLinkService(linkStorage, lsc.BaseURL+"/")

	// Запускаем маршрутизацию
	router := chi.NewRouter()
	router.Mount("/", linkService.Routes())
	router.Mount("/api", linkService.RestRoutes())

	// Запускаем http-сервер
	err := http.ListenAndServe(lsc.Address, mw.Conveyor(router, mw.UnzipRequest, mw.ZipResponse))
	if err != nil {
		log.Printf("Error starting linkService: %s\n", err)
		return
	}
}

func startLinkServiceWithFileStorage(linkStorage *storage.LinkStore) {
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
	router := chi.NewRouter()
	router.Mount("/", linkService.Routes())
	router.Mount("/api", linkService.RestRoutes())

	// Запускаем http-сервер
	err = http.ListenAndServe(lsc.Address, mw.Conveyor(router, mw.UnzipRequest, mw.ZipResponse))
	if err != nil {
		log.Printf("Error starting linkService: %s\n", err)
		return
	}
}
