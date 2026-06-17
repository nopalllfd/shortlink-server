package response

import (
	"github.com/gin-gonic/gin"
	"github.com/nopalllfd/shortlink-server/internal/dto"
)

func Success(ctx *gin.Context, code int, message string, data any) {
	ctx.JSON(code, dto.Response{
		Message: message,
		Data:    data,
		Success: true,
	})
}

func Error(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, dto.Response{
		Message: message,
		Success: false,
	})
}
