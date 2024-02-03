package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Err struct {
	statusCode int // http response status code

	logFunc func() // a function used for logging if needed

	ErrorCode int    `json:"error_code,omitempty"` // application-specific error code
	ErrorMsg  string `json:"error"`                // user-facing error message
}

func RenderErr(ctx *gin.Context, e *Err) {
	if e.logFunc != nil {
		e.logFunc()
	}

	ctx.AbortWithStatusJSON(e.statusCode, e)
}

func ErrBadRequest(err error) *Err {
	return &Err{
		statusCode: http.StatusBadRequest,
		ErrorMsg:   err.Error(),
	}
}
