package federated_aws_session_dto

import (
  "github.com/gin-gonic/gin"
  http_error2 "leapp_daemon/infrastructure/http/http_error"
)

type StopFederatedAwsSessionRequestDto struct {
	Id string `uri:"id" binding:"required"`
}

func (requestDto *StopFederatedAwsSessionRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return http_error2.NewBadRequestError(err)
	} else {
		return nil
	}
}