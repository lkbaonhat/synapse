package apierror

import "net/http"

// APIError is a domain error that carries an HTTP status code and a
// user-facing message. Handlers should attach these via c.Error(); the
// ErrorHandler middleware translates them to JSON responses.
type APIError struct {
	Code    int
	Message string
}

func (e *APIError) Error() string { return e.Message }

// Constructors for common errors.

func NotFound(msg string) *APIError          { return &APIError{Code: http.StatusNotFound, Message: msg} }
func BadRequest(msg string) *APIError        { return &APIError{Code: http.StatusBadRequest, Message: msg} }
func Unauthorized(msg string) *APIError      { return &APIError{Code: http.StatusUnauthorized, Message: msg} }
func Forbidden(msg string) *APIError         { return &APIError{Code: http.StatusForbidden, Message: msg} }
func Conflict(msg string) *APIError          { return &APIError{Code: http.StatusConflict, Message: msg} }
func UnprocessableEntity(msg string) *APIError {
	return &APIError{Code: http.StatusUnprocessableEntity, Message: msg}
}
func Internal(msg string) *APIError {
	return &APIError{Code: http.StatusInternalServerError, Message: msg}
}
