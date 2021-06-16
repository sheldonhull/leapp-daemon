package aws_plain_session_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/domain/session"
)

// swagger:response getAwsPlainSessionResponse
type AwsGetPlainSessionResponseWrapper struct {
	// in: body
	Body AwsGetPlainSessionResponse
}

type AwsGetPlainSessionResponse struct {
	Message string
	Data    session.AwsPlainSession
}

func (responseDto *AwsGetPlainSessionResponse) ToMap() gin.H {
	return gin.H{
		"message": responseDto.Message,
		"data":    responseDto.Data,
	}
}
