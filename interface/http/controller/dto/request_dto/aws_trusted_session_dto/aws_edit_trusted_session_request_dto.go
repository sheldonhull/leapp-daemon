package aws_trusted_session_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/infrastructure/http/http_error"
)

// swagger:parameters editAwsTrustedSession
type AwsEditTrustedSessionParamsWrapper struct {
	// This text will appear as description of your request body.
	// in: body
	Body AwsEditTrustedSessionRequestDto
}

// swagger:parameters editAwsTrustedSession
type AwsEditTrustedSessionUriRequestDto struct {
	// the id of the aws trusted session
	// in: path
	// required: true
	Id string `uri:"id" binding:"required"`
}

type AwsEditTrustedSessionRequestDto struct {
	// the parent session id, can be an aws plain or federated session
	// it's generated with an uuid v4
	ParentId string `json:"parentId"`

	// the name which will be displayed
	AccountName string `json:"accountName"`

	// the account number of the aws account related to the role
	AccountNumber string `json:"accountNumber" binding:"numeric,len=12"`

	// the role name
	RoleName string `json:"roleName"`

	// the region on which the session will be initialized
	Region string `json:"region" binding:"awsregion"`
}

func (requestDto *AwsEditTrustedSessionRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}

func (requestDto *AwsEditTrustedSessionUriRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return http_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
