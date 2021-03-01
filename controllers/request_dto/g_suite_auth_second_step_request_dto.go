package request_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_errors"
)

type GSuiteAuthSecondStepRequestDto struct {
	Captcha string `json:"captcha"`
	CaptchaInputId string `json:"captchaInputId"`
	CaptchaUrl string `json:"captchaUrl"`
	CaptchaForm string `json:"captchaForm"`
	Password string `json:"password"`
	LoginForm string `json:"loginForm"`
	LoginUrl string `json:"loginUrl"`
}

func (requestDto *GSuiteAuthSecondStepRequestDto) Build(context *gin.Context) error {
	err := custom_errors.NewBadRequestError(context.ShouldBindJSON(requestDto))
	return err
}
