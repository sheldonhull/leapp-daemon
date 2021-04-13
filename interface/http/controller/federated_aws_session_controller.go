package controller

import (
  "github.com/gin-gonic/gin"
  "leapp_daemon/infrastructure/logging"
  "leapp_daemon/interfaces/http/controller/dto/request_dto/federated_aws_session_dto"
  "leapp_daemon/interfaces/http/controller/dto/response_dto"
  service2 "leapp_daemon/use_cases/service"
  "net/http"
)

// TODO: should pass DTOs to controllers, not *gin.Context

func GetFederatedAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := federated_aws_session_dto.GetFederatedAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	sess, err := service2.GetFederatedAwsSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: *sess}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func CreateFederatedAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := federated_aws_session_dto.CreateFederatedAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = service2.CreateFederatedAwsSession(requestDto.Name, requestDto.AccountNumber, requestDto.RoleName,
		                                    requestDto.RoleArn, requestDto.IdpArn, requestDto.Region, requestDto.SsoUrl,
		                                    requestDto.ProfileName)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func EditFederatedAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestUriDto := federated_aws_session_dto.EditFederatedAwsSessionUriRequestDto{}
	err := (&requestUriDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	requestDto := federated_aws_session_dto.EditFederatedAwsSessionRequestDto{}
	err = (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = service2.UpdateFederatedAwsSession(
		requestUriDto.Id,
		requestDto.Name,
		requestDto.AccountNumber,
		requestDto.RoleName,
		requestDto.RoleArn,
		requestDto.IdpArn,
		requestDto.Region,
		requestDto.SsoUrl,
		requestDto.ProfileName)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func DeleteFederatedAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := federated_aws_session_dto.DeleteFederatedAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = service2.DeleteFederatedAwsSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func StartFederatedAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := federated_aws_session_dto.StartFederatedAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = service2.StartFederatedAwsSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func StopFederatedAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := federated_aws_session_dto.StopFederatedAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = service2.StopFederatedAwsSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
