package errors

import "fmt"

// ErrResponseClosed indicates that a http.ResponseWriter was closed before
// the http.HandlerFunc finished handling
type ErrResponseClosed struct {
	wraps error
}

// NewErrResponseClosed returns an ErrResponseClosed instance wrapping `wraps`
func NewErrResponseClosed(wraps error) *ErrResponseClosed {
	return &ErrResponseClosed{
		wraps: wraps,
	}
}

func (e *ErrResponseClosed) Error() string {
	return fmt.Sprintf("responsed closed before handler finished: %s", e.wraps.Error())
}

func (e *ErrResponseClosed) Unwrap() error {
	return e.wraps
}
