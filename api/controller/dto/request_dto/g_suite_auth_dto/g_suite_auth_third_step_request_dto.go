package g_suite_auth_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_error"
)

type GSuiteAuthThirdStepRequestDto struct {
	IsMfaTokenRequested bool `json:"isMfaTokenRequested"`
	ResponseForm string `json:"responseForm"`
	SubmitURL string `json:"submitURL"`
	Token string `json:"token"`
}

func (requestDto *GSuiteAuthThirdStepRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return custom_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
