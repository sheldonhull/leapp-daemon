package request_dto

import (
	"github.com/gin-gonic/gin"
)

type IRequestDto interface {
	Build(*gin.Context) error
	// GetFieldByName(fieldName string) interface{}
	// ToMap() map[string]interface{}
	ToServiceRequest() interface{} // service_requests.IServiceRequest
}
