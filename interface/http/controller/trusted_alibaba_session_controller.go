package controller

import (
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/keychain"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/interface/http/controller/dto/request_dto/trusted_alibaba_session_dto"
	"leapp_daemon/interface/http/controller/dto/response_dto"
	"leapp_daemon/use_case"
	"net/http"

	"github.com/gin-gonic/gin"
)

// swagger:response getTrustedAlibabaSessionResponse
type getTrustedAlibabaSessionResponseWrapper struct {
	// in: body
	Body getTrustedAlibabaSessionResponse
}

type getTrustedAlibabaSessionResponse struct {
	Message string
	Data    session.TrustedAlibabaSession
}

func CreateTrustedAlibabaSessionController(context *gin.Context) {
	// swagger:route POST /session/trusted session-trusted-alibaba createTrustedAlibabaSession
	// Create a new trusted alibaba session
	//   Responses:
	//     200: messageResponse

	logging.SetContext(context)

	requestDto := trusted_alibaba_session_dto.CreateTrustedAlibabaSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	trustedAlibabaSessionService := use_case.TrustedAlibabaSessionService{
		Keychain: &keychain.Keychain{},
	}

	err = trustedAlibabaSessionService.Create(requestDto.ParentId, requestDto.AccountName, requestDto.AccountNumber, requestDto.RoleName, requestDto.Region)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func GetTrustedAlibabaSessionController(context *gin.Context) {
	// swagger:route GET /session/trusted/{id} session-trusted-alibaba getTrustedAlibabaSession
	// Get a Trusted AWS Session
	//   Responses:
	//     200: getTrustedAlibabaSessionResponse

	logging.SetContext(context)

	requestDto := trusted_alibaba_session_dto.GetTrustedAlibabaSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	trustedAlibabaSessionService := use_case.TrustedAlibabaSessionService{
		Keychain: &keychain.Keychain{},
	}

	sess, err := trustedAlibabaSessionService.Get(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: *sess}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func EditTrustedAlibabaSessionController(context *gin.Context) {
	// swagger:route PUT /session/trusted/{id} session-trusted-alibaba editTrustedAlibabaSession
	// Edit a Trusted AWS Session
	//   Responses:
	//     200: messageResponse

	logging.SetContext(context)

	requestUriDto := trusted_alibaba_session_dto.EditTrustedAlibabaSessionUriRequestDto{}
	err := (&requestUriDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	requestDto := trusted_alibaba_session_dto.EditTrustedAlibabaSessionRequestDto{}
	err = (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	trustedAlibabaSessionService := use_case.TrustedAlibabaSessionService{
		Keychain: &keychain.Keychain{},
	}

	err = trustedAlibabaSessionService.Update(
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

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func DeleteTrustedAlibabaSessionController(context *gin.Context) {
	// swagger:route DELETE /session/trusted/{id} session-trusted-alibaba deleteTrustedAlibabaSession
	// Delete a Trusted AWS Session
	//   Responses:
	//     200: messageResponse

	logging.SetContext(context)

	requestDto := trusted_alibaba_session_dto.DeleteTrustedAlibabaSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	trustedAlibabaSessionService := use_case.TrustedAlibabaSessionService{
		Keychain: &keychain.Keychain{},
	}

	err = trustedAlibabaSessionService.Delete(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func StartTrustedAlibabaSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := trusted_alibaba_session_dto.StartTrustedAlibabaSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	trustedAlibabaSessionService := use_case.TrustedAlibabaSessionService{
		Keychain: &keychain.Keychain{},
	}

	err = trustedAlibabaSessionService.Start(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func StopTrustedAlibabaSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := trusted_alibaba_session_dto.StopTrustedAlibabaSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	trustedAlibabaSessionService := use_case.TrustedAlibabaSessionService{
		Keychain: &keychain.Keychain{},
	}

	err = trustedAlibabaSessionService.Stop(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
