package controller

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/api/controller/dto/request_dto"
	"leapp_daemon/api/controller/dto/request_dto/mfa_token_confirm"
	"leapp_daemon/api/controller/dto/response_dto"
	"leapp_daemon/logging"
	"leapp_daemon/service"
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

func MfaTokenConfirmController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := mfa_token_confirm.MfaTokenConfirmRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = service.ConfirmRotateSessionWithMfaToken(requestDto.SessionId, requestDto.MfaToken)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: requestDto.SessionId }
	context.JSON(http.StatusOK, responseDto.ToMap())
}