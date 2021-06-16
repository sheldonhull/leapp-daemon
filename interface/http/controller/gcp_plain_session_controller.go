package controller

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/interface/http/controller/dto/request_dto/gcp_plain_session_dto"
	"leapp_daemon/interface/http/controller/dto/response_dto"
	gcp_plain_session_dto2 "leapp_daemon/interface/http/controller/dto/response_dto/gcp_plain_session_dto"
	"net/http"
)

func (controller *EngineController) GetGcpOauthUrl(context *gin.Context) {
	// swagger:route GET /gcp/oauth/url gcpPlainSession getGcpOauthUrl
	// Get the GCP OAuth url
	//   Responses:
	//     200: GcpOauthUrlResponse

	logging.SetContext(context)

	actions := controller.Providers.GetGcpPlainSessionActions()

	url, err := actions.GetOAuthUrl()
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := gcp_plain_session_dto2.GcpOauthUrlResponse{
		Message: "success",
		Data:    url,
	}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) CreateGcpPlainSession(context *gin.Context) {
	// swagger:route POST /gcp/session/plain gcpPlainSession createGcpPlainSession
	// Create a new GCP Plain Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := gcp_plain_session_dto.GcpCreatePlainSessionRequest{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetGcpPlainSessionActions()

	err = actions.CreateSession(requestDto.Name, requestDto.AccountId, requestDto.ProjectName,
		requestDto.ProfileName, requestDto.OauthCode)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) GetGcpPlainSession(context *gin.Context) {
	// swagger:route GET /gcp/session/plain/{id} gcpPlainSession gcpGetPlainSession
	// Get a GCP Plain Session
	//   Responses:
	//     200: GcpGetPlainSessionResponse

	logging.SetContext(context)

	requestDto := gcp_plain_session_dto.GcpGetPlainSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetGcpPlainSessionActions()
	gcpPlainSession, err := actions.GetSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := gcp_plain_session_dto2.GcpGetPlainSessionResponse{
		Message: "success",
		Data:    gcpPlainSession,
	}

	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) StartGcpPlainSession(context *gin.Context) {
	// swagger:route GET /gcp/session/plain/{id}/start gcpPlainSession startGcpPlainSession
	// Start a GCP Plain Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := gcp_plain_session_dto.GcpStartPlainSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetGcpPlainSessionActions()
	err = actions.StartSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) StopGcpPlainSession(context *gin.Context) {
	// swagger:route GET /gcp/session/plain/{id}/stop gcpPlainSession stopGcpPlainSession
	// Stop a GCP Plain Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := gcp_plain_session_dto.GcpStopPlainSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetGcpPlainSessionActions()
	err = actions.StopSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) DeleteGcpPlainSession(context *gin.Context) {
	// swagger:route DELETE /gcp/session/plain/{id} gcpPlainSession deleteGcpPlainSession
	// Delete a GCP Plain Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := gcp_plain_session_dto.GcpDeletePlainSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetGcpPlainSessionActions()
	err = actions.DeleteSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) EditGcpPlainSession(context *gin.Context) {
	// swagger:route PUT /gcp/session/plain/{id} gcpPlainSession editGcpPlainSession
	// Edit a GCP Plain Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestUriDto := gcp_plain_session_dto.GcpEditPlainSessionUriRequest{}
	err := (&requestUriDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	requestDto := gcp_plain_session_dto.GcpEditPlainSessionRequest{}
	err = (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetGcpPlainSessionActions()
	err = actions.EditSession(requestUriDto.Id, requestDto.Name, requestDto.ProjectName, requestDto.ProfileName)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
