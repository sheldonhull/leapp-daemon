package controllers

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/controllers/response_dto"
	"leapp_daemon/logging"
	"leapp_daemon/services"
	"net/http"
)

func CreateConfigurationController(context *gin.Context) {
	logging.SetContext(context)

	err := services.CreateConfiguration()
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func ReadConfigurationController(context *gin.Context) {
	logging.SetContext(context)

	configuration, err := services.ReadConfiguration()
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: configuration}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
