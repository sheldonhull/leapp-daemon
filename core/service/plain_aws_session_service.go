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

func GetPlainAwsSession(id string) (*session.PlainAwsSession, error) {
	var sess *session.PlainAwsSession

	config, err := configuration.ReadConfiguration()
	if err != nil {
		return sess, err
	}

	sess, err = session.GetPlainAwsSession(config, id)
	if err != nil {
		return sess, err
	}

	return sess, nil
}

func UpdatePlainAwsSession(sessionId string, name string, accountNumber string, region string, user string,
	awsAccessKeyId string, awsSecretAccessKey string, mfaDevice string) error {

	config, err := configuration.ReadConfiguration()
	if err != nil {
		return err
	}

	err = session.UpdatePlainAwsSession(config, sessionId, name, accountNumber, region, user, awsAccessKeyId, awsSecretAccessKey, mfaDevice)
	if err != nil {
		return err
	}

	err = config.Update()
	if err != nil {
		return err
	}

	return nil
}

func DeletePlainAwsSession(sessionId string) error {
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return err
	}

	err = session.DeletePlainAwsSession(config, sessionId)
	if err != nil {
		return err
	}

	err = config.Update()
	if err != nil {
		return err
	}

	return nil
}
