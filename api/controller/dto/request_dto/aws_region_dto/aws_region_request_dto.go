package aws_region_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_error"
)

type AwsRegionRequestDto struct {
	Region string `json:"region"`
	SessionId  string `json:"session_id"`
}

func (requestDto *AwsRegionRequestDto) Build(context *gin.Context) error {
	err := custom_error.NewBadRequestError(context.ShouldBindJSON(requestDto))
	return err
}