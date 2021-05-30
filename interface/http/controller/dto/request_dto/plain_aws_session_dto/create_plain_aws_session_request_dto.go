package plain_aws_session_dto

import (
  "github.com/gin-gonic/gin"
  "leapp_daemon/infrastructure/http/http_error"
)

// swagger:parameters createPlainAwsSession
type CreatePlainAwsSessionRequestWrapper struct {
	// plain aws session create body
	// in:body
	Body CreatePlainAwsSessionRequest
}

type CreatePlainAwsSessionRequest struct {
	// the name which will be displayed
	// required: true
	Name string `json:"name" binding:"required"`

	// the account number of the aws account related to the role
	// required: true
	AccountNumber string `json:"accountNumber" binding:"required,numeric,len=12"`

	// the region on which the session will be initialized
	// required: true
	Region string `json:"region" binding:"required"`

	User               string `json:"user" binding:"required"`
	MfaDevice          string `json:"mfaDevice"`
	AwsAccessKeyId     string `json:"awsAccessKeyId" binding:"required"`
	AwsSecretAccessKey string `json:"awsSecretAccessKey" binding:"required"`
	ProfileName        string `json:"profileName"`
}

func (requestDto *CreatePlainAwsSessionRequest) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
