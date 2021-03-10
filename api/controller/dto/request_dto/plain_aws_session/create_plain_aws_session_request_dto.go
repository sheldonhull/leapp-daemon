package plain_aws_session

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/shared/custom_error"
)

type CreatePlainAwsSessionRequestDto struct {
	Name string `json:"name" binding:"required"`
	AccountNumber string `json:"accountNumber" binding:"required"`
	Region string `json:"region" binding:"required"`
	User string `json:"user" binding:"required"`
	MfaDevice string `json:"mfaDevice"`
}

func (requestDto *CreatePlainAwsSessionRequestDto) Build(context *gin.Context) error {
	err := custom_error.NewBadRequestError(context.ShouldBindJSON(requestDto))
	return err
}