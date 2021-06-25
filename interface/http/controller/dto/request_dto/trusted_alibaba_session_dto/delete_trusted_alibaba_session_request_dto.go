package trusted_alibaba_session_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/infrastructure/http/http_error"
)

// swagger:parameters deleteTrustedAlibabaSession
type DeleteTrustedAlibabaSessionRequestDto struct {
	// the id of the trusted alibaba session
	// in: path
	// required: true
	Id string `uri:"id" binding:"required"`
}

func (requestDto *DeleteTrustedAlibabaSessionRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
