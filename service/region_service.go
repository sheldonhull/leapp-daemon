package service

import (
	"github.com/pkg/errors"
	"leapp_daemon/constant"
	"leapp_daemon/core/aws_client"
	"leapp_daemon/core/configuration"
)

type SessionAndRegion struct {
	SessionId string
	SessionType string
	AwsRegion string
}

func EditAwsSessionRegion(sessionId string, region string) (*SessionAndRegion, error) {
	// Check if the region is valid: we check it instantly as it is very simple
	isRegionValid := aws_client.IsRegionValid(region)
	if !isRegionValid {
		return nil, errors.New("Region " + region + " is not a valid AWS region or is not currently in the supported list")
	}

	// Get configuration
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return nil, err
	}

	// Find a valid Aws Session
	for _, plainSession := range config.PlainAwsSessions {
		if plainSession.Id == sessionId {
			plainSession.Account.Region = region
			err = configuration.UpdateConfiguration(config, false)
			if err != nil { return nil, err }

			return &SessionAndRegion{  SessionId: sessionId, SessionType: constant.SessionTypePlain, AwsRegion: region }, nil
		}
	}

	for _, federatedSession := range config.FederatedAwsSessions {
		if federatedSession.Id == sessionId {
			federatedSession.Account.Region = region
			err = configuration.UpdateConfiguration(config, false)
			if err != nil { return nil, err }

			return &SessionAndRegion{  SessionId: sessionId, SessionType: constant.SessionTypeFederated, AwsRegion: region }, nil
		}
	}

	return nil, errors.New("No valid AWS session found for editing region")
}
