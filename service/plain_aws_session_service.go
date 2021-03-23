package service

import (
	"leapp_daemon/core/configuration"
	"leapp_daemon/core/session"
)

func CreatePlainAwsSession(name string, accountNumber string, region string, user string,
	awsAccessKeyId string, awsSecretAccessKey string, mfaDevice string) error {

	config, err := configuration.ReadConfiguration()
	if err != nil {
		return err
	}

	err = session.CreatePlainAwsSession(
		config,
		name,
		accountNumber,
		region,
		user,
		awsAccessKeyId,
		awsSecretAccessKey,
		mfaDevice)

	err = config.Update()
	if err != nil {
		return err
	}

	return nil
}
