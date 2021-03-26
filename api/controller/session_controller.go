package controller

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/api/controller/dto/request_dto"
	"leapp_daemon/api/controller/dto/request_dto/confirm_mfa_token_dto"
	"leapp_daemon/api/controller/dto/response_dto"
	"leapp_daemon/core/service"
	"leapp_daemon/logging"
	"net/http"
)

func ListSessionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := request_dto.ListSessionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	listType := requestDto.Type
	query := requestDto.Query

	sessionList, err2 := service.ListAllSessions(query, listType)
	if err2 != nil {
		_ = context.Error(err2)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: sessionList }
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func ConfirmMfaTokenController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := confirm_mfa_token_dto.MfaTokenConfirmRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = service.RotateSessionCredentialsWithMfaToken(requestDto.SessionId, requestDto.MfaToken)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: requestDto.SessionId }
	context.JSON(http.StatusOK, responseDto.ToMap())
}