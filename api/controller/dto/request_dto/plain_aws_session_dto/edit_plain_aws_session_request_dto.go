package plain_aws_session_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_error"
)

type EditPlainAwsSessionUriRequestDto struct {
	Id string `uri:"id" binding:"required"`
}

type EditPlainAwsSessionRequestDto struct {
	Name               string `json:"name" binding:"required"`
	AccountNumber      string `json:"accountNumber" binding:"required"`
	Region             string `json:"region" binding:"required,awsregion"`
	User               string `json:"user" binding:"required"`
	MfaDevice          string `json:"mfaDevice"`
	AwsAccessKeyId     string `json:"awsAccessKeyId" binding:"required"`
	AwsSecretAccessKey string `json:"awsSecretAccessKey" binding:"required"`
	ProfileName        string `json:"profileName"`
}

func (requestDto *EditPlainAwsSessionRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return custom_error.NewBadRequestError(err)
	} else {
		return nil
	}
}

func (requestDto *EditPlainAwsSessionUriRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return custom_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
