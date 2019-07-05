package api

import (
	"encoding/json"
	_ "github.com/lib/pq"
	"io/ioutil"
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

	err = json.Unmarshal(data, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tx, err := api.database.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	user, err := model.NewUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := user.Create(tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contact = model.Contact{
		AccountID: id,
		Email:     req.Email,
	}
	err = contact.Create(tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session := model.Session{}
	token, err := session.Create(tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	oneYearFromNow := time.Now().Add((356 / 2) * 24 * time.Hour)

	http.SetCookie(w, &http.Cookie{
		Name:     "__hiddenalphabet_session",
		Value:    token,
		Expires:  oneYearFromNow,
		Secure:   true,
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(APIResponse{Status: "successful"})

	w.WriteHeader(http.StatusOK)
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

	req := model.UpdateRequest{SessionToken: cookie.Value()}

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
	err = user.UpdatePassword(req)
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
