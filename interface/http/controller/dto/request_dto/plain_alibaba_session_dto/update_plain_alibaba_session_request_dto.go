package plain_alibaba_session_dto

import (
	http_error2 "leapp_daemon/infrastructure/http/http_error"

	"github.com/gin-gonic/gin"
)

// swagger:parameters updatePlainAlibabaSession
type UpdatePlainAlibabaSessionUriRequestWrapper struct {
	// plain alibaba session update uri body
	// in:body
	Body UpdatePlainAlibabaSessionUriRequest
}

// swagger:parameters updatePlainAlibabaSession
type UpdatePlainAlibabaSessionRequestWrapper struct {
	// plain alibaba session update uri body
	// in:body
	Body UpdatePlainAlibabaSessionRequest
}

type UpdatePlainAlibabaSessionUriRequest struct {
	Id string `uri:"id" binding:"required"`
}

type UpdatePlainAlibabaSessionRequest struct {
	Name   string `json:"name" binding:"required"`
	Region string `json:"region" binding:"required"`
	//User string `json:"user" binding:"required"`
	AlibabaAccessKeyId     string `json:"alibabaAccessKeyId" binding:"required"`
	AlibabaSecretAccessKey string `json:"alibabaSecretAccessKey" binding:"required"`
	ProfileName            string `json:"profileName"`
}

func (requestDto *UpdatePlainAlibabaSessionRequest) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return http_error2.NewBadRequestError(err)
	} else {
		return nil
	}
}

func (requestDto *UpdatePlainAlibabaSessionUriRequest) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return http_error2.NewBadRequestError(err)
	} else {
		return nil
	}
}
