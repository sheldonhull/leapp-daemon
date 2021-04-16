// Package controller Plain AWS Sessions API
//
// Documentation for Plain AWS Accounts API
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
	"leapp_daemon/infrastructure/encryption"
	"leapp_daemon/infrastructure/file_system"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/interface/http/controller/dto/request_dto/plain_aws_session_dto"
	"leapp_daemon/interface/http/controller/dto/response_dto"
	"leapp_daemon/interface/repository"
	"leapp_daemon/use_case"
	"net/http"
)

// swagger:response getPlainAwsSessionResponse
type getPlainAwsSessionResponseWrapper struct {
	// in: body
	Body getPlainAwsSessionResponse
}

type getPlainAwsSessionResponse struct {
	Message string
	Data    session.PlainAwsSession
}

func CreatePlainAwsSessionController(context *gin.Context) {
	// swagger:route POST /plain session-plain-aws createPlainAwsSession
	// Create a new plain aws session
	//   Responses:
	//     200: messageResponse

	logging.SetContext(context)

	requestDto := plain_aws_session_dto.CreatePlainAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	configurationService := use_case.ConfigurationService{
		ConfigurationRepository: &repository.FileConfigurationRepository{
			FileSystem: &file_system.FileSystem{},
			Encryption: &encryption.Encryption{},
		},
	}

	configuration, err := configurationService.Get()
	if err != nil {
		_ = context.Error(err)
		return
	}

	plainAwsSessionService := use_case.PlainAwsSessionService{
		PlainAwsSessionContainer: &configuration,
		NamedProfileContainer:    &configuration,
	}

	err = plainAwsSessionService.Create(requestDto.Name, requestDto.AccountNumber, requestDto.Region, requestDto.User,
		requestDto.AwsAccessKeyId, requestDto.AwsSecretAccessKey, requestDto.MfaDevice, requestDto.ProfileName)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = configurationService.Update(configuration)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func GetPlainAwsSessionController(context *gin.Context) {
	// swagger:route GET /plain/{id} session-plain-aws getPlainAwsSession
	// Get a Plain AWS Session
	//   Responses:
	//     200: getPlainAwsSessionResponse

	logging.SetContext(context)

	requestDto := plain_aws_session_dto.GetPlainAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	sess, err := use_case.GetPlainAwsSession(requestDto.Id)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: *sess}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func UpdatePlainAwsSessionController(context *gin.Context) {
	// swagger:route PUT /plain/{id} session-plain-aws editPlainAwsSession
	// Edit a Plain AWS Session
	//   Responses:
	//     200: messageResponse

	logging.SetContext(context)

	requestUriDto := plain_aws_session_dto.UpdatePlainAwsSessionUriRequestDto{}
	err := (&requestUriDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	requestDto := plain_aws_session_dto.UpdatePlainAwsSessionRequestDto{}
	err = (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.UpdatePlainAwsSession(
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

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func DeletePlainAwsSessionController(context *gin.Context) {
	// swagger:route DELETE /plain/{id} session-plain-aws deletePlainAwsSession
	// Delete a Plain AWS Session
	//   Responses:
	//     200: messageResponse

	/*
		logging.SetContext(context)

		requestDto := plain_aws_session_dto.DeletePlainAwsSessionRequestDto{}
		err := (&requestDto).Build(context)
		if err != nil {
			_ = context.Error(err)
			return
		}

	  configurationService := use_case.ConfigurationService{
	    ConfigurationRepository: &repository.FileConfigurationRepository{
	      FileSystem: &file_system.FileSystem{},
	      Encryption: &encryption.Encryption{},
	    },
	  }

	  configuration, err := configurationService.Get()
	  if err != nil {
	    _ = context.Error(err)
	    return
	  }

	  configuration.

		responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
		context.JSON(http.StatusOK, responseDto.ToMap())
	*/
}

func StartPlainAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := plain_aws_session_dto.StartPlainAwsSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = use_case.StartPlainAwsSession(requestDto.Id)

	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func StopPlainAwsSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := plain_aws_session_dto.StopPlainAwsSessionRequestDto{}
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

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}
