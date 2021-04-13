package trusted_aws_session_dto

import (
  "github.com/gin-gonic/gin"
  http_error2 "leapp_daemon/infrastructure/http/http_error"
)

type EditTrustedAwsSessionUriRequestDto struct {
	Id string `uri:"id" binding:"required"`
}

type EditTrustedAwsSessionRequestDto struct {
	ParentId      string `json:"parentId"`
	AccountName   string `json:"accountName"`
	AccountNumber string `json:"accountNumber" binding:"numeric,len=12"`
	RoleName      string `json:"roleName"`
	Region        string `json:"region"`
}

func (requestDto *EditTrustedAwsSessionRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return http_error2.NewBadRequestError(err)
	} else {
		return nil
	}
}

func (requestDto *EditTrustedAwsSessionUriRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return http_error2.NewBadRequestError(err)
	} else {
		return nil
	}
}
