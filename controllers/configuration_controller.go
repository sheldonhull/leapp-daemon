package controllers

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/controllers/response_dto"
	"leapp_daemon/error_handling"
	"leapp_daemon/services"
	"net/http"
)

func CreateConfigurationController(context *gin.Context) {
	err := services.CreateConfiguration()
	if err != nil {
		error_handling.ErrorHandler.Handle(context, err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func ReadConfigurationController(context *gin.Context) {
	configuration, err := services.ReadConfiguration()
	if err != nil {
		error_handling.ErrorHandler.Handle(context, err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: configuration}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
