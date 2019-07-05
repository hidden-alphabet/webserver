package model

import (
	"database/sql"
)

type Contact struct {
	ID                int
	AccountID         int
	Email             string
	HasConfirmedEmail bool
}

func (c *Contact) Create(tx *sql.Tx) error {
	query := "" +
		"INSERT INTO user.contact (account_id, email, has_confirmed_email)" +
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
		"UPDATE user.account AS ua " +
		"SET has_confirmed_email = $1 " +
		"FROM user.session AS us " +
		"WHERE ua.id = us.account_id " +
		"AND ua.email = $2 " +
		"AND us.token = $3"

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
		"UPDATE user.account AS ua " +
		"SET email = $1 " +
		"FROM user.session AS us " +
		"WHERE ua.id = us.account_id " +
		"AND us.token = $2"

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	confirmation, ok := bool.(req.New)
	if !ok {
		return errors.New("New value is not a boolean.")
	}

	_, err = stmt.Exec(confirmation, req.SessionToken)
	if err != nil {
		return err
	}

	return nil
}
