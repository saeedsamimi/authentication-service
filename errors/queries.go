package project_errors

type PrepareError struct {
	Query string
	Err   error
}

func (e *PrepareError) Error() string {
	return "failed to prepare query: " + e.Query + ": " + e.Err.Error()
}
