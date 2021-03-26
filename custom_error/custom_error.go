package custom_error

import (
	"fmt"
	"net/http"
)

type CustomError struct{
	StatusCode int
	Message    error
}

func (err CustomError) Error() string {
	return fmt.Sprintf(`%+v`, err.Message)
}

func NewBadRequestError(err error) error {
	if err == nil { return nil }
	return CustomError{ StatusCode: http.StatusBadRequest, Message: err }
}

func NewUnprocessableEntityError(err error) error {
	if err == nil { return nil }
	return CustomError{ StatusCode: http.StatusUnprocessableEntity, Message: err }
}

func NewNotFoundError(err error) error {
	if err == nil { return nil }
	return CustomError{ StatusCode: http.StatusNotFound, Message: err }
}
