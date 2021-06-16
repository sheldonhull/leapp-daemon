package aws_plain_session_dto

import (
	"github.com/gin-gonic/gin"
	http_error2 "leapp_daemon/infrastructure/http/http_error"
)

// swagger:parameters updateAwsPlainSession
type AwsUpdatePlainSessionUriRequestWrapper struct {
	// aws plain session update uri body
	// in:body
	Body AwsUpdatePlainSessionUriRequest
}

// swagger:parameters updateAwsPlainSession
type AwsUpdatePlainSessionRequestWrapper struct {
	// aws plain session update uri body
	// in:body
	Body AwsUpdatePlainSessionRequest
}

type AwsUpdatePlainSessionUriRequest struct {
	Id string `uri:"id" binding:"required"`
}

type AwsUpdatePlainSessionRequest struct {
	Name               string `json:"name" binding:"required"`
	AccountNumber      string `json:"accountNumber" binding:"required"`
	Region             string `json:"region" binding:"required"`
	User               string `json:"user" binding:"required"`
	MfaDevice          string `json:"mfaDevice"`
	AwsAccessKeyId     string `json:"awsAccessKeyId" binding:"required"`
	AwsSecretAccessKey string `json:"awsSecretAccessKey" binding:"required"`
	ProfileName        string `json:"profileName"`
}

func (requestDto *AwsUpdatePlainSessionRequest) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return http_error2.NewBadRequestError(err)
	} else {
		return nil
	}
}

func (requestDto *AwsUpdatePlainSessionUriRequest) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return http_error2.NewBadRequestError(err)
	} else {
		return nil
	}
}
