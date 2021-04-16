package trusted_aws_session_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/infrastructure/http/http_error"
)

// swagger:parameters getTrustedAwsSession
type GetTrustedAwsSessionRequestDto struct {
	// the id of the trusted aws session
	// in: path
	// required: true
	Id string `uri:"id" binding:"required"`
}

func (requestDto *GetTrustedAwsSessionRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
