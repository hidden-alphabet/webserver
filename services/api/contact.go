package api

import (
	"database/sql"
	"net/http"

	"io/ioutil"
)

func (api *API) HandleUpdateEmail(w http.ResponseWriter, r *http.Request) {
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

	if req.New == "" {
		http.Error(w, "No new email given.", http.StatusPreconditionFailed)
		return
	}

	if req.Old == "" {
		http.Error(w, "Unsure which email to change.", http.StatusPreconditionFailed)
		return
	}

	tx, err := api.database.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	contact := model.Contact{}
	err = contact.UpdateEmail(req, tx)
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
