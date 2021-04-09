package sts

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"leapp_daemon/core/aws/regional_endpoint"
)

func GetStaticCredentialsClient(accessKeyId string, secretAccessKey string, region *string) (*sts.STS, error) {
	endpoint, err := regional_endpoint.GetRegionalEndpoint(region)
	if err != nil {
		return nil, err
	}

	stsConfig := &aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKeyId, secretAccessKey, ""),
		Region: region,
		Endpoint: endpoint,
	}

	sess, err := session.NewSession(stsConfig)
	stsClient := sts.New(session.Must(sess, err))

	return stsClient, nil
}
