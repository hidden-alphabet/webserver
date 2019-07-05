package model

import (
	"database/sql"
	"github.com/satori/go.uuid"
)

type Session struct {
	ID          int
	AccountID   int
	Token       []byte
	Active      bool
	CreatedAt   string
	CompletedAt string
}

func (s *Session) Create(tx *sql.Tx) (string, error) {
	query := "" +
		"INSERT INTO user.session (account_id, token, active) " +
		"VALUES ($1, $2, $3)"

	token, err := uuid.NewV4()
	if err != nil {
		return make(string), err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return make(string), err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.AccountID, token.String(), false)
	if err != nil {
		return make(string), err
	}

	return token.String(), nil
}

func (s *Session) Delete(req *UpdateRequest, tx *sql.Tx) error {
	query := "" +
		"UPDATE user.session AS as " +
		"SET active = $1, completed_at = now() " +
		"WHERE as.token = $2 "

	stmt, err := tx.Prepare(query)
	if err != nil {
		return make(string), err
	}
	defer stmt.Close()

	active, ok := bool.(req.New)
	if !ok {
		return make(string), err
	}

	_, err = stmt.Exec(active, req.SessionToken, false)
	if err != nil {
		return "", err
	}

	return token.String(), nil
}
