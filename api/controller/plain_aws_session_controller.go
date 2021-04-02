package controller

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/api/controller/dto/request_dto/plain_aws_session_dto"
	"leapp_daemon/api/controller/dto/response_dto"
	"leapp_daemon/core/service"
	"leapp_daemon/logging"
	"net/http"
)

func GetPlainAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := plain_aws_session_dto.GetPlainAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	sess, err := service.GetPlainAwsSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: *sess}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func CreatePlainAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := plain_aws_session_dto.CreatePlainAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = service.CreatePlainAwsSession(
		requestDto.Name,
		requestDto.AccountNumber,
		requestDto.Region,
		requestDto.User,
		requestDto.AwsAccessKeyId,
		requestDto.AwsSecretAccessKey,
		requestDto.MfaDevice,
		requestDto.ProfileName)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func EditPlainAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestUriDto := plain_aws_session_dto.EditPlainAwsSessionUriRequestDto{}
	err := (&requestUriDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	requestDto := plain_aws_session_dto.EditPlainAwsSessionRequestDto{}
	err = (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = service.UpdatePlainAwsSession(
		requestUriDto.Id,
		requestDto.Name,
		requestDto.AccountNumber,
		requestDto.Region,
		requestDto.User,
		requestDto.AwsAccessKeyId,
		requestDto.AwsSecretAccessKey,
		requestDto.MfaDevice,
		requestDto.ProfileName)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func DeletePlainAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := plain_aws_session_dto.DeletePlainAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = service.DeletePlainAwsSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func StartPlainAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := plain_aws_session_dto.StartPlainAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = service.StartPlainAwsSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func StopPlainAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := plain_aws_session_dto.StopPlainAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = service.StopPlainAwsSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}