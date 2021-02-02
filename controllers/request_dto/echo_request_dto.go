package request_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_errors"
)

type EchoRequestDto struct {
	Text string `uri:"Text" binding:"required"`
}

func (requestDto *EchoRequestDto) Build(context *gin.Context) error {
	err := custom_errors.NewBadRequestError(context.ShouldBindUri(requestDto))
	return err
}