package main

import (
	"github.com/ShishkovEM/amazing-shortener/cmd/pkg/linkserver"
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	server := linkserver.New()

	mux.HandleFunc("/", server.LinkHandler)

	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}
