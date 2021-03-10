package request_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/shared/custom_error"
)

type EchoRequestDto struct {
	Text string `uri:"text" binding:"required"`
}

func (requestDto *EchoRequestDto) Build(context *gin.Context) error {
	err := custom_error.NewBadRequestError(context.ShouldBindUri(requestDto))
	return err
}