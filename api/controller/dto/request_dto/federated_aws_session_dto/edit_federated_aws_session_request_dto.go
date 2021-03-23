package federated_aws_session_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_error"
)

type EditFederatedAwsSessionUriRequestDto struct {
	Id string `uri:"id" binding:"required"`
}

type EditFederatedAwsSessionRequestDto struct {
	Name string `json:"name" binding:"required"`
	AccountNumber string `json:"accountNumber" binding:"required"`
	RoleName string `json:"roleName" binding:"required"`
	RoleArn string `json:"roleArn" binding:"required"`
	IdpArn string `json:"idpArn" binding:"required"`
	Region string `json:"region" binding:"required"`
	SsoUrl string `json:"ssoUrl" binding:"required"`
}

func (requestDto *EditFederatedAwsSessionRequestDto) Build(context *gin.Context) error {
	err := custom_error.NewBadRequestError(context.ShouldBindJSON(requestDto))
	return err
}

func (requestDto *EditFederatedAwsSessionUriRequestDto) Build(context *gin.Context) error {
	err := custom_error.NewBadRequestError(context.ShouldBindUri(requestDto))
	return err
}