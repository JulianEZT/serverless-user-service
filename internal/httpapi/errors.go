package httpapi

import "fmt"

// HTTPError represents an HTTP-level error with status code and message.
type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("http %d: %s", e.StatusCode, e.Message)
}

// NewHTTPError returns an HTTPError.
func NewHTTPError(statusCode int, message string) *HTTPError {
	return &HTTPError{StatusCode: statusCode, Message: message}
}
