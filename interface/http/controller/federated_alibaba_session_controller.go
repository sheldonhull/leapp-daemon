package controller

import (
	"leapp_daemon/infrastructure/keychain"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/interface/http/controller/dto/request_dto/federated_alibaba_session_dto"
	"leapp_daemon/interface/http/controller/dto/response_dto"
	"leapp_daemon/use_case"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetFederatedAlibabaSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := federated_alibaba_session_dto.GetFederatedAlibabaSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	federatedAlibabaSessionService := use_case.FederatedAlibabaSessionService{
		Keychain: &keychain.Keychain{},
	}

	sess, err := federatedAlibabaSessionService.Get(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: *sess}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func CreateFederatedAlibabaSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := federated_alibaba_session_dto.CreateFederatedAlibabaSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	federatedAlibabaSessionService := use_case.FederatedAlibabaSessionService{
		Keychain: &keychain.Keychain{},
	}

	err = federatedAlibabaSessionService.Create(requestDto.Name, requestDto.AccountNumber, requestDto.RoleName,
		requestDto.RoleArn, requestDto.IdpArn, requestDto.Region, requestDto.SsoUrl,
		requestDto.ProfileName)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func EditFederatedAlibabaSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestUriDto := federated_alibaba_session_dto.EditFederatedAlibabaSessionUriRequestDto{}
	err := (&requestUriDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	requestDto := federated_alibaba_session_dto.EditFederatedAlibabaSessionRequestDto{}
	err = (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	federatedAlibabaSessionService := use_case.FederatedAlibabaSessionService{
		Keychain: &keychain.Keychain{},
	}

	err = federatedAlibabaSessionService.Update(
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

func DeleteFederatedAlibabaSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := federated_alibaba_session_dto.DeleteFederatedAlibabaSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	federatedAlibabaSessionService := use_case.FederatedAlibabaSessionService{
		Keychain: &keychain.Keychain{},
	}

	err = federatedAlibabaSessionService.Delete(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func StartFederatedAlibabaSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := federated_alibaba_session_dto.StartFederatedAlibabaSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	federatedAlibabaSessionService := use_case.FederatedAlibabaSessionService{
		Keychain: &keychain.Keychain{},
	}

	err = federatedAlibabaSessionService.Start(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func StopFederatedAlibabaSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := federated_alibaba_session_dto.StopFederatedAlibabaSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	federatedAlibabaSessionService := use_case.FederatedAlibabaSessionService{
		Keychain: &keychain.Keychain{},
	}

	err = federatedAlibabaSessionService.Stop(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
