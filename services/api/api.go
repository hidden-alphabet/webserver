package api

import (
	"github.com/gorilla/mux"
)

type API struct {
	router   *mux.Router
	database *sql.DB
}

func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func New(database *sql.DB) *API {
	api := &API{
		router:   mux.NewRouter(),
		database: database,
	}

	api.router.HandleFunc("/", index)
	api.router.HandleFunc("/user/create", HandleCreateUser).Methods("POST")
	api.router.HandleFunc("/user/update/password", HandleUpdatePassword).Methods("PUT")
	api.router.HandleFunc("/contact/update/email", HandleUpdateEmail).Methods("PUT")
	api.router.HandleFunc("/contact/update/confirmation", HandleUpdateEmail).Methods("GET")

	return api
}
