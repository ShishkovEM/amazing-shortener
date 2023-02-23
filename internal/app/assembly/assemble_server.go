package assembly

import (
	"github.com/ShishkovEM/amazing-shortener/internal/app/config"
	mw "github.com/ShishkovEM/amazing-shortener/internal/app/middleware"
	"github.com/ShishkovEM/amazing-shortener/internal/app/models"
	"github.com/ShishkovEM/amazing-shortener/internal/app/repository"
	"github.com/ShishkovEM/amazing-shortener/internal/app/service"
	"github.com/ShishkovEM/amazing-shortener/internal/app/storage"
	"github.com/ShishkovEM/amazing-shortener/internal/app/workerpool"
	"log"
	"net/http"
	"runtime"

	"github.com/go-chi/chi/v5"
)

var (
	standAloneDatabaseServiceConfigs config.LinkServiceConfigForStandaloneDB
	serviceConfigsWithFileStorage    config.LinkServiceConfigWithFileStorage
)

func AssembleAndStartAppWithStandAloneDB(allConfigs config.LinkServiceConfig) {
	log.Print("Assembling services with stand-alone database...")

	// Считываем конфигерацию для сервиса, работающего с БД
	standAloneDatabaseServiceConfigs.GetConfig(allConfigs)

	// Создаём модель БД
	dbModel := models.NewDB(standAloneDatabaseServiceConfigs.DatabaseDSN)
	err := dbModel.Migrate()
	if err != nil {
		log.Printf("Error applying DB migrations")
		return
	}

	// Создаём воркер-пул для обработки DELETE-запросов
	var allDeletionTasks []*workerpool.DeletionTask
	workerDeletionPool := workerpool.NewDeletionPool(allDeletionTasks, runtime.NumCPU())
	go workerDeletionPool.RunBackground(dbModel)

	// Создаём репозиторий
	linkStorage := repository.NewDBURLStorage(dbModel, workerDeletionPool)

	// Инициализируем сервис
	linkService := service.NewStandAloneDBService(linkStorage, standAloneDatabaseServiceConfigs.BaseURL+"/")

	// Запускаем маршрутизацию
	router := chi.NewRouter()
	router.Use(mw.GenerateAuthToken())
	router.Mount("/", linkService.Routes())

	// Запускаем http-сервер
	log.Print("Starting http server...")
	err = http.ListenAndServe(standAloneDatabaseServiceConfigs.Address, mw.Conveyor(router, mw.UnzipRequest, mw.ZipResponse))
	if err != nil {
		log.Printf("Error starting services with stand-alone database: %s\n", err)
		return
	}
}

func AssembleAndStartAppWithFileStorage(allConfigs config.LinkServiceConfig) {
	log.Print("Assembling services with file and in-memory storage...")

	// Считываем конфигерацию для сервиса, работающего без БД
	serviceConfigsWithFileStorage.GetConfig(allConfigs)

	// Объявляем in-memory хранилище для ссылок
	var linkStorage *storage.LinkStore

	// Создаём файловый репозиторий и хранилище ссылок
	if serviceConfigsWithFileStorage.FileStoragePath != "" {
		linkFileRepository, err := repository.NewLinkFileRepository(serviceConfigsWithFileStorage.FileStoragePath)
		if err != nil {
			log.Printf("Error creating linkRepository: %s\n", err)
			return
		}
		linkStorage = storage.NewLinkStore(linkFileRepository)
	} else {
		linkStorage = storage.NewLinkStoreInMemory()
	}

	// Создаём сервис для обработки create- и read- операций
	linkService := service.NewLinkService(linkStorage, serviceConfigsWithFileStorage.BaseURL+"/")

	// Запускаем маршрутизацию
	router := chi.NewRouter()
	router.Use(mw.GenerateAuthToken())
	router.Mount("/", linkService.Routes())
	router.Mount("/api", linkService.RestRoutes())
	router.Mount("/api/user", linkService.UserLinkRoutes())

	// Запускаем http-сервер
	log.Print("Starting http server...")
	err := http.ListenAndServe(serviceConfigsWithFileStorage.Address, mw.Conveyor(router, mw.UnzipRequest, mw.ZipResponse))
	if err != nil {
		log.Printf("Error starting services with file and in-memory storage: %s\n", err)
		return
	}
}
