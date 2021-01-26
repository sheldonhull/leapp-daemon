package controllers

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/controllers/request_dto"
	"leapp_daemon/controllers/response_dto"
	"leapp_daemon/services"
	"leapp_daemon/services/service_requests"
)

var HomeController = BaseController(&request_dto.HomeRequestDto{}, func(ctx *gin.Context, requestDto request_dto.IRequestDto) (response_dto.IResponseDto, error) {
	// serviceRequest := service_requests.HomeServiceRequest{}
	// err := (&serviceRequest).Build(requestDto.ToMap())

	// if err != nil {
	// 	return nil, err
	// }

	serviceResponse, err2 := services.HomeService(*requestDto.ToServiceRequest().(*service_requests.HomeServiceRequest))
	if err2 != nil { return nil, err2 }
	return serviceResponse.ToDto(), nil
})