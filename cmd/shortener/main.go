package main

import (
	"log"
	"net/http"

	"github.com/ShishkovEM/amazing-shortener/cmd/pkg/linkservice"
	"github.com/ShishkovEM/amazing-shortener/internal/app/linkstore"
)

const (
	serverAddress = "http://localhost"
	serverPort    = "8080"
)

func main() {
	linkStorage := linkstore.NewLinkStore()
	service := linkservice.NewLinkService(linkStorage)
	log.Fatal(
		http.ListenAndServe(
			service.Run(serverAddress, serverPort), service.Routes(),
		),
	)
}
