package middleware

import (
	"database/sql"
	"log"

	"github.com/joho/godotenv"
)

func createConnection() *sql.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
