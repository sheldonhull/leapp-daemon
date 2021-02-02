package controllers

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/controllers/request_dto"
	"leapp_daemon/controllers/response_dto"
	"leapp_daemon/services"
	"log"
	"net/http"
)

func EchoController(context *gin.Context) {
	requestDto := request_dto.EchoRequestDto{}
	err := (&requestDto).Build(context)
	if err != nil {
		log.Println("ciao0")
		_ = context.Error(err)
		//custom_errors.ErrorHandler.Handle(context, err)
		return
	}

	serviceResponse, err2 := services.Echo(requestDto.Text)
	if err2 != nil {
		_ = context.Error(err)
		//custom_errors.ErrorHandler.Handle(context, err2)
		return
	}

	responseDto := response_dto.MessageAndDataResponseDto{Message: "success", Data: serviceResponse}
	context.JSON(http.StatusOK, responseDto.ToMap())
}