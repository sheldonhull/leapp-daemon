package gcp_plain_session_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/infrastructure/http/http_error"
)

// swagger:parameters editGcpPlainSession
type GcpEditPlainSessionUriRequestWrapper struct {
	// gcp plain session edit uri
	// in:body
	Body GcpEditPlainSessionUriRequest
}

// swagger:parameters editGcpPlainSession
type GcpEditPlainSessionRequestWrapper struct {
	// gcp plain session edit body
	// in:body
	Body GcpEditPlainSessionRequest
}

type GcpEditPlainSessionUriRequest struct {
	// the Id of the session
	//required: true
	Id string `uri:"id" binding:"required"`
}

type GcpEditPlainSessionRequest struct {
	// the name which will be displayed
	// required: true
	Name string `json:"name" binding:"required"`

	// the name of the gcp project
	// required: true
	ProjectName string `json:"projectName" binding:"required"`
}

func (requestDto *GcpEditPlainSessionUriRequest) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}

func (requestDto *GcpEditPlainSessionRequest) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
