package g_suite_auth

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_errors"
)

type GSuiteAuthThirdStepRequestDto struct {
	IsMfaTokenRequested bool `json:"isMfaTokenRequested"`
	ResponseForm string `json:"responseForm"`
	SubmitURL string `json:"submitURL"`
	Token string `json:"token"`
}

func (requestDto *GSuiteAuthThirdStepRequestDto) Build(context *gin.Context) error {
	err := custom_errors.NewBadRequestError(context.ShouldBindJSON(requestDto))
	return err
}
