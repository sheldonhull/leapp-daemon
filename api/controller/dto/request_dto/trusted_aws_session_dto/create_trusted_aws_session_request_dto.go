package trusted_aws_session_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_error"
)

type CreateTrustedAwsSessionRequestDto struct {
	ParentId      string `json:"parentId" binding:"required"`
	AccountName   string `json:"accountName" binding:"required"`
	AccountNumber string `json:"accountNumber" binding:"required,numeric,len=12"`
	RoleName      string `json:"roleName" binding:"required"`
	Region        string `json:"region"`
}

func (requestDto *CreateTrustedAwsSessionRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return custom_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
