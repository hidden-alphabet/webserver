package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/argon2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Server struct {
	router   *mux.Router
	database *sql.DB
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (s *Server) HandleUser(w http.ResponseWriter, r *http.Request) {
	user := User{}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}

	err = json.Unmarshal(data, &user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}

	tx, err := s.database.Begin()
	if err != nil {
		log.Println(err)
	}
	defer tx.Rollback()

	err = func(u *User) error {
		switch r.Method {
		case "POST":
			return user.Create(tx)
			//      case "PUT": return user.Update(tx)
			//      case "DELETE": return user.Delete(tx)
			//      case "GET": fallthrough
			//      default: return user.Read(tx)
		}
		return nil
	}(&user)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusExpectationFailed)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (u *User) Create(tx *sql.Tx) error {
	salt := make([]byte, 2056)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}

	hash := argon2.IDKey([]byte(u.Password), salt, 1, 64*1024, 4, 32)

	query := "INSERT INTO web.user (name, email, hash, salt) VALUES ($1, $2, $3, $4)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(u.Username, u.Email, hash, salt)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "__hiddenalphabet_session", Value: "astaxie", Expires: expiration}
	http.SetCookie(w, &cookie)

	return nil
}

func (u *User) Update(tx *sql.Tx) error {
	query := "UPDATE web.user SET  WHERE condition;"

	return nil
}

func resolveDatabaseURL(env string) string {
	template := "postgres://%s:%s@%s:%s/%s?sslmode=disable"

	switch env {
	case "production":
		host := os.Getenv("PG_HOST")
		port := os.Getenv("PG_PORT")
		username := os.Getenv("PG_USERNAME")
		password := os.Getenv("PG_PASSWORD")
		database := os.Getenv("PG_DATABASE")

		return fmt.Sprintf(template, username, password, host, port, database)
	case "development":
		fallthrough
	default:
		username, exists := os.LookupEnv("PG_USERNAME")
		if !exists {
			username = "postgres"
		}
		password, exists := os.LookupEnv("PG_PASSWORD")
		if !exists {
			password = ""
		}
		database, exists := os.LookupEnv("PG_DATABASE")
		if !exists {
			database = "postgres"
		}

		return fmt.Sprintf(template, username, password, "localhost", "5432", database)
	}
}

func main() {
	env, exists := os.LookupEnv("ENV")
	if !exists {
		log.Println("Could not determine environment. Defaulting to development.")
		env = "development"
	}

	url := resolveDatabaseURL(env)
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to database.")

	router := mux.NewRouter()

	s := &Server{router, db}

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/user", s.HandleUser).Methods("POST")

	http.Handle("/", handlers.LoggingHandler(os.Stdout, router))
	log.Println("Server started at 0.0.0.0:8081.")

	http.ListenAndServe(":8081", nil)
}
