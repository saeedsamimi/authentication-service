package models_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	project_errors "github.com/saeedsamimi/authentication-service/errors"
	"github.com/saeedsamimi/authentication-service/models"
	"github.com/stretchr/testify/assert"
)

func TestAuthEntryModel(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	model := models.NewAuthEntryModel(db)

	t.Run("Create", func(t *testing.T) {
		entry := models.AuthEntryCreate{
			UserId:   "user123",
			Email:    "test@example.com",
			Password: "password123",
		}

		expectedID := "auth123"
		expectedTime := time.Now()

		mock.ExpectPrepare("INSERT INTO auth_entries").
			ExpectQuery().
			WithArgs(entry.UserId, entry.Email, entry.Password).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "last_login"}).
				AddRow(expectedID, expectedTime, expectedTime, sql.NullTime{}))

		result, err := model.Create(entry)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedID, result.ID)
		assert.Equal(t, entry.UserId, result.UserId)
		assert.Equal(t, entry.Email, result.Email)
		assert.Equal(t, entry.Password, result.Password)
	})

	t.Run("Create_AlreadyExists", func(t *testing.T) {
		entry := models.AuthEntryCreate{
			UserId:   "user123",
			Email:    "email@email.com",
			Password: "password123",
		}

		mock.ExpectPrepare("INSERT INTO auth_entries").
			ExpectQuery().
			WithArgs(entry.UserId, entry.Email, entry.Password).
			WillReturnError(&pq.Error{Code: "23505"})

		var expectedErr *project_errors.ModelError

		result, err := model.Create(entry)

		assert.ErrorAs(t, err, &expectedErr)
		assert.Equal(t, err.(*project_errors.ModelError).Code, project_errors.ErrCodeAleadyExists)
		assert.Nil(t, result)
	})

	t.Run("Get", func(t *testing.T) {
		userID := "user123"
		email := "test@example.com"
		expectedEntry := &models.AuthEntry{
			ID:        "auth123",
			UserId:    userID,
			Email:     email,
			Password:  "password123",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			LastLogin: sql.NullTime{},
		}

		query := models.AuthEntryQuery{
			UserId: &userID,
			Email:  &email,
		}

		mock.ExpectPrepare("SELECT .+ FROM auth_entries").
			ExpectQuery().
			WithArgs(userID, email).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "email", "password", "created_at", "updated_at", "last_login"}).
				AddRow(expectedEntry.ID, expectedEntry.UserId, expectedEntry.Email, expectedEntry.Password,
					expectedEntry.CreatedAt, expectedEntry.UpdatedAt, expectedEntry.LastLogin))

		result, err := model.Get(query)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedEntry.ID, result.ID)
		assert.Equal(t, expectedEntry.UserId, result.UserId)
		assert.Equal(t, expectedEntry.Email, result.Email)
	})

	t.Run("Get_NoArgs", func(t *testing.T) {
		query := models.AuthEntryQuery{}

		result, err := model.Get(query)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "at least one field must be specified for query")
		assert.Nil(t, result)
	})

	t.Run("Get_NotFound", func(t *testing.T) {
		userID := "nonexistent"
		query := models.AuthEntryQuery{
			UserId: &userID,
		}

		mock.ExpectPrepare("SELECT (.+) FROM auth_entries").
			ExpectQuery().
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		result, err := model.Get(query)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
