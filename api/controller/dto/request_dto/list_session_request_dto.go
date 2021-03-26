package request_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_error"
)

type ListSessionRequestDto struct {
	Query string `query:"query"`
	Type  string `query:"type"`
}

func (requestDto *ListSessionRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return custom_error.NewBadRequestError(err)
	} else {
		return nil
	}
}