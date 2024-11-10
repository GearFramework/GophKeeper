package client

import "net/http"

// Error struct of server errors
type Error struct {
	code int
}

// NewError return new client error
func NewError(httpStatus int) error {
	return &Error{code: httpStatus}
}

// Error return error message
func (err *Error) Error() string {
	if err.code == http.StatusUnauthorized {
		return "invalid authorization, username or password not accepted"
	}
	if err.code == http.StatusNoContent {
		return "empty list entities"
	}
	if err.code == http.StatusInternalServerError {
		return "internal server error"
	}
	if err.code == http.StatusBadRequest {
		return "invalid request to server"
	}
	if err.code == http.StatusNotFound {
		return "entity not found"
	}
	if err.code == http.StatusUnprocessableEntity {
		return "unsupported entity type"
	}
	if err.code == http.StatusConflict {
		return "user already exists"
	}
	if err.code == http.StatusNotAcceptable {
		return "operation not acceptable"
	}
	return "unrecognized error"
}
