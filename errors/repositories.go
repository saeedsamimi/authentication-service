package project_errors

import "fmt"

type RepositoryError struct {
	Code       string
	Repository string
	Err        error
}

func (e *RepositoryError) Error() string {
	return fmt.Sprintf("ModelError: Code=%s, Repository=%s, Err=%v", e.Code, e.Repository, e.Err)
}
