package request_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_errors"
)

type DeletePlainAwsSessionRequestDto struct {
	Id string `uri:"id" binding:"required"`
}

func (requestDto *DeletePlainAwsSessionRequestDto) Build(context *gin.Context) error {
	err := custom_errors.NewBadRequestError(context.ShouldBindUri(requestDto))
	return err
}