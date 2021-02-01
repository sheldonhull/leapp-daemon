package controllers

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/controllers/request_dto"
	"leapp_daemon/controllers/response_dto"
	"leapp_daemon/error_handling"
	"leapp_daemon/services"
	"net/http"
)

func CreateFederatedAccountController(context *gin.Context) {
	requestDto := request_dto.CreateFederatedAccountRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		error_handling.ErrorHandler.Handle(context, err)
		return
	}

	err2 := services.CreateFederatedAwsAccount(requestDto.Name, requestDto.AccountNumber,
		requestDto.RoleName, requestDto.RoleArn, requestDto.IdpArn, requestDto.Region, requestDto.SsoUrl)
	if err2 != nil {
		error_handling.ErrorHandler.Handle(context, err2)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
