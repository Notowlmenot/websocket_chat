package main

import (
	"chat/internal/server"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	//Получение порта из env
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}
	parentDir := filepath.Dir(currentDir)
	envFilePath := filepath.Join(parentDir, ".env")
	err = godotenv.Load(envFilePath)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	PORT := os.Getenv("PORT")

	err = server.RunServer(PORT)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
