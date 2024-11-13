package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	errs := []struct {
		Code            int
		ExpectedMessage string
	}{
		{Code: http.StatusUnauthorized, ExpectedMessage: "invalid authorization, username or password not accepted"},
		{Code: http.StatusNoContent, ExpectedMessage: "empty list entities"},
		{Code: http.StatusInternalServerError, ExpectedMessage: "internal server error"},
		{Code: http.StatusBadRequest, ExpectedMessage: "invalid request to server"},
		{Code: http.StatusNotFound, ExpectedMessage: "entity not found"},
		{Code: http.StatusUnprocessableEntity, ExpectedMessage: "unsupported entity type"},
		{Code: http.StatusConflict, ExpectedMessage: "user already exists"},
		{Code: http.StatusNotAcceptable, ExpectedMessage: "operation not acceptable"},
		{Code: http.StatusBadGateway, ExpectedMessage: "unrecognized error"},
	}
	for _, e := range errs {
		err := NewError(e.Code)
		assert.Equal(t, e.ExpectedMessage, err.Error())
	}
}
