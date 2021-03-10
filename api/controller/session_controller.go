package controller

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/api/controller/dto/request_dto"
	"leapp_daemon/api/controller/dto/response_dto"
	"leapp_daemon/core/service/session"
	"leapp_daemon/shared/logging"
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

	sessionList, err2 := session.ListSessions(query, listType)
	if err2 != nil {
		_ = context.Error(err2)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: sessionList }
	context.JSON(http.StatusOK, responseDto.ToMap())
}