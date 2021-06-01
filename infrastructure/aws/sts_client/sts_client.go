package sts_client

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/credentials"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/sts"
  "leapp_daemon/domain/constant"
)

func GetStaticCredentialsClient(accessKeyId string, secretAccessKey string, region *string) (*sts.STS, error) {
	endpoint, err := GetRegionalEndpoint(region)
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

func GenerateAccessToken(region string, mfaDevice string, mfaToken *string, accessKeyId string,
  secretAccessKey string) (*sts.Credentials, error) {
  stsClient, err := GetStaticCredentialsClient(accessKeyId, secretAccessKey, &region)
  if err != nil {
    return nil, err
  }

  durationSeconds := constant.SessionTokenDurationInSeconds
  var getSessionTokenInput sts.GetSessionTokenInput

  if mfaToken == nil {
    getSessionTokenInput = sts.GetSessionTokenInput{
      DurationSeconds: &durationSeconds,
      SerialNumber:    nil,
      TokenCode:       nil,
    }
  } else {
    getSessionTokenInput = sts.GetSessionTokenInput{
      DurationSeconds: &durationSeconds,
      SerialNumber:    &mfaDevice,
      TokenCode:       mfaToken,
    }
  }

  getSessionTokenOutput, err := stsClient.GetSessionToken(&getSessionTokenInput)
  if err != nil {
    return nil, err
  }

  return getSessionTokenOutput.Credentials, nil
}
