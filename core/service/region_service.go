package service

import (
	"github.com/pkg/errors"
	"leapp_daemon/core/aws/aws_client"
	"leapp_daemon/core/configuration"
	"leapp_daemon/custom_error"
)

func EditAwsSessionRegion(sessionId string, region string) error {
	isRegionValid := aws_client.IsRegionValid(region)
	if !isRegionValid {
		return custom_error.NewUnprocessableEntityError(errors.New("Region " + region + " not valid"))
	}

	config, err := configuration.ReadConfiguration()
	if err != nil {
		return err
	}

	// Find a valid Aws Session
	for _, plainSession := range config.PlainAwsSessions {
		if plainSession.Id == sessionId {
			plainSession.Account.Region = region
			err = configuration.UpdateConfiguration(config, false)
			if err != nil { return err }
			return nil
		}
	}

	for _, federatedSession := range config.FederatedAwsSessions {
		if federatedSession.Id == sessionId {
			federatedSession.Account.Region = region
			err = configuration.UpdateConfiguration(config, false)
			if err != nil { return err }
			return nil
		}
	}

	return custom_error.NewUnprocessableEntityError(errors.New("No valid AWS session found for editing region"))
}
