package controller

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/api/controller/dto/request_dto/plain_aws_session"
	"leapp_daemon/api/controller/dto/response_dto"
	"leapp_daemon/core/service/session"
	"leapp_daemon/shared/logging"
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

	session, err2 := session.GetPlainAwsSession(requestDto.Id)
	if err2 != nil {
		_ = context.Error(err2)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: session}
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

	err2 := session.CreatePlainAwsSession(
		requestDto.Name,
		requestDto.AccountNumber,
		requestDto.Region,
		requestDto.User,
		requestDto.AwsAccessKeyId,
		requestDto.AwsSecretAccessKey,
		requestDto.MfaDevice)

	if err2 != nil {
		_ = context.Error(err2)
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

	err2 := session.EditPlainAwsSession(
		requestUriDto.Id,
		requestDto.Name,
		requestDto.AccountNumber,
		requestDto.Region,
		requestDto.User,
		requestDto.AwsAccessKeyId,
		requestDto.AwsSecretAccessKey,
		requestDto.MfaDevice)

	if err2 != nil {
		_ = context.Error(err2)
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

	err2 := session.DeletePlainAwsSession(requestDto.Id)

	if err2 != nil {
		_ = context.Error(err2)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}