package plain_aws_session_dto

import (
  "github.com/gin-gonic/gin"
  "leapp_daemon/domain/session"
)

// swagger:response getPlainAwsSessionResponse
type GetPlainAwsSessionResponseWrapper struct {
  // in: body
  Body GetPlainAwsSessionResponse
}

type GetPlainAwsSessionResponse struct {
  Message string
  Data    session.PlainAwsSession
}

func (responseDto *GetPlainAwsSessionResponse) ToMap() gin.H {
  return gin.H{
    "message": responseDto.Message,
    "data": responseDto.Data,
  }
}
