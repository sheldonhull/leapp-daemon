package request_dto

import (
	"github.com/gin-gonic/gin"
)

type IRequestDto interface {
	Build(*gin.Context) error
	ToServiceRequest() interface{}
}
