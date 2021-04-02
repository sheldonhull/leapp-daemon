package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"leapp_daemon/api/util"
	"leapp_daemon/custom_error"
	"leapp_daemon/logging"
	"net/http"
	"runtime"
)

type errorHandler struct{}

func (*errorHandler) Handle(context *gin.Context) {
	var code int
	var err error
	var errString string
	var stackTrace []runtime.Frame

	defer func() {
		panicErr := recover()

		if panicErr != nil || len(context.Errors) > 0 {
			if panicErr != nil {
				code = http.StatusInternalServerError
				errString = fmt.Sprintf("%s", panicErr)
				stackTrace = custom_error.GetStackTrace()
			} else if err != nil {
				switch err.(type) {
				case custom_error.CustomError:
					errString = err.Error()
					stackTrace = err.(custom_error.CustomError).StackTrace
				default:
					errString = fmt.Sprintf("%+v", err)
					stackTrace = make([]runtime.Frame, 0)
				}
			}

			errorMap := gin.H{"statusCode": code, "error": errString, "stackTrace": stackTrace, "context": util.NewContext(context)}

			logging.Entry().WithField("stackTrace", stackTrace).Error(errString)
			context.JSON(code, errorMap)
		}
	}()

	context.Next()

	if len(context.Errors) > 0 {
		err = context.Errors[0]
		err = err.(*gin.Error).Err
	} else {
		return
	}

	switch err.(type) {
	case custom_error.CustomError:
		code = err.(custom_error.CustomError).StatusCode
	default:
		code = http.StatusInternalServerError
	}
}

var ErrorHandler = &errorHandler{}