package controller

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/logging"
	plain_aws_session_request_dto "leapp_daemon/interface/http/controller/dto/request_dto/plain_aws_session_dto"
	"leapp_daemon/interface/http/controller/dto/response_dto"
	plain_aws_session_response_dto "leapp_daemon/interface/http/controller/dto/response_dto/plain_aws_session_dto"
	"leapp_daemon/use_case"
	"net/http"
)

func (env *EngineEnvironment) CreatePlainAwsSessionController(context *gin.Context) {
	// swagger:route POST /session/plain plainAwsSession createPlainAwsSession
	// Create a new Plain AWS Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := plain_aws_session_request_dto.CreatePlainAwsSessionRequest{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := env.Providers.GetAwsPlainSessionActions()

	err = actions.Create(requestDto.Name, requestDto.AwsAccessKeyId, requestDto.AwsSecretAccessKey,
		requestDto.MfaDevice, requestDto.Region, requestDto.ProfileName)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (env *EngineEnvironment) GetPlainAwsSessionController(context *gin.Context) {
	// swagger:route GET /session/plain/{id} plainAwsSession getPlainAwsSession
	// Get a Plain AWS Session
	//   Responses:
	//     200: GetPlainAwsSessionResponse

	logging.SetContext(context)

	requestDto := plain_aws_session_request_dto.GetPlainAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := env.Providers.GetAwsPlainSessionActions()

	sess, err := actions.GetPlainAwsSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := plain_aws_session_response_dto.GetPlainAwsSessionResponse{
		Message: "success",
		Data:    *sess,
	}

	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (env *EngineEnvironment) UpdatePlainAwsSessionController(context *gin.Context) {
	// swagger:route PUT /session/plain/{id} plainAwsSession updatePlainAwsSession
	// Edit a Plain AWS Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestUriDto := plain_aws_session_request_dto.UpdatePlainAwsSessionUriRequest{}
	err := (&requestUriDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	requestDto := plain_aws_session_request_dto.UpdatePlainAwsSessionRequest{}
	err = (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := env.Providers.GetAwsPlainSessionActions()

	err = actions.UpdatePlainAwsSession(
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

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (env *EngineEnvironment) DeletePlainAwsSessionController(context *gin.Context) {
	// swagger:route DELETE /session/plain/{id} plainAwsSession deletePlainAwsSession
	// Delete a Plain AWS Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := plain_aws_session_request_dto.DeletePlainAwsSessionRequest{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = session.GetPlainAwsSessionsFacade().RemovePlainAwsSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (env *EngineEnvironment) StartPlainAwsSessionController(context *gin.Context) {
	// swagger:route POST /session/plain/{id}/start plainAwsSession startPlainAwsSession
	// Start a Plain AWS Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := plain_aws_session_request_dto.StartPlainAwsSessionRequest{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := env.Providers.GetAwsPlainSessionActions()

	err = actions.StartPlainAwsSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (env *EngineEnvironment) StopPlainAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := plain_aws_session_request_dto.StopPlainAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.StopPlainAwsSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
