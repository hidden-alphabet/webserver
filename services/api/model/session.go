package model

import (
	"database/sql"
	"github.com/satori/go.uuid"
	"strconv"
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
		"INSERT INTO web.session (account_id, token, active) " +
		"VALUES ($1, $2, $3)"

	token, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	_, err = stmt.Exec(s.AccountID, token.Bytes(), true)
	if err != nil {
		return "", err
	}

	return token.String(), nil
}

func (s *Session) Delete(req *UpdateRequest, tx *sql.Tx) error {
	query := "" +
		"UPDATE web.session AS as " +
		"SET active = $1, completed_at = now() " +
		"WHERE as.token = $2 "

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	active, err := strconv.ParseBool(req.New)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(active, req.SessionToken, false)
	if err != nil {
		return err
	}

	return nil
}
