package controller

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/api/controller/dto/request_dto/plain_aws_session"
	"leapp_daemon/api/controller/dto/response_dto"
	"leapp_daemon/core/configuration"
	"leapp_daemon/logging"
	"leapp_daemon/service"
	"net/http"
)

func GetPlainAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := plain_aws_session.GetPlainAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	session, err := service.GetPlainAwsSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: *session}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func CreatePlainAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := plain_aws_session.CreatePlainAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = configuration.CreatePlainAwsSession(
		requestDto.Name,
		requestDto.AccountNumber,
		requestDto.Region,
		requestDto.User,
		requestDto.AwsAccessKeyId,
		requestDto.AwsSecretAccessKey,
		requestDto.MfaDevice)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func EditPlainAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestUriDto := plain_aws_session.EditPlainAwsSessionUriRequestDto{}
	err := (&requestUriDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	requestDto := plain_aws_session.EditPlainAwsSessionRequestDto{}
	err = (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = configuration.UpdatePlainAwsSession(
		requestUriDto.Id,
		requestDto.Name,
		requestDto.AccountNumber,
		requestDto.Region,
		requestDto.User,
		requestDto.AwsAccessKeyId,
		requestDto.AwsSecretAccessKey,
		requestDto.MfaDevice)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func DeletePlainAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := plain_aws_session.DeletePlainAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = configuration.DeletePlainAwsSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}