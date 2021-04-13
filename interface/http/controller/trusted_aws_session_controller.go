package controller

import (
  "github.com/gin-gonic/gin"
  logging2 "leapp_daemon/infrastructure/logging"
  federated_aws_session_dto2 "leapp_daemon/interfaces/http/controller/dto/request_dto/federated_aws_session_dto"
  trusted_aws_session_dto2 "leapp_daemon/interfaces/http/controller/dto/request_dto/trusted_aws_session_dto"
  response_dto2 "leapp_daemon/interfaces/http/controller/dto/response_dto"
  service2 "leapp_daemon/use_cases/service"
  "net/http"
)

func CreateTrustedAwsSessionController(context *gin.Context) {
	logging2.SetContext(context)

	requestDto := trusted_aws_session_dto2.CreateTrustedAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = service2.CreateTrustedAwsSession(requestDto.ParentId, requestDto.AccountName, requestDto.AccountNumber, requestDto.RoleName, requestDto.Region)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto2.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func GetTrustedAwsSessionController(context *gin.Context) {
	logging2.SetContext(context)

	requestDto := trusted_aws_session_dto2.GetTrustedAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	sess, err := service2.GetTrustedAwsSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto2.MessageAndDataResponseDto{Message: "success", Data: *sess}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func EditTrustedAwsSessionController(context *gin.Context) {
	logging2.SetContext(context)

	requestUriDto := trusted_aws_session_dto2.EditTrustedAwsSessionUriRequestDto{}
	err := (&requestUriDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	requestDto := trusted_aws_session_dto2.EditTrustedAwsSessionRequestDto{}
	err = (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = service2.UpdateTrustedAwsSession(
		requestUriDto.Id,
		requestDto.ParentId,
		requestDto.AccountName,
		requestDto.AccountNumber,
		requestDto.RoleName,
		requestDto.Region)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto2.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func DeleteTrustedAwsSessionController(context *gin.Context) {
	logging2.SetContext(context)

	requestDto := federated_aws_session_dto2.DeleteFederatedAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = service2.DeleteTrustedAwsSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto2.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
