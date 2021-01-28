package controllers

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/rest_api/controllers/request_dto"
	"leapp_daemon/rest_api/controllers/response_dto"
	"leapp_daemon/rest_api/error_handling"
	"leapp_daemon/services"
	"net/http"
)

func HomeController(context *gin.Context) {
	requestDto := request_dto.HomeRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		error_handling.ErrorHandler.Handle(context, err)
		return
	}

	serviceResponse, err2 := services.Home(requestDto.Name)
	if err2 != nil {
		error_handling.ErrorHandler.Handle(context, err2)
		return
	}

	responseDto := response_dto.HomeResponseDto{Message: "success", Data: serviceResponse}
	context.JSON(http.StatusOK, responseDto.ToMap())
}