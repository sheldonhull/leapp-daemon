package federated_aws_session_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_error"
)

type StartFederatedAwsSessionRequestDto struct {
	Id string `uri:"id" binding:"required"`
}

func (requestDto *StartFederatedAwsSessionRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return custom_error.NewBadRequestError(err)
	} else {
		return nil
	}
}