package confirm_mfa_token_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_error"
)

type MfaTokenConfirmRequestDto struct {
	SessionId string `json:"sessionId" binding:"required"`
	MfaToken string `json:"mfaToken" binding:required`
}

func (requestDto *MfaTokenConfirmRequestDto) Build(context *gin.Context) error {
	err := custom_error.NewBadRequestError(context.ShouldBindJSON(requestDto))
	return err
}