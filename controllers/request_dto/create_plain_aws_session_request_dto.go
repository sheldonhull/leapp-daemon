package request_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_errors"
)

type CreatePlainAwsSessionRequestDto struct {
	Name string `json:"name" binding:"required"`
	AccountNumber string `json:"accountNumber" binding:"required"`
	Region string `json:"region" binding:"required"`
	User string `json:"user" binding:"required"`
	MfaDevice string `json:"mfaDevice"`
}

func (requestDto *CreatePlainAwsSessionRequestDto) Build(context *gin.Context) error {
	err := custom_errors.NewBadRequestError(context.ShouldBindJSON(requestDto))
	return err
}