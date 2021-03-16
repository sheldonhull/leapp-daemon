package g_suite_auth

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_error"
)

type GSuiteAuthFirstStepRequestDto struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (requestDto *GSuiteAuthFirstStepRequestDto) Build(context *gin.Context) error {
	err := custom_error.NewBadRequestError(context.ShouldBindJSON(requestDto))
	return err
}
