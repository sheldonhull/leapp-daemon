package trusted_alibaba_session_dto

import (
	"leapp_daemon/infrastructure/http/http_error"

	"github.com/gin-gonic/gin"
)

// swagger:parameters editTrustedAlibabaSession
type EditTrustedAlibabaSessionParamsWrapper struct {
	// This text will appear as description of your request body.
	// in: body
	Body EditTrustedAlibabaSessionRequestDto
}

// swagger:parameters editTrustedAlibabaSession
type EditTrustedAlibabaSessionUriRequestDto struct {
	// the id of the trusted alibaba session
	// in: path
	// required: true
	Id string `uri:"id" binding:"required"`
}

type EditTrustedAlibabaSessionRequestDto struct {
	// the parent session id, can be an alibaba plain or federated session
	// it's generated with an uuid v4
	ParentId string `json:"parentId"`

	// the name which will be displayed
	AccountName string `json:"accountName"`

	// the account number of the alibaba account related to the role
	AccountNumber string `json:"accountNumber" binding:"numeric,len=16"`

	// the role name
	RoleName string `json:"roleName"`

	// the region on which the session will be initialized
	Region string `json:"region"`
}

func (requestDto *EditTrustedAlibabaSessionRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}

func (requestDto *EditTrustedAlibabaSessionUriRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
