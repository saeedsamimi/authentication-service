package repositories_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	project_errors "github.com/saeedsamimi/authentication-service/errors"
	"github.com/saeedsamimi/authentication-service/models"
	"github.com/saeedsamimi/authentication-service/repositories"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type mockedAuthEntryModel struct {
	t        *testing.T
	DoCreate func(t *testing.T, entry models.AuthEntryCreate) (*models.AuthEntry, error)
	DoGet    func(t *testing.T, query models.AuthEntryQuery) (*models.AuthEntry, error)
}

func (m *mockedAuthEntryModel) Create(entry models.AuthEntryCreate) (*models.AuthEntry, error) {
	return m.DoCreate(m.t, entry)
}

func (m *mockedAuthEntryModel) Get(query models.AuthEntryQuery) (*models.AuthEntry, error) {
	return m.DoGet(m.t, query)
}

func TestAuthEntryRepository(t *testing.T) {
	mockedModel := &mockedAuthEntryModel{
		t: t,
	}

	repo := repositories.NewAuthEntryRepository(mockedModel)

	t.Run("Create", func(t *testing.T) {
		expectedOutput := &models.AuthEntry{
			ID:        "test-user-id",
			UserId:    "test-userId",
			Email:     "example@email.com",
			Password:  "SecurePassword123",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			LastLogin: sql.NullTime{},
		}

		mockedModel.DoCreate = func(t *testing.T, entry models.AuthEntryCreate) (*models.AuthEntry, error) {
			assert.Equal(t, expectedOutput.UserId, entry.UserId, "Expected UserId to match")
			assert.Equal(t, expectedOutput.Email, entry.Email, "Expected Email to match")
			assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(entry.Password), []byte(expectedOutput.Password)), "Expected Password to be hashed correctly")
			expectedOutput.Password = "..."
			return &models.AuthEntry{
				ID:        entry.UserId,
				UserId:    entry.UserId,
				Email:     entry.Email,
				Password:  entry.Password,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				LastLogin: sql.NullTime{},
			}, nil
		}

		result, err := repo.Create(&models.AuthEntryCreate{
			UserId:   expectedOutput.UserId,
			Email:    expectedOutput.Email,
			Password: expectedOutput.Password,
		})

		assert.NoError(t, err, "Expected no error on Create")
		assert.Equal(t, expectedOutput.UserId, result.UserId, "Expected UserId to match")
	})

	t.Run("Create_AlreadyExists", func(t *testing.T) {
		mockedModel.DoCreate = func(t *testing.T, entry models.AuthEntryCreate) (*models.AuthEntry, error) {
			return nil, &project_errors.ModelError{
				Code:  project_errors.ErrCodeAlreadyExists,
				Model: "AuthEntryModel",
				Err:   nil,
			}
		}

		_, err := repo.Create(&models.AuthEntryCreate{
			UserId:   "test-user-id",
			Email:    "example.com",
			Password: "SecurePassword123",
		})

		assert.Error(t, err, "Expected error on Create for already existing entry")
		var repoErr *project_errors.RepositoryError
		assert.ErrorAs(t, err, &repoErr, "Expected error to be of type RepositoryError")
		assert.Equal(t, project_errors.ErrCodeAlreadyExists, repoErr.Code, "Expected error code to be ErrCodeAlreadyExists")
	})

	t.Run("Create_ProcessError", func(t *testing.T) {
		mockedModel.DoCreate = func(t *testing.T, entry models.AuthEntryCreate) (*models.AuthEntry, error) {
			return nil, fmt.Errorf("process error")
		}

		_, err := repo.Create(&models.AuthEntryCreate{
			UserId:   "test-user-id",
			Email:    "example.com",
			Password: "SecurePassword123",
		})

		assert.Error(t, err, "Expected error on Create for process error")
		var repoErr *project_errors.RepositoryError
		assert.ErrorAs(t, err, &repoErr, "Expected error to be of type RepositoryError")
		assert.Equal(t, project_errors.ErrCodeProcessError, repoErr.Code, "Expected error code to be ErrCodeProcessError")
	})
}
