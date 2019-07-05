package server

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/argon2"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

func (s *Server) AuthorizationRequired(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("__hiddenalphabet_session")
		if err != nil {
			switch err {
			case http.ErrNoCookie:
				w.WriteHeader(http.StatusUnauthorized)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			h(w, r)
		}
	}
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (u *User) FromRequest(r *http.Request) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, u)
	if err != nil {
		return err
	}

	return nil
}

/*
  Add a new user.

  Route: /user/create
  Method(s): POST
  Example Request:
    {
      "username": "foo",
      "password": "bar",
      "email": "bax@baz.com
    }
  Example Response:
    {
      "status": "successful"
    }
*/
func (s *Server) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	user := User{}

	err := user.FromRequest(r)
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
		"INSERT INTO user.account (name, email, hash, salt) " +
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
		"INSERT INTO web.meta (account_id, is_active, email_confirmed, email_confirmation_path)" +
		"VALUES ($1, $2, $3, $4)"
	createUserMetadataStmt, err := tx.Prepare(createUserMetadataQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	emailConfirmationKeyData := make([]byte, 64)
	_, err = rand.Read(salt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	emailConfirmationKey := url.QueryEscape(base64.StdEncoding.EncodeToString(emailConfirmationKeyData))
	emailConfirmationPath := fmt.Sprintf("/email/update/confirmation/%s", emailConfirmationKey)

	_, err = createUserMetadataStmt.Exec(id, true, false, emailConfirmationPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	createSessionQuery := "" +
		"INSERT INTO user.session (account_id, active, token) " +
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

/*
  Add a new user.

  Route: /user/update/email
  Method(s): PUT
  Example Request:
    {
      "new": "test@gmail.com"
    }
  Example Response:
    {
      "status": "successful"
    }
*/
func (s *Server) HandleUpdateUserEmail(w http.ResponseWriter, r *http.Request) {
	user := User{}

	err := user.FromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.Email == "" {
		http.Error(w, "No email provided.", http.StatusPreconditionFailed)
		return
	}

	tx, err := s.database.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	updateEmailQuery := "" +
		"UPDATE user.account AS ua " +
		"SET email = $1 " +
		"FROM user.session AS us " +
		"WHERE ua.id = us.account_id " +
		"AND us.token = $2"

	updateEmailStmt, err := tx.Prepare(updateEmailQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer updateEmailStmt.Close()

	cookie, err := r.Cookie("__hiddenalphabet_session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

type PasswordUpdate struct {
	Old string "json:`old`"
	New string "json:`new`"
}

func (p *PasswordUpdate) FromRequest(r *http.Request) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	return nil
}

/*
  Add a new user.

  Route: /user/update/password
  Method(s): PUT
  Example Request:
    {
      "old": "old-password",
      "new": "new-password"
    }
  Example Response:
    {
      "status": "successful"
    }
*/
func (s *Server) HandleUpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	update := PasswordUpdate{}

	err := update.FromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if update.Old == "" || update.New == "" {
		http.Error(w, "Cannot have an empty password.", http.StatusPreconditionFailed)
		return
	}

	tx, err := s.database.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	getUserHashQuery := "" +
		"SELECT hash, salt " +
		"FROM user.account AS ua " +
		"INNER JOIN user.session AS us " +
		"ON ua.id = us.account_id " +
		"WHERE us.token = $1"

	getUserHashStmt, err := tx.Prepare(getUserHashQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie, err := r.Cookie("__hiddenalphabet_session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var oldHash []byte
	var oldSalt []byte
	err = getUserHashStmt.QueryRow(cookie.Value).Scan(&oldHash, &oldSalt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if reflect.DeepEqual(oldHash, argon2.IDKey([]byte(update.Old), oldSalt, 1, 64*1024, 4, 32)) {
		http.Error(w, "Incorrect Password", http.StatusUnauthorized)
		return
	}

	updatePasswordQuery := "" +
		"UPDATE user.account AS ua " +
		"SET hash = $1, salt = $2 " +
		"FROM user.session AS us " +
		"WHERE ua.id = us.account_id " +
		"AND us.token = $3"

	updatePasswordStmt, err := tx.Prepare(updatePasswordQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer updatePasswordStmt.Close()

	newSalt := make([]byte, 2056)
	_, err = rand.Read(newSalt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newHash := argon2.IDKey([]byte(update.New), newSalt, 1, 64*1024, 4, 32)

	_, err = updatePasswordStmt.Exec(newHash, newSalt, cookie.Value)
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
