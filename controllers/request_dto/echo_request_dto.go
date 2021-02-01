package request_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/error_handling"
)

type EchoRequestDto struct {
	Text string `uri:"text" binding:"required"`
}

func (requestDto *EchoRequestDto) Build(context *gin.Context) error {
	err := error_handling.NewBadRequestError(context.ShouldBindUri(requestDto))
	return err
}