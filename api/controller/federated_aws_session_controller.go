package controller

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/api/controller/dto/request_dto/federated_aws_session_dto"
	"leapp_daemon/api/controller/dto/response_dto"
	"leapp_daemon/logging"
	"leapp_daemon/service"
	"net/http"
)

func CreateFederatedAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := federated_aws_session_dto.CreateFederatedAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = service.CreateFederatedAwsSession(requestDto.Name, requestDto.AccountNumber, requestDto.RoleName,
		                                    requestDto.RoleArn, requestDto.IdpArn, requestDto.Region, requestDto.SsoUrl)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
