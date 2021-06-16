package controller

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/logging"
	plain_aws_session_request_dto "leapp_daemon/interface/http/controller/dto/request_dto/aws_plain_session_dto"
	"leapp_daemon/interface/http/controller/dto/response_dto"
	plain_aws_session_response_dto "leapp_daemon/interface/http/controller/dto/response_dto/aws_plain_session_dto"
	"leapp_daemon/use_case"
	"net/http"
)

func (controller *EngineController) CreateAwsPlainSession(context *gin.Context) {
	// swagger:route POST /session/plain awsPlainSession createAwsPlainSession
	// Create a new AWS Plain Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := plain_aws_session_request_dto.AwsCreatePlainSessionRequest{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAwsPlainSessionActions()

	err = actions.Create(requestDto.Name, requestDto.AwsAccessKeyId, requestDto.AwsSecretAccessKey,
		requestDto.MfaDevice, requestDto.Region, requestDto.ProfileName)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) GetAwsPlainSession(context *gin.Context) {
	// swagger:route GET /session/plain/{id} awsPlainSession getAwsPlainSession
	// Get a AWS Plain Session
	//   Responses:
	//     200: AwsGetPlainSessionResponse

	logging.SetContext(context)

	requestDto := plain_aws_session_request_dto.AwsGetPlainSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAwsPlainSessionActions()

	sess, err := actions.GetAwsPlainSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := plain_aws_session_response_dto.AwsGetPlainSessionResponse{
		Message: "success",
		Data:    *sess,
	}

	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) UpdateAwsPlainSession(context *gin.Context) {
	// swagger:route PUT /session/plain/{id} awsPlainSession updateawsPlainSession
	// Edit a AWS Plain Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestUriDto := plain_aws_session_request_dto.AwsUpdatePlainSessionUriRequest{}
	err := (&requestUriDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	requestDto := plain_aws_session_request_dto.AwsUpdatePlainSessionRequest{}
	err = (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAwsPlainSessionActions()

	err = actions.UpdateAwsPlainSession(
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

func (controller *EngineController) DeleteAwsPlainSession(context *gin.Context) {
	// swagger:route DELETE /session/plain/{id} awsPlainSession deleteawsPlainSession
	// Delete a AWS Plain Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := plain_aws_session_request_dto.AwsDeletePlainSessionRequest{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = session.NewAwsPlainSessionsFacade().RemoveSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) StartAwsPlainSession(context *gin.Context) {
	// swagger:route POST /session/plain/{id}/start awsPlainSession startawsPlainSession
	// Start a AWS Plain Session
	//   Responses:
	//     200: MessageResponse

	logging.SetContext(context)

	requestDto := plain_aws_session_request_dto.AwsStartPlainSessionRequest{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	actions := controller.Providers.GetAwsPlainSessionActions()

	err = actions.StartAwsPlainSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) StopAwsPlainSession(context *gin.Context) {
	logging.SetContext(context)

	requestDto := plain_aws_session_request_dto.AwsStopPlainSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.StopAwsPlainSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
