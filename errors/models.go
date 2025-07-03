package errors

import "fmt"

const (
	// ErrCodePrepare is the error code for prepare errors.
	ErrCodePrepare = "prepare_error"

	ErrCodeAleadyExists = "already_exists"
	ErrCodeNotFound     = "not_found"
)

type ModelError struct {
	Code  string
	Model string
	Err   error
}

func (e *ModelError) Error() string {
	return fmt.Sprintf("ModelError: Code=%s, Model=%s, Err=%v", e.Code, e.Model, e.Err)
}
