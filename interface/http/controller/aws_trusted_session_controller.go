package controller

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/interface/http/controller/dto/request_dto/aws_trusted_session_dto"
	"leapp_daemon/interface/http/controller/dto/response_dto"
	"leapp_daemon/use_case"
	"net/http"
)

// swagger:response getAwsTrustedSessionResponse
type getAwsTrustedSessionResponseWrapper struct {
	// in: body
	Body getAwsTrustedSessionResponse
}

type getAwsTrustedSessionResponse struct {
	Message string
	Data    session.AwsTrustedSession
}

func (controller *EngineController) CreateAwsTrustedSession(context *gin.Context) {
	// swagger:route POST /session/trusted session-trusted-aws createAwsTrustedSession
	// Create a new aws trusted session
	//   Responses:
	//     200: messageResponse

	logging.SetContext(context)

	requestDto := aws_trusted_session_dto.AwsCreateTrustedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.CreateAwsTrustedSession(requestDto.ParentId, requestDto.AccountName, requestDto.AccountNumber, requestDto.RoleName, requestDto.Region)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) GetAwsTrustedSession(context *gin.Context) {
	// swagger:route GET /session/trusted/{id} session-trusted-aws getAwsTrustedSession
	// Get a AWS Trusted Session
	//   Responses:
	//     200: getAwsTrustedSessionResponse

	logging.SetContext(context)

	requestDto := aws_trusted_session_dto.AwsGetTrustedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	sess, err := use_case.GetAwsTrustedSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: *sess}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func (controller *EngineController) EditAwsTrustedSession(context *gin.Context) {
	// swagger:route PUT /session/trusted/{id} session-trusted-aws editAwsTrustedSession
	// Edit a AWS Trusted Session
	//   Responses:
	//     200: messageResponse

	logging.SetContext(context)

	requestUriDto := aws_trusted_session_dto.AwsEditTrustedSessionUriRequestDto{}
	err := (&requestUriDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	requestDto := aws_trusted_session_dto.AwsEditTrustedSessionRequestDto{}
	err = (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.UpdateAwsTrustedSession(
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

func (controller *EngineController) DeleteAwsTrustedSession(context *gin.Context) {
	// swagger:route DELETE /session/trusted/{id} session-trusted-aws deleteAwsTrustedSession
	// Delete a AWS Trusted Session
	//   Responses:
	//     200: messageResponse

	logging.SetContext(context)

	requestDto := aws_trusted_session_dto.AwsDeleteTrustedSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.DeleteAwsTrustedSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageResponse{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
