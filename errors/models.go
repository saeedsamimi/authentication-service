package project_errors

import "fmt"

type ModelError struct {
	Code  string
	Model string
	Err   error
}

func (e *ModelError) Error() string {
	return fmt.Sprintf("ModelError: Code=%s, Model=%s, Err=%v", e.Code, e.Model, e.Err)
}
