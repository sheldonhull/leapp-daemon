// Package controller Trusted AWS Sessions API
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
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/interface/http/controller/dto/request_dto/trusted_aws_session_dto"
	"leapp_daemon/interface/http/controller/dto/response_dto"
	"leapp_daemon/use_case"
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

	err = use_case.CreateTrustedAwsSession(requestDto.ParentId, requestDto.AccountName, requestDto.AccountNumber, requestDto.RoleName, requestDto.Region)
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
	//   Responses:
	//     200: getTrustedAwsSessionResponse

	logging.SetContext(context)

	requestDto := trusted_aws_session_dto.GetTrustedAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	sess, err := use_case.GetTrustedAwsSession(requestDto.Id)
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

	err = use_case.UpdateTrustedAwsSession(
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

	requestDto := trusted_aws_session_dto.DeleteTrustedAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.DeleteTrustedAwsSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
