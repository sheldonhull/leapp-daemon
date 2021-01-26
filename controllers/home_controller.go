package controllers

import (
	"github.com/gin-gonic/gin"
	"leapp-daemon/controllers/request_dto"
	"leapp-daemon/controllers/response_dto"
	"leapp-daemon/services"
	"leapp-daemon/services/service_requests"
)

var HomeController = BaseController(&request_dto.HomeRequestDto{}, func(ctx *gin.Context, requestDto request_dto.IRequestDto) (response_dto.IResponseDto, error) {
	serviceResponse, err2 := services.HomeService(*requestDto.ToServiceRequest().(*service_requests.HomeServiceRequest))
	if err2 != nil { return nil, err2 }
	return serviceResponse.ToDto(), nil
})