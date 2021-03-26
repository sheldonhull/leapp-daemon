package service

import (
	"github.com/pkg/errors"
	"leapp_daemon/core/configuration"
)

func RotateAllSessionsCredentials() error {
	config, err := configuration.ReadConfiguration()
	if err != nil { return err }

	for i := range config.PlainAwsSessions {
		sess := config.PlainAwsSessions[i]

		err = sess.RotateCredentials(nil)
		if err != nil {
			return err
		}
	}

	for i := range config.FederatedAwsSessions {
		sess := config.FederatedAwsSessions[i]

		err = sess.RotateCredentials(nil)
		if err != nil {
			return err
		}
	}

	err = config.Update()
	if err != nil {
		return err
	}

	return nil
}

func RotateSessionCredentialsWithMfaToken(sessionId string, mfaToken string) error {
	found := false

	config, err := configuration.ReadConfiguration()
	if err != nil { return err }

	for i := range config.PlainAwsSessions {
		sess := config.PlainAwsSessions[i]
		if sess.Id == sessionId {
			err = sess.RotateCredentials(&mfaToken)
			if err != nil {
				return err
			}
			found = true
			break
		}
	}

	if !found {
		return errors.New("No session found with id " + sessionId)
	}

	err = config.Update()
	if err != nil {
		return err
	}

	return nil
}
