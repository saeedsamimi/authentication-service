package repositories_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/saeedsamimi/authentication-service/models"
	"github.com/saeedsamimi/authentication-service/repositories"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type mockedAuthEntryModel struct {
	t              *testing.T
	expectedOutput *models.AuthEntry
}

func (m *mockedAuthEntryModel) Create(entry models.AuthEntryCreate) (*models.AuthEntry, error) {
	compareResult := bcrypt.CompareHashAndPassword([]byte(entry.Password), []byte(m.expectedOutput.Password))
	assert.NoError(m.t, compareResult, "Expected password to be hashed")
	return m.expectedOutput, nil
}

func (m *mockedAuthEntryModel) Get(query models.AuthEntryQuery) (*models.AuthEntry, error) {
	return m.expectedOutput, nil
}

func TestAuthEntryRepository(t *testing.T) {
	mockedModel := &mockedAuthEntryModel{
		t: t,
	}

	repo := repositories.NewAuthEntryRepository(mockedModel)

	t.Run("Create", func(t *testing.T) {
		mockedModel.expectedOutput = &models.AuthEntry{
			ID:        "USERID001",
			UserId:    "user123",
			Email:     "email@email.com",
			Password:  "hashed_password",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			LastLogin: sql.NullTime{},
		}

		toInsertEntry := &models.AuthEntryCreate{
			UserId:   mockedModel.expectedOutput.UserId,
			Email:    mockedModel.expectedOutput.Email,
			Password: mockedModel.expectedOutput.Password,
		}

		result, err := repo.Create(toInsertEntry)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		expectedOutput := *mockedModel.expectedOutput
		expectedOutput.Password = "..."

		assert.Equal(t, expectedOutput, *result, "Expected result to match mocked output")
	})
}
