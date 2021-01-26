package response_dto

import "github.com/gin-gonic/gin"

type HomeResponseDto struct {
	Message string
	Data interface{}
}

func (responseDto *HomeResponseDto) ToMap() gin.H {
	return gin.H{
		"message": responseDto.Message,
		"data": responseDto.Data,
	}
}
