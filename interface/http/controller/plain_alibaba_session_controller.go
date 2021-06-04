package controller

import (
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/keychain"
	"leapp_daemon/infrastructure/logging"
	plain_alibaba_session_request_dto "leapp_daemon/interface/http/controller/dto/request_dto/plain_alibaba_session_dto"
	"leapp_daemon/interface/http/controller/dto/response_dto"
	plain_alibaba_session_response_dto "leapp_daemon/interface/http/controller/dto/response_dto/plain_alibaba_session_dto"
	"leapp_daemon/use_case"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreatePlainAlibabaSessionController(context *gin.Context) {
	// swagger:route POST /plain/alibaba/session/ plainAlibabaSession createPlainAlibabaSession
	// Create a new Plain Alibaba Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := plain_alibaba_session_request_dto.CreatePlainAlibabaSessionRequest{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	plainAlibabaSessionService := use_case.PlainAlibabaSessionService{
		Keychain: &keychain.Keychain{},
	}

	err = plainAlibabaSessionService.Create(requestDto.Name, requestDto.AlibabaAccessKeyId, requestDto.AlibabaSecretAccessKey,
		requestDto.Region, requestDto.ProfileName)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func GetPlainAlibabaSessionController(context *gin.Context) {
	// swagger:route GET /plain/alibaba/session/{id} plainAlibabaSession getPlainAlibabaSession
	// Get a Plain Alibaba Session
	//   Responses:
	//     200: GetPlainAlibabaSessionResponse

	logging.SetContext(context)

	requestDto := plain_alibaba_session_request_dto.GetPlainAlibabaSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	plainAlibabaSessionService := use_case.PlainAlibabaSessionService{
		Keychain: &keychain.Keychain{},
	}

	sess, err := plainAlibabaSessionService.GetPlainAlibabaSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := plain_alibaba_session_response_dto.GetPlainAlibabaSessionResponse{
		Message: "success",
		Data:    *sess,
	}

	context.JSON(http.StatusOK, responseDto.ToMap())
}

func UpdatePlainAlibabaSessionController(context *gin.Context) {
	// swagger:route PUT /plain/alibaba/session/{id} plainAlibabaSession updatePlainAlibabaSession
	// Edit a Plain Alibaba Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestUriDto := plain_alibaba_session_request_dto.UpdatePlainAlibabaSessionUriRequest{}
	err := (&requestUriDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	requestDto := plain_alibaba_session_request_dto.UpdatePlainAlibabaSessionRequest{}
	err = (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	plainAlibabaSessionService := use_case.PlainAlibabaSessionService{
		Keychain: &keychain.Keychain{},
	}

	err = plainAlibabaSessionService.UpdatePlainAlibabaSession(
		requestUriDto.Id,
		requestDto.Name,
		requestDto.Region,
		requestDto.User,
		requestDto.AlibabaAccessKeyId,
		requestDto.AlibabaSecretAccessKey,
		requestDto.ProfileName)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func DeletePlainAlibabaSessionController(context *gin.Context) {
	// swagger:route DELETE /plain/alibaba/session/{id} plainAlibabaSession deletePlainAlibabaSession
	// Delete a Plain Alibaba Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := plain_alibaba_session_request_dto.DeletePlainAlibabaSessionRequest{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = session.GetPlainAlibabaSessionsFacade().RemovePlainAlibabaSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func StartPlainAlibabaSessionController(context *gin.Context) {
	// swagger:route POST /plain/alibaba/session/{id}/start plainAlibabaSession startPlainAlibabaSession
	// Start a Plain Alibaba Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := plain_alibaba_session_request_dto.StartPlainAlibabaSessionRequest{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	plainAlibabaSessionService := use_case.PlainAlibabaSessionService{
		Keychain: &keychain.Keychain{},
	}

	err = plainAlibabaSessionService.StartPlainAlibabaSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func StopPlainAlibabaSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := plain_alibaba_session_request_dto.StopPlainAlibabaSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.StopPlainAlibabaSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
