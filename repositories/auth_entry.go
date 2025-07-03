package repositories

import (
	"errors"

	project_errors "github.com/saeedsamimi/authentication-service/errors"
	"github.com/saeedsamimi/authentication-service/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthEntryRepository struct {
	model models.IAuthEntryModel
}

const name = "AuthEntryRepository"

func NewAuthEntryRepository(model models.IAuthEntryModel) *AuthEntryRepository {
	return &AuthEntryRepository{model: model}
}

func (r *AuthEntryRepository) Create(entry *models.AuthEntryCreate) (*models.AuthEntry, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(entry.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, &project_errors.RepositoryError{
			Code:       project_errors.ErrCodeProcessError,
			Repository: name,
			Err:        err,
		}
	}

	authEntry := models.AuthEntryCreate{
		UserId:   entry.UserId,
		Email:    entry.Email,
		Password: string(hashedPassword),
	}

	authEntryResult, err := r.model.Create(authEntry)

	if err != nil {
		var projError *project_errors.ModelError
		if errors.As(err, &projError) {
			if projError.Code == project_errors.ErrCodeAlreadyExists {
				return nil, &project_errors.RepositoryError{
					Code:       project_errors.ErrCodeAlreadyExists,
					Repository: name,
					Err:        projError,
				}
			}
			return nil, &project_errors.RepositoryError{
				Code:       projError.Code,
				Repository: name,
				Err:        projError,
			}
		}
		return nil, &project_errors.RepositoryError{
			Code:       project_errors.ErrCodeProcessError,
			Repository: name,
			Err:        err,
		}
	}

	authEntryResult.Password = "..."

	return authEntryResult, nil
}
