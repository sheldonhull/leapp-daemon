package controller

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/api/controller/dto/response_dto"
	"leapp_daemon/core/configuration"
	"leapp_daemon/shared/logging"
	"net/http"
)

func CreateConfigurationController(context *gin.Context) {
	logging.SetContext(context)

	err := configuration.CreateConfiguration()
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func ReadConfigurationController(context *gin.Context) {
	logging.SetContext(context)

	configuration, err := configuration.ReadConfiguration()
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: configuration}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
