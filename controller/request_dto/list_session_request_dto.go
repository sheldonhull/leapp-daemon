package request_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_errors"
)

type ListSessionRequestDto struct {
	Query string `query:"query"`
	Type  string `query:"type"`
}

func (requestDto *ListSessionRequestDto) Build(context *gin.Context) error {
	err := custom_errors.NewBadRequestError(context.ShouldBindJSON(requestDto))
	return err
}