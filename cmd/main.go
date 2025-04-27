package main

import (
	"chat/internal/database"
	"chat/internal/server"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	//Получение данных из env
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
	connectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"), os.Getenv("DB_SSLMODE"),
	)

	//Загрузка БД
	err = database.ConnectDB(connectionString)
	if err != nil {
		log.Fatalf("Error connection to DB: %v", err)
	}

	//Запуск сервера
	err = server.RunServer(PORT)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer database.DB.Close()
}
