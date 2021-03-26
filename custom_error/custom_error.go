package custom_error

import (
	"fmt"
	"net/http"
	"runtime"
)

type CustomError struct{
	StatusCode int
	Message    error
	StackTrace []runtime.Frame
}

type stack []uintptr

func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

func GetStackTrace() []runtime.Frame {
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

func (err CustomError) Error() string {
	return fmt.Sprintf(`%+v`, err.Message)
}

func NewBadRequestError(err error) CustomError {
	return CustomError{ StatusCode: http.StatusBadRequest, Message: err, StackTrace: GetStackTrace()}
}

func NewUnprocessableEntityError(err error) CustomError {
	return CustomError{ StatusCode: http.StatusUnprocessableEntity, Message: err, StackTrace: GetStackTrace() }
}

func NewNotFoundError(err error) CustomError {
	return CustomError{ StatusCode: http.StatusNotFound, Message: err, StackTrace: GetStackTrace() }
}
