// Package classification of Trusted AWS Accounts API
//
// Documentation for Trusted AWS Accounts API
//
//  Schemes: http
//  Host: localhost
//  BasePath: /api/v1
//
//  Consumes:
//   - application/json
//
//  Produces:
//   - application/json
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

// swagger:response createTrustedAwsSessionResponse
type createTrustedAwsSessionResponse struct {
	Message string
	Data    session.TrustedAwsAccount
}

// CreateTrustedAwsSessionController returns a new AWS Trusted session
func CreateTrustedAwsSessionController(context *gin.Context) {

	// swagger:route POST /session/trusted createTrustedAwsSession
	//
	// Create a new AWS Trusted session
	//
	// Region is optional
	//
	//  Responses:
	//    200: createTrustedAwsSessionResponse

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
