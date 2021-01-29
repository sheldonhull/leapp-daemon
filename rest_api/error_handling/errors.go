package error_handling

import (
	"fmt"
	"net/http"
)

type BadRequestError struct{
	StatusCode int
	Message    error
}

func (err BadRequestError) Error() string {
	return fmt.Sprintf(`[%d] %+v`, err.StatusCode, err.Message)
}

func NewBadRequestError(err error) error {
	if err == nil { return nil }
	return BadRequestError{ StatusCode: http.StatusBadRequest, Message: err }
}
