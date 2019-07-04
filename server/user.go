package server

import (
	"crypto/rand"
	"encoding/json"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/argon2"
	"io/ioutil"
	"net/http"
	"time"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func userFromHttpRequest(r *http.Request) (*User, error) {
	user := User{}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Server) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	user, err := userFromHttpRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	salt := make([]byte, 2056)
	_, err = rand.Read(salt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hash := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)

	tx, err := s.database.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	createUserQuery := "" +
		"INSERT INTO web.user (name, email, hash, salt) " +
		"VALUES ($1, $2, $3, $4) " +
		"RETURNING id"
	createUserStmt, err := tx.Prepare(createUserQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer createUserStmt.Close()

	var id int64
	err = createUserStmt.QueryRow(user.Username, user.Email, hash, salt).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	createUserMetadataQuery := "" +
		"INSERT INTO web.meta (user_id, is_active, email_confirmed, email_confirmation_path)" +
		"VALUES ($1, $2, $3, $4)"
	createUserMetadataStmt, err := tx.Prepare(createUserMetadataQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = createUserMetadataStmt.Exec(id, true, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	createSessionQuery := "" +
		"INSERT INTO web.session (user_id, active, token) " +
		"VALUES ($1, $2, $3)"
	createSessionStmt, err := tx.Prepare(createSessionQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer createUserStmt.Close()

	uuidToken, err := uuid.NewV4()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := uuidToken.String()

	_, err = createSessionStmt.Exec(id, true, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	oneYearFromNow := time.Now().Add(356 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:     "__hiddenalphabet_session",
		Value:    token,
		Expires:  oneYearFromNow,
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
}

func (s *Server) HandleUpdateUserEmail(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("__hiddenalphabet_session")
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			w.WriteHeader(http.StatusUnauthorized)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	user, err := userFromHttpRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tx, err := s.database.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	updateEmailQuery := "" +
		"UPDATE web.user AS wu " +
		"SET email = $1 " +
		"FROM web.session AS ws " +
		"WHERE wu.id = ws.user_id " +
		"AND ws.token = $2"

	updateEmailStmt, err := tx.Prepare(updateEmailQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer updateEmailStmt.Close()

	_, err = updateEmailStmt.Exec(user.Email, cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) HandleUpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("__hiddenalphabet_session")
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			w.WriteHeader(http.StatusUnauthorized)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	user, err := userFromHttpRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.Password == "" {

	}

	tx, err := s.database.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	getUserHashQuery := "" +
		"SELECT salt " +
		"FROM web.user AS wu " +
		"INNER JOIN web.session AS ws " +
		"ON wu.id = ws.user_id " +
		"WHERE ws.token = $1"
	getUserHashStmt, err := tx.Prepare(getUserHashQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var salt []byte
	err = getUserHashStmt.QueryRow().Scan(&salt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatePasswordQuery := "" +
		"UPDATE web.user AS wu " +
		"SET password = $1 " +
		"FROM web.session AS ws " +
		"WHERE wu.id = ws.user_id " +
		"AND ws.token = $2 " +
		"AND wu.hash = $3"

	updatePasswordStmt, err := tx.Prepare(updatePasswordQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer updateEmailStmt.Close()

	hash := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)

	_, err = updatePasswordStmt.Exec(user.Email, cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) HandleReadUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
