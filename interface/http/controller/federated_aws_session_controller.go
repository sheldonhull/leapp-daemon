package controller

import (
  "github.com/gin-gonic/gin"
  "leapp_daemon/infrastructure/logging"
  "leapp_daemon/interface/http/controller/dto/request_dto/federated_aws_session_dto"
  "leapp_daemon/interface/http/controller/dto/response_dto"
  "leapp_daemon/use_case"
  "net/http"
)

func GetFederatedAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := federated_aws_session_dto.GetFederatedAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	sess, err := use_case.GetFederatedAwsSession(requestDto.Id)
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

	err = use_case.CreateFederatedAwsSession(requestDto.Name, requestDto.AccountNumber, requestDto.RoleName,
		                                    requestDto.RoleArn, requestDto.IdpArn, requestDto.Region, requestDto.SsoUrl,
		                                    requestDto.ProfileName)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
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

	err = use_case.UpdateFederatedAwsSession(
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

func DeleteFederatedAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := federated_aws_session_dto.DeleteFederatedAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.DeleteFederatedAwsSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
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

	err = use_case.StartFederatedAwsSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
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

	err = use_case.StopFederatedAwsSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
