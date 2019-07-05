package main

import (
	api "../.."
	"database/sql"
	"fmt"
	"github.com/gorilla/handlers"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"net/http"
	"os"
)

func main() {
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	username := os.Getenv("PG_USERNAME")
	password := os.Getenv("PG_PASSWORD")
	database := os.Getenv("PG_DATABASE")

	template := "postgres://%s:%s@%s:%s/%s?sslmode=disable"
	url := fmt.Sprintf(template, username, password, host, port, database)

	postgres, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
	}

	api := api.New(postgres)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, api.Router)))
}
