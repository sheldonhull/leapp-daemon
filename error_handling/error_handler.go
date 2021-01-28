package error_handling

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"leapp_daemon/controllers/utils"
	"log"
	"net/http"
)

type errorHandler struct{}

func (*errorHandler) Handle(context *gin.Context, err error) {
	log.Println(fmt.Sprintf("%+v", err))
	var code int

	switch err.(type) {
	case BadRequestError:
		code = http.StatusBadRequest
	default:
		code = http.StatusInternalServerError
	}

	context.JSON(code, gin.H{ "statuscode": code, "error": fmt.Sprintf("%+v", err), "context": utils.NewContext(context) })
}

var ErrorHandler = &errorHandler{}
