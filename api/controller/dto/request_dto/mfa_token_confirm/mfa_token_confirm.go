package mfa_token_confirm

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_error"
)

type MfaTokenConfirmRequestDto struct {
	SessionId string `json:"sesionId" binding:"required"`
	MfaToken string `json:"mfaToken" binding:required`
}

func (requestDto *MfaTokenConfirmRequestDto) Build(context *gin.Context) error {
	err := custom_error.NewBadRequestError(context.ShouldBindJSON(requestDto))
	return err
}