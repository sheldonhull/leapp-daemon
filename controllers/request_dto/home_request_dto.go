package request_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/services/service_requests"
)

type HomeRequestDto struct {
	Name string `uri:"name" binding:"required"`
}

func (requestDto *HomeRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	return err
}

func (requestDto *HomeRequestDto) ToServiceRequest() interface{} {
	return &service_requests.HomeServiceRequest{Name: requestDto.Name}
}
