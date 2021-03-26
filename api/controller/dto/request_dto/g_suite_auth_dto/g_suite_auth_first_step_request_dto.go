package g_suite_auth_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_error"
)

type GSuiteAuthFirstStepRequestDto struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (requestDto *GSuiteAuthFirstStepRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return custom_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
