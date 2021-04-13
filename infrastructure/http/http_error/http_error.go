package http_error

import (
  "net/http"
  "runtime"
)

const BaseNumberOfStackFramesToSkip = 3
const DefaultNumberOfStackFramesToSkip = 5

type stack []uintptr

func callers(skip int) *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

func GetStackTrace(skip int) []runtime.Frame {
	clrs := *callers(skip)

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

func NewCustomError(statusCode int, err error) CustomError {
  var stackTrace []runtime.Frame

  customErr, ok := err.(CustomError)
  if ok {
    stackTrace = customErr.StackTrace
  } else {
    stackTrace = GetStackTrace(DefaultNumberOfStackFramesToSkip)
  }

  return CustomError{ StatusCode: statusCode, Err: err, StackTrace: stackTrace }
}

func NewBadRequestError(err error) CustomError {
	return NewCustomError(http.StatusBadRequest, err)
}

func NewUnprocessableEntityError(err error) CustomError {
  return NewCustomError(http.StatusUnprocessableEntity, err)
}

func NewNotFoundError(err error) CustomError {
  return NewCustomError(http.StatusNotFound, err)
}

func NewInternalServerError(err error) CustomError {
  return NewCustomError(http.StatusInternalServerError, err)
}

func NewConflictError(err error) CustomError {
  return NewCustomError(http.StatusConflict, err)
}
