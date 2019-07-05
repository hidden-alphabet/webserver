package model

import (
	"database/sql"
	"strconv"
)

type Contact struct {
	ID                int
	AccountID         int
	Email             string
	HasConfirmedEmail bool
}

func (c *Contact) Create(tx *sql.Tx) error {
	query := "" +
		"INSERT INTO web.contact (account_id, email, has_confirmed_email)" +
		"VALUES ($1, $2, $3) "

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.AccountID, c.Email, false)
	if err != nil {
		return err
	}

	return nil
}

func (c *Contact) UpdateEmail(req *UpdateRequest, tx *sql.Tx) error {
	query := "" +
		"UPDATE web.account AS ua " +
		"SET email = $1 " +
		"FROM web.session AS us " +
		"WHERE ua.id = us.account_id " +
		"AND us.token = $2"

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(req.New, req.Old, req.SessionToken)
	if err != nil {
		return err
	}

	return nil
}

func (c *Contact) UpdateEmailConfirmation(req *UpdateRequest, tx *sql.Tx) error {
	query := "" +
		"UPDATE web.account AS ua " +
		"SET has_confirmed_email = $1 " +
		"FROM web.session AS us " +
		"WHERE ua.id = us.account_id " +
		"AND ua.email = $2 " +
		"AND us.token = $3"

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	status, err := strconv.ParseBool(req.New)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(status, req.SessionToken)
	if err != nil {
		return err
	}

	return nil
}
