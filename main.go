package main

import (
  "os"
  "log"
  "fmt"
  "net/http"
  "io/ioutil"
  "database/sql"
  _ "github.com/lib/pq"
  "github.com/gorilla/mux"
  "github.com/gorilla/handlers"
)

type Server struct {
  router *mux.Router
  database *sql.DB
}



func resolveDatabaseURL(env string) string {
  template := "postgres://%s:%s@%s:%s/%s?sslmode=disable"

  switch env {
  case "production":
    host := os.Getenv("PG_HOST")
    port := os.Getenv("PG_PORT")
    username := os.Getenv("PG_USERNAME")
    password := os.Getenv("PG_PASSWORD")
    database := os.Getenv("PG_DATABASE")

    return fmt.Sprintf(template, username, password, host, port, database)
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

    return fmt.Sprintf(template, username, password, "localhost", "5432", database)
  }
}

func main() {
  env, exists := os.LookupEnv("ENV")
  if !exists {
    log.Println("Could not determine environment. Defaulting to development.")
    env = "development"
  }

  url := resolveDatabaseURL(env)
  db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
	}

  log.Println("Connected to database.")

  router := mux.NewRouter()

  s := &Server{ router, db }

  api := router.PathPrefix("/api").Subrouter()

  http.Handle("/", handlers.LoggingHandler(os.Stdout, router))
  log.Println("Server started at 0.0.0.0:8081.")

  http.ListenAndServe(":8081", nil)
}
