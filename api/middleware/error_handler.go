package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"leapp_daemon/api/util"
	"leapp_daemon/shared/custom_error"
	"leapp_daemon/shared/logging"
	"net/http"
)

type errorHandler struct{}

func (*errorHandler) Handle(context *gin.Context) {
	var code int
	var err error
	var errString string

	defer func() {
		panicErr := recover()

		if panicErr != nil || len(context.Errors) > 0 {
			if panicErr != nil {
				code = http.StatusInternalServerError
				errString = fmt.Sprintf("%s", panicErr)
			} else if len(context.Errors) > 0 {
				errString = err.Error()
			}

			errorMap := gin.H{"statusCode": code, "error": errString, "context": util.NewContext(context)}

			logging.Entry().
				WithFields(logrus.Fields{"statusCode": code}).
				Error(fmt.Sprintf("%s", err))

			context.JSON(code, errorMap)
		}
	}()

	context.Next()

	if len(context.Errors) > 0 {
		err = context.Errors[0]
	} else {
		return
	}

	switch err.(type) {
	case custom_error.BadRequestError:
		code = http.StatusBadRequest
	case custom_error.UnprocessableEntityError:
		code = http.StatusUnprocessableEntity
	default:
		code = http.StatusInternalServerError
	}
}

var ErrorHandler = &errorHandler{}
