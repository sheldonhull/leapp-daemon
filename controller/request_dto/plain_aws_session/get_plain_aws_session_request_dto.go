package plain_aws_session

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_errors"
)

type GetPlainAwsSessionRequestDto struct {
	Id string `uri:"id" binding:"required"`
}

func (requestDto *GetPlainAwsSessionRequestDto) Build(context *gin.Context) error {
	err := custom_errors.NewBadRequestError(context.ShouldBindUri(requestDto))
	return err
}