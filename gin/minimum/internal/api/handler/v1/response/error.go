package response

import (
	"net/http"
)

type ErrResponse struct {
	StatusCode int `json:"status"` // http response status code

	ErrorCode int    `json:"error_code,omitempty"` // application-specific error code
	ErrorMsg  string `json:"error"`                // user-facing error message

	StackErr error `json:"-"` // stack error for logging
}

func NewBadRequest(err error) *ErrResponse {
	return &ErrResponse{
		StackErr:   nil, // here we don't want to log err to StackErr.
		StatusCode: http.StatusBadRequest,
		ErrorMsg:   err.Error(),
	}
}
