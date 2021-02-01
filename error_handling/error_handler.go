package error_handling

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"leapp_daemon/controllers/utils"
	"leapp_daemon/logging"
	"net/http"
)

type errorHandler struct{}

func (*errorHandler) Handle(context *gin.Context, err error) {
	var code int

	switch err.(type) {
	case BadRequestError:
		code = http.StatusBadRequest
	default:
		code = http.StatusInternalServerError
	}

	errorMap := gin.H{ "statusCode": code, "error": err.Error(), "context": utils.NewContext(context) }

	logging.CtxLogger(context).
		WithFields(logrus.Fields{"statusCode": code}).
		Error(fmt.Sprintf("%s", err.Error()))
	context.JSON(code, errorMap)
}

var ErrorHandler = &errorHandler{}
