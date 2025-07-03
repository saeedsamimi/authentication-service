package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	project_errors "github.com/saeedsamimi/authentication-service/errors"
	"github.com/saeedsamimi/authentication-service/queries"
)

type AuthEntry struct {
	ID        string
	UserId    string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	LastLogin sql.NullTime
}

type AuthEntryCreate struct {
	UserId   string
	Email    string
	Password string
}

type AuthEntryModel struct {
	DB *sql.DB
}

const name = "AuthEntryModel"

func NewAuthEntryModel(db *sql.DB) *AuthEntryModel {
	return &AuthEntryModel{DB: db}
}

func (m *AuthEntryModel) Create(entry AuthEntryCreate) (*AuthEntry, error) {
	query := queries.CreateAuthEntryQuery

	var authEntry AuthEntry = AuthEntry{
		UserId:   entry.UserId,
		Email:    entry.Email,
		Password: entry.Password,
	}

	err := m.DB.QueryRow(query, entry.UserId, entry.Email, entry.Password).Scan(
		&authEntry.ID,
		&authEntry.CreatedAt,
		&authEntry.UpdatedAt,
		&authEntry.LastLogin,
	)

	if err != nil {
		var pqError *pq.Error
		if errors.As(err, &pqError) {
			if pqError.Code.Name() == "unique_violation" { // Unique violation
				return nil, &project_errors.ModelError{
					Code:  project_errors.ErrCodeAlreadyExists,
					Model: name,
					Err:   err,
				}
			}
		}
		return nil, err
	}

	return &authEntry, nil
}

type AuthEntryQuery struct {
	ID     *string
	UserId *string
	Email  *string
}

func (m *AuthEntryModel) Get(query AuthEntryQuery) (*AuthEntry, error) {
	fields := []string{}
	qs := []*string{}

	if query.ID != nil {
		fields = append(fields, "id")
		qs = append(qs, query.ID)
	}
	if query.UserId != nil {
		fields = append(fields, "user_id")
		qs = append(qs, query.UserId)
	}
	if query.Email != nil {
		fields = append(fields, "email")
		qs = append(qs, query.Email)
	}

	if len(fields) == 0 {
		return nil, &project_errors.ModelError{
			Code:  project_errors.ErrCodeInvalidArgument,
			Model: name,
			Err:   fmt.Errorf("at least one field must be specified for query"),
		}
	}

	dbQuery := queries.GetAuthEntryQuery(fields)

	var authEntry AuthEntry
	args := make([]any, len(qs))
	for i, v := range qs {
		args[i] = *v
	}

	err := m.DB.QueryRow(dbQuery, args...).Scan(
		&authEntry.ID,
		&authEntry.UserId,
		&authEntry.Email,
		&authEntry.Password,
		&authEntry.CreatedAt,
		&authEntry.UpdatedAt,
		&authEntry.LastLogin,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &project_errors.ModelError{
				Code:  project_errors.ErrCodeNotFound,
				Model: name,
				Err:   err,
			}
		}
		return nil, err
	}

	return &authEntry, nil
}
