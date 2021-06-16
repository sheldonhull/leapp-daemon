package controller

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/interface/http/controller/dto/request_dto/aws_federated_session_dto"
	"leapp_daemon/interface/http/controller/dto/response_dto"
	"leapp_daemon/use_case"
	"net/http"
)

func (controller *EngineController) GetAwsFederatedSession(context *gin.Context) {
	logging.SetContext(context)

	requestDto := aws_federated_session_dto.AwsGetFederatedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	sess, err := use_case.GetAwsFederatedSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: *sess}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) CreateAwsFederatedSession(context *gin.Context) {
	logging.SetContext(context)

	requestDto := aws_federated_session_dto.AwsCreateFederatedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.CreateAwsFederatedSession(requestDto.Name, requestDto.AccountNumber, requestDto.RoleName,
		requestDto.RoleArn, requestDto.IdpArn, requestDto.Region, requestDto.SsoUrl,
		requestDto.ProfileName)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) EditAwsFederatedSession(context *gin.Context) {
	logging.SetContext(context)

	requestUriDto := aws_federated_session_dto.AwsEditFederatedSessionUriRequestDto{}
	err := (&requestUriDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	requestDto := aws_federated_session_dto.AwsEditFederatedSessionRequestDto{}
	err = (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.UpdateAwsFederatedSession(
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

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) DeleteAwsFederatedSession(context *gin.Context) {
	logging.SetContext(context)

	requestDto := aws_federated_session_dto.AwsDeleteFederatedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.DeleteAwsFederatedSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) StartAwsFederatedSession(context *gin.Context) {
	logging.SetContext(context)

	requestDto := aws_federated_session_dto.AwsStartFederatedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.StartAwsFederatedSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) StopAwsFederatedSession(context *gin.Context) {
	logging.SetContext(context)

	requestDto := aws_federated_session_dto.AwsStopFederatedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.StopAwsFederatedSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
