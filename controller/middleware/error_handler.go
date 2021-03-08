package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"leapp_daemon/controller/util"
	"leapp_daemon/custom_errors"
	"leapp_daemon/logging"
	"net/http"
)

type errorHandler struct{}

func (*errorHandler) Handle(context *gin.Context) {
	var code int
	var err error

	context.Next()

	if len(context.Errors) > 0 {
		err = context.Errors[0]
	} else {
		return
	}

	switch err.(type) {
	case custom_errors.BadRequestError:
		code = http.StatusBadRequest
	default:
		code = http.StatusInternalServerError
	}

	errorMap := gin.H{ "statusCode": code, "error": err.Error(), "context": util.NewContext(context) }

	logging.CtxEntry().
		WithFields(logrus.Fields{"statusCode": code}).
		Error(fmt.Sprintf("%s", err.Error()))

	context.JSON(code, errorMap)
}

var ErrorHandler = &errorHandler{}
