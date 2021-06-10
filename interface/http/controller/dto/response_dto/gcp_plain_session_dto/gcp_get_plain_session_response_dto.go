package gcp_plain_session_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/domain/session"
)

// swagger:response gcpGetPlainSessionResponse
type GcpGetPlainSessionResponseWrapper struct {
	// in: body
	Body GcpGetPlainSessionResponse
}

type GcpGetPlainSessionResponse struct {
	Message string
	Data    session.GcpPlainSession
}

func (responseDto *GcpGetPlainSessionResponse) ToMap() gin.H {
	return gin.H{
		"message": responseDto.Message,
		"data":    responseDto.Data,
	}
}
