package trusted_aws_session_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_error"
)

type CreateTrusterAwsSessionRequestDto struct {
	AccountName string `json:"accountName" binding:"required"`
	AccountNumber string `json:"accountNumber" binding:"required"`
	RoleName string `json:"roleName" binding:"required"`
	Region string `json:"region"`
}

func (requestDto *CreateTrusterAwsSessionRequestDto) Build(context *gin.Context) error {
	err := custom_error.NewBadRequestError(context.ShouldBindJSON(requestDto))
	return err
}
