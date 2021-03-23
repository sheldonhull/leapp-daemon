package federated_aws_session_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_error"
)

type DeleteFederatedAwsSessionRequestDto struct {
	Id string `uri:"id" binding:"required"`
}

func (requestDto *DeleteFederatedAwsSessionRequestDto) Build(context *gin.Context) error {
	err := custom_error.NewBadRequestError(context.ShouldBindUri(requestDto))
	return err
}