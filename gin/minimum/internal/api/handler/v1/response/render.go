package response

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RenderError(ctx *gin.Context, e *ErrResponse) {
	if e.StackErr != nil {
		zap.L().Error(e.StackErr.Error())
	}

	ctx.JSON(e.StatusCode, e)
}
