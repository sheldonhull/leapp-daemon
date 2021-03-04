package controllers

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/controllers/request_dto"
	"leapp_daemon/controllers/response_dto"
	"leapp_daemon/services/sessions"
	"leapp_daemon/logging"
	"net/http"
)

func CreateFederatedAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := request_dto.CreateFederatedAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err2 := sessions.CreateFederatedAwsSession(requestDto.Name, requestDto.AccountNumber,
		requestDto.RoleName, requestDto.RoleArn, requestDto.IdpArn, requestDto.Region, requestDto.SsoUrl)
	if err2 != nil {
		_ = context.Error(err2)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
