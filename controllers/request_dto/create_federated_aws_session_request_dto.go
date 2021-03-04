package request_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_errors"
)

type CreateFederatedAwsSessionRequestDto struct {
	Name string `json:"name" binding:"required"`
	AccountNumber string `json:"accountNumber" binding:"required"`
	RoleName string `json:"roleName" binding:"required"`
	RoleArn string `json:"roleArn" binding:"required"`
	IdpArn string `json:"idpArn" binding:"required"`
	Region string `json:"region" binding:"required"`
	SsoUrl string `json:"ssoUrl" binding:"required"`
}

func (requestDto *CreateFederatedAwsSessionRequestDto) Build(context *gin.Context) error {
	err := custom_errors.NewBadRequestError(context.ShouldBindJSON(requestDto))
	return err
}