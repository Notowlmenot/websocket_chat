package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB(connStr string) error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("Error connecting to DB: %v", err)
	}
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("error connecting to DB: %w", err)
	}
	defer db.Close()

	log.Printf("Connected to database succesfully")

	DB = db
	return nil
}
