package server

import (
	"database/sql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

var database *sql.DB

type Server struct {
	router   *mux.Router
	database *sql.DB
}

func New(database *sql.DB) *Server {
  return &Server{
    router: mux.NewRouter(),
    database: database,
  }
}

func (s *Server) Start() {
	api := s.router.PathPrefix("/api").Subrouter()

	api.HandleFunc("/user/create", s.HandleCreateUser).Methods("POST")
  api.HandleFunc("/user/read", s.HandleReadUser).Methods("GET")
  api.HandleFunc("/user/update", s.HandleUpdateUser).Methods("PUT")
  api.HandleFunc("/user/delete", s.HandleDeleteUser).Methods("DELETE")
  s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
  })

	log.Println("Server started at 0.0.0.0:8081.")

	log.Fatal(http.ListenAndServe(":8081", handlers.LoggingHandler(os.Stdout, s.router)))
}
