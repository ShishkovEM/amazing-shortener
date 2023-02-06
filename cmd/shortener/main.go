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

	// Проверяем, требуется ли подключение к БД
	var DB storage.DB
	var dbModel = &DB

	if lsc.DatabaseDSN != "" {
		dbModel = storage.NewDB(lsc.DatabaseDSN)
		dbModel.CreateTables()

		linkStorage := repository.NewDBURLStorage(dbModel)
		linkService := service.NewStandAloneDBService(linkStorage, lsc.BaseURL+"/")

		// Запускаем маршрутизацию
		router := chi.NewRouter()
		router.Use(mw.GenerateAuthToken())
		router.Mount("/", linkService.Routes())

		// Запускаем http-сервер
		err := http.ListenAndServe(lsc.Address, mw.Conveyor(router, mw.UnzipRequest, mw.ZipResponse))
		if err != nil {
			log.Printf("Error starting linkService: %s\n", err)
			return
		}

	} else {
		db := storage.NewDB(lsc.DatabaseDSN)

		var linkStorage *storage.LinkStore

		// Создаём сервис для работы с БД
		dataBaseService := service.NewDataBaseService(db)

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
		router.Use(mw.GenerateAuthToken())
		router.Mount("/", linkService.Routes())
		router.Mount("/api", linkService.RestRoutes())
		router.Mount("/api/user", linkService.UserLinkRoutes())
		router.Mount("/ping", dataBaseService.Routes())

		// Запускаем http-сервер
		err := http.ListenAndServe(lsc.Address, mw.Conveyor(router, mw.UnzipRequest, mw.ZipResponse))
		if err != nil {
			log.Printf("Error starting linkService: %s\n", err)
			return
		}
	}
}
