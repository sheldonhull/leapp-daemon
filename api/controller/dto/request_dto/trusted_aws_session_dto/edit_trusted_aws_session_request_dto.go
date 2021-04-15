package trusted_aws_session_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_error"
)

// swagger:parameters editTrustedAwsSession
type EditTrustedAwsSessionParamsWrapper struct {
	// This text will appear as description of your request body.
	// in: body
	Body EditTrustedAwsSessionRequestDto
}

type EditTrustedAwsSessionUriRequestDto struct {
	Id string `uri:"id" binding:"required"`
}

type EditTrustedAwsSessionRequestDto struct {
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

func (requestDto *EditTrustedAwsSessionRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return custom_error.NewBadRequestError(err)
	} else {
		return nil
	}
}

func (requestDto *EditTrustedAwsSessionUriRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindUri(requestDto)
	if err != nil {
		return custom_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
