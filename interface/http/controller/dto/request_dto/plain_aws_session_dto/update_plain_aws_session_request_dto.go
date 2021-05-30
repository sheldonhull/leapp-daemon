package plain_aws_session_dto

import (
  "github.com/gin-gonic/gin"
  http_error2 "leapp_daemon/infrastructure/http/http_error"
)

// swagger:parameters updatePlainAwsSession
type UpdatePlainAwsSessionUriRequestWrapper struct {
  // plain aws session update uri body
  // in:body
  Body UpdatePlainAwsSessionUriRequest
}

// swagger:parameters updatePlainAwsSession
type UpdatePlainAwsSessionRequestWrapper struct {
  // plain aws session update uri body
  // in:body
  Body UpdatePlainAwsSessionRequest
}

type UpdatePlainAwsSessionUriRequest struct {
	Id string `uri:"id" binding:"required"`
}

type UpdatePlainAwsSessionRequest struct {
	Name string `json:"name" binding:"required"`
	AccountNumber string `json:"accountNumber" binding:"required"`
	Region string `json:"region" binding:"required"`
	User string `json:"user" binding:"required"`
	MfaDevice string `json:"mfaDevice"`
	AwsAccessKeyId string `json:"awsAccessKeyId" binding:"required"`
	AwsSecretAccessKey string `json:"awsSecretAccessKey" binding:"required"`
	ProfileName string `json:"profileName"`
}

func (requestDto *UpdatePlainAwsSessionRequest) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return http_error2.NewBadRequestError(err)
	} else {
		return nil
	}
}

func (requestDto *UpdatePlainAwsSessionUriRequest) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return http_error2.NewBadRequestError(err)
	} else {
		return nil
	}
}
