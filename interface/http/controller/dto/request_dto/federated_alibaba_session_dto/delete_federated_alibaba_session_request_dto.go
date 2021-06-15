package federated_alibaba_session_dto

import (
	http_error2 "leapp_daemon/infrastructure/http/http_error"

	"github.com/gin-gonic/gin"
)

type DeleteFederatedAlibabaSessionRequestDto struct {
	Id string `uri:"id" binding:"required"`
}

func (requestDto *DeleteFederatedAlibabaSessionRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return http_error2.NewBadRequestError(err)
	} else {
		return nil
	}
}
