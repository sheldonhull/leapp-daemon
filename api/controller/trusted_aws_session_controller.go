// Package controller Controller of Trusted AWS Sessions API
//
// Documentation for Trusted AWS Accounts API
//
//  Schemes: http
//  Host: localhost
//  BasePath: /api/v1/session
//
// swagger:meta
package controller

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/api/controller/dto/request_dto/federated_aws_session_dto"
	"leapp_daemon/api/controller/dto/request_dto/trusted_aws_session_dto"
	"leapp_daemon/api/controller/dto/response_dto"
	"leapp_daemon/core/service"
	"leapp_daemon/core/session"
	"leapp_daemon/logging"
	"net/http"
)

// swagger:response getTrustedAwsSessionResponse
type getTrustedAwsSessionResponseWrapper struct {
	// in: body
	Body getTrustedAwsSessionResponse
}

type getTrustedAwsSessionResponse struct {
	Message string
	Data    session.TrustedAwsSession
}

// swagger:parameters getTrustedAwsSession deleteTrustedAwsSession editTrustedAwsSession
type idTrustedAwsSessionParameterWrapper struct {
	// the id of the session
	// in: path
	// required: true
	Id string `json:"id"`
}

func CreateTrustedAwsSessionController(context *gin.Context) {
	// swagger:route POST /trusted session-trusted-aws createTrustedAwsSession
	// Create a new trusted aws session
	//   Responses:
	//     200: messageResponse

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

func GetTrustedAwsSessionController(context *gin.Context) {
	// swagger:route GET /trusted/{id} session-trusted-aws getTrustedAwsSession
	// Get a Trusted AWS Session
	//  Responses:
	//    200: getTrustedAwsSessionResponse

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

func EditTrustedAwsSessionController(context *gin.Context) {
	// swagger:route PUT /trusted/{id} session-trusted-aws editTrustedAwsSession
	// Edit a Trusted AWS Session
	//   Responses:
	//     200: messageResponse

	logging.SetContext(context)

	requestUriDto := trusted_aws_session_dto.EditTrustedAwsSessionUriRequestDto{}
	err := (&requestUriDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	requestDto := trusted_aws_session_dto.EditTrustedAwsSessionRequestDto{}
	err = (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = service.UpdateTrustedAwsSession(
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

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func DeleteTrustedAwsSessionController(context *gin.Context) {
	// swagger:route DELETE /trusted/{id} session-trusted-aws deleteTrustedAwsSession
	// Delete a Trusted AWS Session
	//   Responses:
	//     200: messageResponse

	logging.SetContext(context)

	requestDto := federated_aws_session_dto.DeleteFederatedAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = service.DeleteTrustedAwsSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
