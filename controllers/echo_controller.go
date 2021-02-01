package controllers

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/controllers/request_dto"
	"leapp_daemon/controllers/response_dto"
	"leapp_daemon/error_handling"
	"leapp_daemon/services"
	"net/http"
)

func EchoController(context *gin.Context) {
	requestDto := request_dto.EchoRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		error_handling.ErrorHandler.Handle(context, err)
		return
	}

	serviceResponse, err2 := services.Echo(requestDto.Text)
	if err2 != nil {
		error_handling.ErrorHandler.Handle(context, err2)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: serviceResponse}
	context.JSON(http.StatusOK, responseDto.ToMap())
}