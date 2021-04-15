package response_dto

import "github.com/gin-gonic/gin"

// swagger:response messageResponse
type MessageOnlyResponseWrapper struct {
	// in: body
	Body MessageOnlyResponseDto
}

type MessageOnlyResponseDto struct {
	Message string
}

func (responseDto *MessageOnlyResponseDto) ToMap() gin.H {
	return gin.H{
		"message": responseDto.Message,
	}
}
