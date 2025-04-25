package server

import (
	"chat/internal/handlers"
	"chat/internal/ws"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Define routes here
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	router.GET("/register", handlers.Register)
	router.GET("/ws", ws.WebsocketHandler)
}

func RunServer(PORT string) error {
	router := gin.Default()
	SetupRoutes(router)

	fmt.Printf("Starting server on PORT: %s\n", PORT)
	err := router.Run(":" + PORT)

	return err
}
