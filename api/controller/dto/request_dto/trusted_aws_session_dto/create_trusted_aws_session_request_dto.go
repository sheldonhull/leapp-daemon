package trusted_aws_session_dto

import (
	"github.com/gin-gonic/gin"
	"leapp_daemon/custom_error"
)

// swagger:parameters createTrustedAwsSession
type CreateTrustedAwsSessionParamsWrapper struct {
	// This text will appear as description of your request body.
	// in:body
	Body CreateTrustedAwsSessionRequestDto
}

type CreateTrustedAwsSessionRequestDto struct {
	// the parent session id, can be an aws plain or federated session
	// it's generated with an uuid v4
	// required: true
	ParentId string `json:"parentId" binding:"required,uuid4"`

	// the name which will be displayed
	// required: true
	AccountName string `json:"accountName" binding:"required"`

	// the account number of the aws account related to the role
	// required: true
	AccountNumber string `json:"accountNumber" binding:"required,numeric,len=12"`

	// the role name
	// required: true
	RoleName string `json:"roleName" binding:"required"`

	// the region on which the session will be initialized
	Region string `json:"region" binding:"awsregion"`
}

func (requestDto *CreateTrustedAwsSessionRequestDto) Build(context *gin.Context) error {
	err := context.ShouldBindJSON(requestDto)
	if err != nil {
		return custom_error.NewBadRequestError(err)
	} else {
		return nil
	}
}
