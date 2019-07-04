package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"./server"
)

func main() {
	var postgresURL string

	template := "postgres://%s:%s@%s:%s/%s?sslmode=disable"

	env, exists := os.LookupEnv("ENV")
	if !exists {
		log.Println("Could not determine environment. Defaulting to development.")
		env = "development"
	}

	switch env {
	case "production":
		host := os.Getenv("PG_HOST")
		port := os.Getenv("PG_PORT")
		username := os.Getenv("PG_USERNAME")
		password := os.Getenv("PG_PASSWORD")
		database := os.Getenv("PG_DATABASE")

		postgresURL = fmt.Sprintf(template, username, password, host, port, database)
	case "development":
		fallthrough
	default:
		username, exists := os.LookupEnv("PG_USERNAME")
		if !exists {
			username = "postgres"
		}
		password, exists := os.LookupEnv("PG_PASSWORD")
		if !exists {
			password = ""
		}
		database, exists := os.LookupEnv("PG_DATABASE")
		if !exists {
			database = "postgres"
		}

		postgresURL = fmt.Sprintf(template, username, password, "localhost", "5432", database)
	}

	database, err := sql.Open("postgres", postgresURL)
	if err != nil {
		log.Fatal(err)
	}

	s := server.New(database)
	s.Start()
}
