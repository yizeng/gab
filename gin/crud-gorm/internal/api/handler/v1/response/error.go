package response

import (
	"fmt"
	"net/http"
)

type ErrResponse struct {
	StatusCode int `json:"status"` // http response status code

	ErrorCode int    `json:"error_code,omitempty"` // application-specific error code
	ErrorMsg  string `json:"error"`                // user-facing error message

	StackErr error `json:"-"` // stack error for logging
}

func NewBadRequest(msg string) *ErrResponse {
	return &ErrResponse{
		StatusCode: http.StatusBadRequest,
		ErrorMsg:   msg,
	}
}

func NewInvalidInput(fieldName string, fieldValue any) *ErrResponse {
	msg := fmt.Sprintf("invalid input field %v=%v", fieldName, fieldValue)

	return &ErrResponse{
		StatusCode: http.StatusBadRequest,
		ErrorMsg:   msg,
	}
}

func NewNotFound(resourceName, fieldName string, fieldValue any) *ErrResponse {
	msg := fmt.Sprintf("%v not found (%v=%v)", resourceName, fieldName, fieldValue)

	return &ErrResponse{
		StatusCode: http.StatusNotFound,
		ErrorMsg:   msg,
	}
}

func NewInternalServerError(err error) *ErrResponse {
	return &ErrResponse{
		StackErr:   err,
		StatusCode: http.StatusInternalServerError,
		ErrorMsg:   "something went wrong",
	}
}
