package model

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/argon2"
	"reflect"
)

type User struct {
	ID   int
	Name string
	Hash []byte
	Salt []byte
}

func NewUser(username, password string) (*User, error) {
	salt, hash, err := StringToSaltAndHash(password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Name: username,
		Hash: hash,
		Salt: salt,
	}

	return user, nil
}

/* Add a new user to the database. */
func (u *User) Create(tx *sql.Tx) (int, error) {
	var id int

	query := "" +
		"INSERT INTO web.account (name, hash, salt) " +
		"VALUES ($1, $2, $3) " +
		"RETURNING id"

	stmt, err := tx.Prepare(query)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(u.Name, u.Hash, u.Salt).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (u *User) GetHash(token string, tx *sql.Tx) ([]byte, error) {
	var hash []byte

	query := "" +
		"SELECT hash " +
		"FROM web.account AS ua " +
		"INNER JOIN web.session AS us " +
		"ON ua.id = us.account_id " +
		"WHERE us.token = $1"

	stmt, err := tx.Prepare(query)
	if err != nil {
		return hash, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(token).Scan(&hash)
	if err != nil {
		return hash, err
	}

	return hash, nil
}

func (u *User) GetSalt(token string, tx *sql.Tx) (*[]byte, error) {
	var salt []byte

	query := "" +
		"SELECT salt " +
		"FROM web.account AS ua " +
		"INNER JOIN web.session AS us " +
		"ON ua.id = us.account_id " +
		"WHERE us.token = $1"

	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(token).Scan(&salt)
	if err != nil {
		return nil, err
	}

	return &salt, nil
}

func (u *User) GetSaltAndHash(token string, tx *sql.Tx) ([]byte, []byte, error) {
	var salt []byte
	var hash []byte

	query := "" +
		"SELECT hash, salt " +
		"FROM web.account AS ua " +
		"INNER JOIN web.session AS us " +
		"ON ua.id = us.account_id " +
		"WHERE us.token = $1"

	stmt, err := tx.Prepare(query)
	if err != nil {
		return salt, hash, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(token).Scan(&hash, &salt)
	if err != nil {
		return salt, hash, err
	}

	return salt, hash, nil
}

func (u *User) ValidPassword(password string, token string, tx *sql.Tx) (bool, error) {
	oldSalt, oldHash, err := u.GetSaltAndHash(token, tx)
	if err != nil {
		return false, err
	}

	newHash := argon2.IDKey([]byte(password), oldSalt, 1, 64*1024, 4, 32)

	return reflect.DeepEqual(oldHash, newHash), nil
}

func (u *User) UpdatePassword(req *UpdateRequest, tx *sql.Tx) error {
	query := "" +
		"UPDATE web.account AS ua" +
		"SET hash = $1, salt = $2 " +
		"FROM web.session AS us " +
		"WHERE ua.id = us.account_id " +
		"AND us.token = $3"

	valid, err := u.ValidPassword(req.Old, req.SessionToken, tx)
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("Invalid password")
	}

	salt, hash, err := StringToSaltAndHash(req.New)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(hash, salt, req.SessionToken)
	if err != nil {
		return err
	}

	return nil
}
