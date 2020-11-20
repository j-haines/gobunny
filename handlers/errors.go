package handlers

import "fmt"

// HTTPError is an error wrapping an error message and status code
type HTTPError struct {
	Message string
	Status  int
}

// NewHTTPError returns an HTTPError with the given error message and status code
func NewHTTPError(message string, status int) error {
	return &HTTPError{
		Message: message,
		Status:  status,
	}
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("error %d: '%s'", e.Status, e.Message)
}
