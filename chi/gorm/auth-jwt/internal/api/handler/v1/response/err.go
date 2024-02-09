package response

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type Err struct {
	statusCode int // http response status code

	logFunc func() // a function used for logging if needed

	ErrorCode int    `json:"error_code,omitempty"` // application-specific error code
	ErrorMsg  string `json:"error"`                // user-facing error message
}

func (e *Err) Render(w http.ResponseWriter, r *http.Request) error {
	if e.logFunc != nil {
		e.logFunc()
	}

	render.Status(r, e.statusCode)

	return nil
}

func ErrBadRequest(err error) *Err {
	return &Err{
		statusCode: http.StatusBadRequest,
		ErrorMsg:   err.Error(),
	}
}

func ErrInternalServerError(err error) *Err {
	return &Err{
		statusCode: http.StatusInternalServerError,
		logFunc: func() {
			zap.L().Error(err.Error())
		},
		ErrorMsg: "something went wrong",
	}
}

func ErrInvalidInput(fieldName string, fieldValue any) *Err {
	err := fmt.Errorf("invalid input field %v=%v", fieldName, fieldValue)

	return &Err{
		statusCode: http.StatusBadRequest,
		ErrorMsg:   err.Error(),
	}
}

func ErrNotFound(resourceName, fieldName string, fieldValue any) *Err {
	err := fmt.Errorf("%v not found (%v=%v)", resourceName, fieldName, fieldValue)

	return &Err{
		statusCode: http.StatusNotFound,
		ErrorMsg:   err.Error(),
	}
}

func ErrWrongCredentials(err error) *Err {
	return &Err{
		statusCode: http.StatusUnauthorized,
		logFunc: func() {
			zap.L().Debug("wrong credentials: " + err.Error())
		},
		ErrorMsg: "wrong credentials",
	}
}

func ErrJWTUnverified(err error) *Err {
	return &Err{
		statusCode: http.StatusUnauthorized,
		logFunc: func() {
			zap.L().Debug("unable to verify JWT: " + err.Error())
		},
		ErrorMsg: "please log in",
	}
}

func ErrPermissionDenied(err error) *Err {
	return &Err{
		statusCode: http.StatusForbidden,
		logFunc: func() {
			zap.L().Debug("permission denied: " + err.Error())
		},
		ErrorMsg: "permission denied",
	}
}
