package v1

import (
	"errors"
	"net/http"
)

var (
	// ErrInternalServerError internal server error
	ErrInternalServerError = NewError("internal server error")
	// ErrSignupBadRequest raised if invalid decode json data request
	ErrSignupBadRequest = NewError("bad signup request")
	// ErrSigninBadRequest raised if invalid decode json data request
	ErrSigninBadRequest = NewError("bad signin request")
	// ErrUserAlreadyExists raised if user already exists
	ErrUserAlreadyExists = NewError("user already exists")
	// ErrUnauthorized raised if user unauthorized
	ErrUnauthorized = NewError("unauthorized user")
	// ErrUserNotFound raised if user not found
	ErrUserNotFound = NewError("user not found")
	// ErrUploadBadRequest raised if invalid decode json data request
	ErrUploadBadRequest = NewError("bad upload request")
	// ErrDownloadBadRequest raised if invalid request to download entity
	ErrDownloadBadRequest = NewError("bad download request")
	// ErrEntityNotFound raised if entity not exists in user space
	ErrEntityNotFound = NewError("entity not found")
	// ErrEntityInvalidType raised if entity not supported binary data
	ErrEntityInvalidType = NewError("invalid entity type")
	// ErrEntityDeleted raise if invalid deleted entity
	ErrEntityDeleted = NewError("error entity deleted")
	// ErrUnsupportedEntityType raised if unknown entity type
	ErrUnsupportedEntityType = NewError("unsupported entity type")
)

// ServerError struct of server errors
type ServerError struct {
	m string
}

// NewError return new server error
func NewError(m string) error {
	return &ServerError{m: m}
}

// Error return error message
func (err *ServerError) Error() string {
	return err.m
}

// GetHTTPStatus return http status equal server error
func (err *ServerError) GetHTTPStatus() int {
	if errors.Is(err, ErrUserAlreadyExists) {
		return http.StatusConflict
	}
	if errors.Is(err, ErrSignupBadRequest) ||
		errors.Is(err, ErrSigninBadRequest) ||
		errors.Is(err, ErrUploadBadRequest) ||
		errors.Is(err, ErrDownloadBadRequest) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrUnauthorized) ||
		errors.Is(err, ErrUserNotFound) {
		return http.StatusUnauthorized
	}
	if errors.Is(err, ErrEntityNotFound) {
		return http.StatusNotFound
	}
	if errors.Is(err, ErrEntityInvalidType) ||
		errors.Is(err, ErrEntityDeleted) {
		return http.StatusNotAcceptable
	}
	if errors.Is(err, ErrUnsupportedEntityType) {
		return http.StatusUnprocessableEntity
	}
	return http.StatusInternalServerError
}
