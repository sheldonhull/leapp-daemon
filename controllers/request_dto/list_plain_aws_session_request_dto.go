package request_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_errors"
)

type ListPlainAwsSessionRequestDto struct {
	Query string `query:"query"`
}

func (requestDto *ListPlainAwsSessionRequestDto) Build(context *gin.Context) error {
	err := custom_errors.NewBadRequestError(context.ShouldBindJSON(requestDto))
	return err
}