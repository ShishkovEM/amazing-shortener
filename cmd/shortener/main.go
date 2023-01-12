package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"

	"github.com/ShishkovEM/amazing-shortener/internal/app/service"
	"github.com/ShishkovEM/amazing-shortener/internal/app/storage"
)

const (
	linkServiceHost = "localhost"
	linkServicePort = "8080"
)

func main() {
	linkStorage := storage.NewLinkStore()
	linkService := service.NewLinkService(linkStorage, "http://"+linkServiceHost+":"+linkServicePort+"/")
	router := chi.NewRouter()
	router.Mount("/", linkService.Routes())
	router.Mount("/api", linkService.RestRoutes())
	err := http.ListenAndServe(
		linkServiceHost+":"+linkServicePort, router,
	)
	if err != nil {
		log.Printf("Error starting linkService: %s\n", err)
		return
	}
}
