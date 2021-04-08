package controller

import (
  "github.com/gin-gonic/gin"
  "leapp_daemon/api/controller/dto/request_dto/trusted_aws_session_dto"
  "leapp_daemon/api/controller/dto/response_dto"
  "leapp_daemon/core/service"
  "leapp_daemon/logging"
  "net/http"
)

func GetTrustedAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := trusted_aws_session_dto.GetTrustedAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	sess, err := service.GetTrustedAwsSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: *sess}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func CreateTrustedAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := trusted_aws_session_dto.CreateTrustedAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = service.CreateTrustedAwsSession(requestDto.ParentId, requestDto.AccountName, requestDto.AccountNumber, requestDto.RoleName, requestDto.Region)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

/*func EditTrusterAwsSessionController(context *gin.Context) {
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

	err = service.UpdateFederatedAwsSession(
		requestUriDto.Id,
		requestDto.Name,
		requestDto.AccountNumber,
		requestDto.RoleName,
		requestDto.RoleArn,
		requestDto.IdpArn,
		requestDto.Region,
		requestDto.SsoUrl)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func DeleteTrusterAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := federated_aws_session_dto.DeleteFederatedAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = service.DeleteFederatedAwsSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
*/
