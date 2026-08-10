package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Error(ctx *gin.Context, err error) {
	if respErr, ok := IsResponseError(err); ok {
		ctx.JSON(respErr.Code, gin.H{"message": respErr.Msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "data": nil})
}

func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, gin.H{"message": "ok", "data": data})
}
