package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	project_errors "github.com/saeedsamimi/authentication-service/errors"
	"github.com/saeedsamimi/authentication-service/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthEntryRepository struct {
	model models.IAuthEntryModel
	cache *redis.Client
}

const name = "AuthEntryRepository"

func NewAuthEntryRepository(model models.IAuthEntryModel, cache *redis.Client) *AuthEntryRepository {
	return &AuthEntryRepository{model: model, cache: cache}
}

func (r *AuthEntryRepository) Create(ctx context.Context, entry *models.AuthEntryCreate) (*models.AuthEntry, error) {
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

	authEntryResult, err := r.model.Create(ctx, authEntry)

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

	hashDetails := map[string]string{
		"ID":        authEntryResult.ID,
		"UserId":    authEntryResult.UserId,
		"Email":     authEntryResult.Email,
		"Password":  string(hashedPassword),
		"CreatedAt": authEntryResult.CreatedAt.Format(time.RFC3339),
		"UpdatedAt": authEntryResult.UpdatedAt.Format(time.RFC3339),
	}

	if authEntryResult.LastLogin.Valid {
		hashDetails["LastLogin"] = authEntryResult.LastLogin.Time.Format(time.RFC3339)
	} else {
		hashDetails["LastLogin"] = ""
	}

	go func(key string, dto *map[string]string) {
		val, err := r.cache.HSet(ctx, key, *dto).Result()
		if err != nil || val != 7 {
			fmt.Println("Error setting cache:", err)
		}
	}(authEntryResult.ID, &hashDetails)

	return authEntryResult, nil
}
