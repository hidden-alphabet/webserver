package api

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
)

type API struct {
	Router   *mux.Router
	database *sql.DB
}

type APIResponse struct {
	Status string
}

func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func New(database *sql.DB) *API {
	router := mux.NewRouter()

	api := &API{
		Router:   router,
		database: database,
	}

	router.HandleFunc("/", index)
	router.HandleFunc("/user/create", api.HandleCreateUser).Methods("POST")
	router.HandleFunc("/user/update/password", api.HandleUpdatePassword).Methods("PUT")
	router.HandleFunc("/contact/update/email", api.HandleUpdateEmail).Methods("PUT")
	router.HandleFunc("/contact/update/confirmation", api.HandleUpdateEmail).Methods("GET")

	return api
}
