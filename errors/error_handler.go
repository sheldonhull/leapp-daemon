package errors

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type errorHandler struct{}

func (*errorHandler) Handle(context *gin.Context, err error) {
	switch err.(type) {
	default:
		context.JSON(http.StatusInternalServerError, gin.H{ "error": err.Error() })
	}
}

var ErrorHandler = &errorHandler{}
