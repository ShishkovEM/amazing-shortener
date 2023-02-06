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

	var linkStorage *storage.LinkStore

	// Создаём файловый репозиторий и хранилище ссылок
	if lsc.FileStoragePath != "" {
		linkFileRepository, err := repository.NewLinkFileRepository(lsc.FileStoragePath)
		if err != nil {
			log.Printf("Error creating linkRepository: %s\n", err)
			return
		}
		linkStorage = storage.NewLinkStore(linkFileRepository)
	} else {
		linkStorage = storage.NewLinkStoreInMemory()
	}

	// Создаём сервис для обработки create- и read- операций
	linkService := service.NewLinkService(linkStorage, lsc.BaseURL+"/")

	// Запускаем маршрутизацию
	router := chi.NewRouter()
	router.Use(mw.Authorize())
	router.Mount("/", linkService.Routes())
	router.Mount("/api", linkService.RestRoutes())
	router.Mount("/api/user", linkService.UserLinkRoutes())

	// Запускаем http-сервер
	err := http.ListenAndServe(lsc.Address, mw.Conveyor(router, mw.UnzipRequest, mw.ZipResponse))
	if err != nil {
		log.Printf("Error starting linkService: %s\n", err)
		return
	}
}
