package plain_aws_session

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/shared/custom_error"
)

type GetPlainAwsSessionRequestDto struct {
	Id string `uri:"id" binding:"required"`
}

func (requestDto *GetPlainAwsSessionRequestDto) Build(context *gin.Context) error {
	err := custom_error.NewBadRequestError(context.ShouldBindUri(requestDto))
	return err
}