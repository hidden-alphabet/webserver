package server

import (
	"crypto/rand"
	"encoding/json"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/argon2"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (s *Server) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	user := User{}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(data, &user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tx, err := s.database.Begin()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer tx.Rollback()

	salt := make([]byte, 2056)
	_, err = rand.Read(salt)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	hash := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)

	query := "INSERT INTO web.user (name, email, hash, salt) VALUES ($1, $2, $3, $4)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	stmt.Exec(user.Username, user.Email, hash, salt)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{
		Name:     "__hiddenalphabet_session",
		Value:    "test",
		Expires:  expiration,
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
}

func (s *Server) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) HandleReadUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
