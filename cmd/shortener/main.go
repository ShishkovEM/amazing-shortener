package main

import (
	"github.com/ShishkovEM/amazing-shortener/cmd/pkg/linkserver"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	server := linkserver.New()

	router.POST("/", server.CreateLinkHandler)
	router.GET("/:id", server.GetLinkHandler)

	router.Run("localhost:8080")
}
