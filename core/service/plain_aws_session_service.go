package service

import (
	"leapp_daemon/core/configuration"
)

func IsMfaRequiredForPlainAwsSession(id string) (bool, error) {
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return false, err
	}

	sess, err := configuration.GetById(config, id)
	if err != nil {
		return false, err
	}

	return sess.Account.MfaDevice != "", nil
}

func StartPlainAwsSession(id string, mfaToken *string) error {
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return err
	}

	sess, err := configuration.GetById(config, id)
	if err != nil {
		return err
	}

	err = sess.Rotate(config, nil)
	if err != nil {
		return err
	}

	return nil
}

func GetPlainAwsSession(id string) (*configuration.PlainAwsSession, error) {
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return nil, err
	}

	sess, err := configuration.GetById(config, id)
	if err != nil {
		return nil, err
	}

	return sess, nil
}
