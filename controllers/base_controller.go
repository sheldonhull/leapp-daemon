package controllers

import (
	"github.com/gin-gonic/gin"
	"leapp-daemon/controllers/request_dto"
	"leapp-daemon/controllers/response_dto"
	"leapp-daemon/errors"
	"net/http"
)

func BaseController(requestDto request_dto.IRequestDto, function func(*gin.Context, request_dto.IRequestDto) (response_dto.IResponseDto, error)) func(*gin.Context) {
	var decoratedFunction = func(context *gin.Context) {
		err := requestDto.Build(context)

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{ "error": err.Error() })
			return
		}

		responseDto, err2 := function(context, requestDto)

		if err2 != nil {
			errors.ErrorHandler.Handle(context, err2)
			return
		}

		context.JSON(http.StatusOK, responseDto.ToMap())
	}

	return decoratedFunction
}
