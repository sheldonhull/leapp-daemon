package custom_error

import (
	"net/http"
	"runtime"
)

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

type CustomError struct{
	StatusCode int
	Err        error
	StackTrace []runtime.Frame
}

func (err CustomError) Error() string {
	return err.Err.Error()
}

func NewBadRequestError(err error) CustomError {
	return CustomError{ StatusCode: http.StatusBadRequest, Err: err, StackTrace: GetStackTrace() }
}

func NewUnprocessableEntityError(err error) CustomError {
	return CustomError{ StatusCode: http.StatusUnprocessableEntity, Err: err, StackTrace: GetStackTrace() }
}

func NewNotFoundError(err error) CustomError {
	return CustomError{ StatusCode: http.StatusNotFound, Err: err, StackTrace: GetStackTrace() }
}

func NewInternalServerError(err error) CustomError {
	return CustomError{ StatusCode: http.StatusInternalServerError, Err: err, StackTrace: GetStackTrace() }
}

func NewConflictError(err error) CustomError {
  return CustomError{ StatusCode: http.StatusConflict, Err: err, StackTrace: GetStackTrace() }
}
