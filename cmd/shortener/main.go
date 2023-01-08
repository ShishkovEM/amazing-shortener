package main

import (
	"github.com/ShishkovEM/amazing-shortener/cmd/pkg/linkserver"
	"github.com/ShishkovEM/amazing-shortener/internal/app/linkstore"
)

func main() {
	linkStorage := linkstore.New()
	server := linkserver.NewLinkServer(linkStorage)
	server.Run()
}
