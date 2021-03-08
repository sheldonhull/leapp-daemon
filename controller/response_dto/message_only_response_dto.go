package response_dto

import "github.com/gin-gonic/gin"

type MessageOnlyResponseDto struct {
	Message string
}

func (responseDto *MessageOnlyResponseDto) ToMap() gin.H {
	return gin.H{
		"message": responseDto.Message,
	}
}
