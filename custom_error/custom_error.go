package custom_error

import (
	"fmt"
	"net/http"
)

type BadRequestError struct{
	StatusCode int
	Message    error
}

func (err BadRequestError) Error() string {
	return fmt.Sprintf(`%+v`, err.Message)
}

type UnprocessableEntityError struct{
	StatusCode int
	Message    error
}

func (err UnprocessableEntityError) Error() string {
	return fmt.Sprintf(`%+v`, err.Message)
}

func NewBadRequestError(err error) error {
	if err == nil { return nil }
	return BadRequestError{ StatusCode: http.StatusBadRequest, Message: err }
}

func NewUnprocessableEntityError(err error) error {
	if err == nil { return nil }
	return UnprocessableEntityError{ StatusCode: http.StatusUnprocessableEntity, Message: err }
}
