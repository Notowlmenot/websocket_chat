package server

import (
	"chat/internal/auth"
	"chat/internal/ws"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	router.GET("/ws", ws.WebsocketHandler)
	router.POST("/register", auth.RegisterHandler)
	router.GET("/auth", auth.AuthHandler)
}

func RunServer(PORT string) error {
	router := gin.Default()
	SetupRoutes(router)

	fmt.Printf("Starting server on PORT: %s\n", PORT)
	err := router.Run(":" + PORT)

	return err
}
