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
		router:   mux.NewRouter(),
		database: database,
	}
}

func (s *Server) Start() {
	s.router.HandleFunc("/user/create", s.HandleCreateUser).Methods("POST")
	s.router.HandleFunc("/user/update/email", s.SessionRequired(s.HandleUpdateUserEmail)).Methods("PUT")
	s.router.HandleFunc("/user/update/password", s.SessionRequired(s.HandleUpdateUserPassword)).Methods("PUT")
	s.router.HandleFunc("/user/delete", s.SessionRequired(s.HandleDeleteUser)).Methods("DELETE")
	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Println("Server started at 0.0.0.0:8081.")

	log.Fatal(http.ListenAndServe(":8081", handlers.LoggingHandler(os.Stdout, s.router)))
}
