package main

import (
  "os"
  "log"
  "fmt"
  "net/http"
  "io/ioutil"
  "crypto/rand"
  "database/sql"
  "encoding/json"
  _ "github.com/lib/pq"
  "github.com/gorilla/mux"
  "golang.org/x/crypto/argon2"
  "github.com/gorilla/handlers"
)

type Server struct {
  router *mux.Router
  database *sql.DB
}

type User struct {
  Username string `json:"username"`
  Password string `json:"password"`
  Email    string `json:"email"`
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
  user := User{}
  data, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Println("Error while reading POST body.")
    log.Println(err)
    w.WriteHeader(http.StatusExpectationFailed)
    return
  }

  log.Println("Creating user")

  err = json.Unmarshal(data, &user)
  if err != nil {
    log.Println("Error while parsing JSON.")
    log.Println(err)
    w.WriteHeader(http.StatusExpectationFailed)
    return
  }

  salt := make([]byte, 2056)
  _, err = rand.Read(salt)
  if err != nil {
    log.Println("Error while generating salt.")
    log.Println(err)
    w.WriteHeader(http.StatusExpectationFailed)
    return
  }

  hash := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)

  query := "INSERT INTO web.user (name, email, hash, salt) VALUES ($1, $2, $3, $4)"
  _, err = s.database.Query(query, user.Username, user.Email, hash, salt)
  if err != nil {
    log.Println("Error exeuting SQL query.")
    log.Println(err)
    w.WriteHeader(http.StatusExpectationFailed)
    return
  }

  w.WriteHeader(http.StatusOK)
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
  api.HandleFunc("/user", s.CreateUser).Methods("POST")

  http.Handle("/", handlers.LoggingHandler(os.Stdout, router))
  log.Println("Server started at 0.0.0.0:8081.")

  http.ListenAndServe(":8081", nil)
}
