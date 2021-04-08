package trusted_aws_session_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_error"
)

type EditTrustedAwsSessionUriRequestDto struct {
	Id string `uri:"id" binding:"required"`
}

type EditTrustedAwsSessionRequestDto struct {
  Id string `json:"id" binding:"required"`
  ParentId string `json:"parentId"`
	AccountName string `json:"accountName"`
	AccountNumber string `json:"accountNumber"`
	RoleName string `json:"roleName"`
	Region string `json:"region"`
}

func (requestDto *EditTrustedAwsSessionRequestDto) Build(context *gin.Context) error {
	err := custom_error.NewBadRequestError(context.ShouldBindJSON(requestDto))
	return err
}

func (requestDto *EditTrustedAwsSessionUriRequestDto) Build(context *gin.Context) error {
	err := custom_error.NewBadRequestError(context.ShouldBindUri(requestDto))
	return err
}
