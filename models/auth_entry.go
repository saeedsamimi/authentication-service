package models

import (
	"database/sql"
	"time"

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

func NewAuthEntryModel(db *sql.DB) *AuthEntryModel {
	return &AuthEntryModel{DB: db}
}

func (m *AuthEntryModel) Create(entry AuthEntryCreate) (*AuthEntry, error) {
	stmt, err := queries.CreateAuthEntry(m.DB)
	if err != nil {
		return nil, err
	}

	var authEntry AuthEntry = AuthEntry{
		UserId:   entry.UserId,
		Email:    entry.Email,
		Password: entry.Password,
	}

	err = stmt.QueryRow(entry.UserId, entry.Email, entry.Password).Scan(
		&authEntry.ID,
		&authEntry.CreatedAt,
		&authEntry.UpdatedAt,
		&authEntry.LastLogin,
	)

	if err != nil {
		return nil, err
	}

	return &authEntry, nil
}
