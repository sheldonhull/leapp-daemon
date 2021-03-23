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

	defer func() {
		panicErr := recover()

		if panicErr != nil || len(context.Errors) > 0 {
			if panicErr != nil {
				code = http.StatusInternalServerError
				errString = fmt.Sprintf("%s", panicErr)
			} else if len(context.Errors) > 0 {
				errString = err.Error()
			}

			errorMap := gin.H{"statusCode": code, "error": errString, "stackTrace": getStackTrace(), "context": util.NewContext(context)}

			logging.Entry().Error(errString)
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

type stack []uintptr

func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

func getStackTrace() []runtime.Frame {
	clrs := *callers()

	var frames []runtime.Frame
	callersFrames := runtime.CallersFrames(clrs)

	for {
		fr, more := callersFrames.Next()
		frames = append(frames, fr)
		if !more {
			break
		}
	}

	return frames
}