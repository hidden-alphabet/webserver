package api

import (
	"./model"
	"encoding/json"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (api *API) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	req := CreateUserRequest{}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Reading JSON")

	err = json.Unmarshal(data, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Beginning postgres transaction")

	tx, err := api.database.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	log.Println("Initializing user")

	user, err := model.NewUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Adding new user to postgres")

	id, err := user.Create(tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Adding user contact information to postgres")

	contact := model.Contact{
		AccountID: id,
		Email:     req.Email,
	}
	err = contact.Create(tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Adding a new session")

	session := model.Session{AccountID: id}
	token, err := session.Create(tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Committing transaction to postgres")

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Creating cookie.")

	oneYearFromNow := time.Now().Add((356 / 2) * 24 * time.Hour)

	http.SetCookie(w, &http.Cookie{
		Name:     "__hiddenalphabet_session",
		Value:    token,
		Expires:  oneYearFromNow,
		Secure:   true,
		HttpOnly: true,
	})

	log.Println("Sending response.")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(APIResponse{Status: "successful"})

	w.WriteHeader(http.StatusOK)

	log.Println("Sending response.")
}

func (api *API) HandleUpdatePassword(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("__hiddenalphabet_session")
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			w.WriteHeader(http.StatusUnauthorized)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	req := model.UpdateRequest{SessionToken: cookie.Value}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(data, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Old == "" || req.New == "" {
		http.Error(w, "Invalid request.", http.StatusPreconditionFailed)
		return
	}

	tx, err := api.database.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	user := model.User{}
	err = user.UpdatePassword(&req, tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(APIResponse{Status: "successful"})

	w.WriteHeader(http.StatusOK)
}
