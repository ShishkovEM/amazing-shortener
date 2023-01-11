package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"

	"github.com/ShishkovEM/amazing-shortener/cmd/pkg/linkservice"
	"github.com/ShishkovEM/amazing-shortener/internal/app/linkstore"
)

const (
	linkServiceHost = "localhost"
	linkServicePort = "8080"
)

func main() {
	linkStorage := linkstore.NewLinkStore()
	linkService := linkservice.NewLinkService(linkStorage, "http://"+linkServiceHost+":"+linkServicePort+"/")
	router := chi.NewRouter()
	router.Mount("/", linkService.Routes())
	err := http.ListenAndServe(
		linkServiceHost+":"+linkServicePort, router,
	)
	if err != nil {
		log.Print("Error starting linkService")
		return
	}
}
