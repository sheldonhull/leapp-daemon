package middleware

import (
  "fmt"
  "github.com/gin-gonic/gin"
  "leapp_daemon/infrastructure/http/context"
  "leapp_daemon/infrastructure/http/http_error"
  "leapp_daemon/infrastructure/logging"
  "net/http"
  "runtime"
)

type errorHandler struct{}

func (*errorHandler) Handle(ctx *gin.Context) {
	var code int
	var err error
	var errString string
	var stackTrace []runtime.Frame

	defer func() {
		panicErr := recover()

		if panicErr != nil || len(ctx.Errors) > 0 {
			if panicErr != nil {
				code = http.StatusInternalServerError
				errString = fmt.Sprintf("%s", panicErr)
				stackTrace = http_error.GetStackTrace(http_error.BaseNumberOfStackFramesToSkip)
			} else if err != nil {
				switch err.(type) {
				case http_error.CustomError:
					errString = err.Error()
					stackTrace = err.(http_error.CustomError).StackTrace
				default:
					errString = fmt.Sprintf("%+v", err)
					stackTrace = make([]runtime.Frame, 0)
				}
			}

			errorMap := gin.H{"statusCode": code, "error": errString, "stackTrace": stackTrace, "ctx": context.NewContext(ctx)}

			logging.Entry().WithField("stackTrace", stackTrace).Error(errString)
			ctx.JSON(code, errorMap)
		}
	}()

	ctx.Next()

	if len(ctx.Errors) > 0 {
		err = ctx.Errors[0]
		err = err.(*gin.Error).Err
	} else {
		return
	}

	switch err.(type) {
	case http_error.CustomError:
		code = err.(http_error.CustomError).StatusCode
	default:
		code = http.StatusInternalServerError
	}
}

var ErrorHandler = &errorHandler{}
