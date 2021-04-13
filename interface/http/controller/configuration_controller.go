package controller

import (
  "github.com/gin-gonic/gin"
  "leapp_daemon/infrastructure/encryption"
  "leapp_daemon/infrastructure/file_system"
  "leapp_daemon/infrastructure/logging"
  "leapp_daemon/interface/http/controller/dto/response_dto"
  "leapp_daemon/use_case"
  "net/http"
)

func CreateConfigurationController(context *gin.Context) {
	logging.SetContext(context)

	service := use_case.ConfigurationService{
    FileSystem: &file_system.FileSystem{},
    Encryption: &encryption.Encryption{},
  }

	err := service.Create()
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func ReadConfigurationController(context *gin.Context) {
	logging.SetContext(context)

  service := use_case.ConfigurationService{
    FileSystem: &file_system.FileSystem{},
    Encryption: &encryption.Encryption{},
  }

	config, err := service.Read()
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: config}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
