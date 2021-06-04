package plain_alibaba_session_dto

import (
  "github.com/gin-gonic/gin"
  "leapp_daemon/domain/session"
)

// swagger:response getPlainAlibabaSessionResponse
type GetPlainAlibabaSessionResponseWrapper struct {
  // in: body
  Body GetPlainAlibabaSessionResponse
}

type GetPlainAlibabaSessionResponse struct {
  Message string
  Data    session.PlainAlibabaSession
}

func (responseDto *GetPlainAlibabaSessionResponse) ToMap() gin.H {
  return gin.H{
    "message": responseDto.Message,
    "data": responseDto.Data,
  }
}
