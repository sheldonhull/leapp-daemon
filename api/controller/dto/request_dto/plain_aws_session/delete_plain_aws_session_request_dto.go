package plain_aws_session

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_error"
)

type DeletePlainAwsSessionRequestDto struct {
	Id string `uri:"id" binding:"required"`
}

func (requestDto *DeletePlainAwsSessionRequestDto) Build(context *gin.Context) error {
	err := custom_error.NewBadRequestError(context.ShouldBindUri(requestDto))
	return err
}