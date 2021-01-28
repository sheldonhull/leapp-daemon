package error_handling

import (
	"fmt"
	"net/http"
)

type BadRequestError struct{
	StatusCode int
	Err error
}

func (err BadRequestError) Error() string {
	return fmt.Sprintf(`{ statuscode: %d, error: %+v }`, err.StatusCode, err.Err)
}

func NewBadRequestError(err error) error {
	if err == nil { return nil }
	return BadRequestError{ StatusCode: http.StatusBadRequest, Err: err }
}
