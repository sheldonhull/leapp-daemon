package controller

import (
  "github.com/gin-gonic/gin"
  "leapp_daemon/infrastructure/keychain"
  "leapp_daemon/infrastructure/logging"
  gcp2 "leapp_daemon/interface/gcp"
  "leapp_daemon/interface/http/controller/dto/request_dto/gcp_plain_session_dto"
  "leapp_daemon/interface/http/controller/dto/response_dto"
  gcp_plain_session_dto2 "leapp_daemon/interface/http/controller/dto/response_dto/gcp_plain_session_dto"
  "leapp_daemon/use_case"
  "net/http"
  "sync"
)

func GetGcpOauthUrl(context *gin.Context) {
	// swagger:route GET /gcp/session/plain/oauth/url gcpPlainSession getGcpOauthUrl
	// Get the GCP OAuth url
	//   Responses:
	//     200: GcpOauthUrlResponse

	logging.SetContext(context)

	actions := getActions()

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

func CreateGcpPlainSession(context *gin.Context) {
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

	actions := getActions()

	err = actions.CreateSession(requestDto.Name, requestDto.AccountId, requestDto.ProjectName,
		requestDto.ProfileName, requestDto.OauthCode)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func GetGcpPlainSession(context *gin.Context) {
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

	actions := getActions()
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

func StartGcpPlainSession(context *gin.Context) {
  // swagger:route GET /gcp/session/plain/{id}/start gcpPlainSession startGcpPlainSession
  // Start a GCP Plain Session
  //   Responses:
  //     200: GcpGetPlainSessionResponse

  logging.SetContext(context)

  requestDto := gcp_plain_session_dto.GcpStartPlainSessionRequestDto{}
  err := (&requestDto).Build(context)
  if err != nil {
    _ = context.Error(err)
    return
  }

  actions := getActions()
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

var actionsSingleton *use_case.GcpPlainSessionActions
var actionsMutex sync.Mutex

func getActions() *use_case.GcpPlainSessionActions {
	actionsMutex.Lock()
	defer actionsMutex.Unlock()

	if actionsSingleton == nil {
		actionsSingleton = &use_case.GcpPlainSessionActions{
			Keychain: &keychain.Keychain{},
			GcpApi:   &gcp2.GcpApi{},
		}
	}
	return actionsSingleton
}
