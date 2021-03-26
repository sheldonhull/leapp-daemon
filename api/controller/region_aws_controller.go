package controller

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/api/controller/dto/request_dto/aws_region_dto"
	"leapp_daemon/api/controller/dto/response_dto"
	"leapp_daemon/core/aws/aws_client"
	"leapp_daemon/logging"
	"leapp_daemon/service"
	"net/http"
)

func GetAwsRegionListController(context *gin.Context) {
	logging.SetContext(context)

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: aws_client.GetRegionList() }
	context.JSON(http.StatusOK, responseDto.ToMap())
}

func EditAwsRegionController(context *gin.Context) {
	logging.SetContext(context)

	requestDto := aws_region_dto.AwsRegionRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		_ = context.Error(err)
		return
	}

	err = service.EditAwsSessionRegion(requestDto.SessionId, requestDto.Region)
	if err != nil {
		_ = context.Error(err)
		return
	}

	responseDto := response_dto.MessageOnlyResponseDto{Message: "success"}
	context.JSON(http.StatusOK, responseDto.ToMap())
}