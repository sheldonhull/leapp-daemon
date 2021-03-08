package controllers

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/controllers/request_dto/federated_aws_session"
	"leapp_daemon/controllers/response_dto"
	"leapp_daemon/logging"
	"leapp_daemon/services/sessions"
	"net/http"
)

func CreateFederatedAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := federated_aws_session.CreateFederatedAwsSessionRequestDto{}
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
