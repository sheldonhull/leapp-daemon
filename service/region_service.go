package service

import (
	"github.com/pkg/errors"
	"leapp_daemon/core/aws_client"
	"leapp_daemon/core/configuration"
)

func EditAwsSessionRegion(sessionId string, region string) error {
	// Check if the region is valid: we check it instantly as it is very simple
	isRegionValid := aws_client.IsRegionValid(region)
	if !isRegionValid {
		return errors.New("Region " + region + " not valid")
	}

	// Get configuration
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

	return errors.New("No valid AWS session found for editing region")
}
