package controllers

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/controllers/request_dto"
	"leapp_daemon/controllers/response_dto"
	"leapp_daemon/error_handling"
	"leapp_daemon/services"
	"leapp_daemon/services/service_requests"
	"net/http"
)

func HomeController(context *gin.Context) {
	requestDto := request_dto.HomeRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		error_handling.ErrorHandler.Handle(context, err)
		return
	}

	serviceResponse, err2 := services.HomeService(service_requests.HomeServiceRequest{Name: requestDto.Name})
	if err2 != nil {
		error_handling.ErrorHandler.Handle(context, err2)
		return
	}

	responseDto := response_dto.HomeResponseDto{Message: "success", Data: serviceResponse.Data}
	context.JSON(http.StatusOK, responseDto.ToMap())
}