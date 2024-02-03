package response

import (
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
