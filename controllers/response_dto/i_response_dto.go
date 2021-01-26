package response_dto

import "github.com/gin-gonic/gin"

type IResponseDto interface {
	ToMap() gin.H
}
