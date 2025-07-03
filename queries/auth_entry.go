package queries

import (
	"database/sql"
	"strconv"

	project_errors "github.com/saeedsamimi/authentication-service/errors"
)

var (
	CreateAuthEntryQuery = `INSERT INTO auth_entries (user_id, email, password) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at, last_login;`
	GetAuthEntryQuery    = func(fields []string) string {
		baseQuery := `SELECT id, user_id, email, password, created_at, updated_at, last_login FROM auth_entries WHERE `
		for i, field := range fields {
			if i > 0 {
				baseQuery += " AND "
			}
			baseQuery += field + ` = $` + strconv.Itoa(i+1)
		}
		return baseQuery + `;`
	}
)

func CreateAuthEntry(db *sql.DB) (*sql.Stmt, error) {
	stmt, err := db.Prepare(CreateAuthEntryQuery)
	if err != nil {
		return nil, &project_errors.PrepareError{
			Query: CreateAuthEntryQuery,
			Err:   err,
		}
	}
	return stmt, nil
}

func GetAuthEntryBy(db *sql.DB, fields []string) (*sql.Stmt, error) {
	query := GetAuthEntryQuery(fields)
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, &project_errors.PrepareError{
			Query: query,
			Err:   err,
		}
	}
	return stmt, nil
}
