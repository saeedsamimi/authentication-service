package queries

import (
	"database/sql"

	"github.com/saeedsamimi/authentication-service/errors"
)

const (
	CreateAuthEntryQuery = `INSERT INTO auth_entries (user_id, email, password) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at, last_login;`
)

func CreateAuthEntry(db *sql.DB) (*sql.Stmt, error) {
	stmt, err := db.Prepare(CreateAuthEntryQuery)
	if err != nil {
		return nil, &errors.PrepareError{
			Query: CreateAuthEntryQuery,
			Err:   err,
		}
	}
	return stmt, nil
}
