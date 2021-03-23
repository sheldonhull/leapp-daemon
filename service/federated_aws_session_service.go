package service

import (
	"leapp_daemon/core/configuration"
	"leapp_daemon/core/session"
)

func CreateFederatedAwsSession(name string, accountNumber string, roleName string, roleArn string,
	                           idpArn string, region string, ssoUrl string) error {

	config, err := configuration.ReadConfiguration()
	if err != nil { return err }

	err = session.CreateFederatedAwsSession(config, name, accountNumber, roleName, roleArn, idpArn, region, ssoUrl)
	if err != nil { return err }

	err = config.Update()
	if err != nil { return err }

	return nil
}

func GetFederatedAwsSession(id string) (*session.FederatedAwsSession, error) {
	var sess *session.FederatedAwsSession

	config, err := configuration.ReadConfiguration()
	if err != nil {
		return sess, err
	}

	sess, err = session.GetFederatedAwsSession(config, id)
	if err != nil {
		return sess, err
	}

	return sess, nil
}


func UpdateFederatedAwsSession(sessionId string, name string, accountNumber string, roleName string, roleArn string,
							   idpArn string, region string, ssoUrl string) error {

	config, err := configuration.ReadConfiguration()
	if err != nil {
		return err
	}

	err = session.UpdateFederatedAwsSession(config, sessionId, name, accountNumber, roleName, roleArn, idpArn, region, ssoUrl)
	if err != nil {
		return err
	}

	err = config.Update()
	if err != nil {
		return err
	}

	return nil
}

func DeleteFederatedAwsSession(sessionId string) error {
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return err
	}

	err = session.DeleteFederatedAwsSession(config, sessionId)
	if err != nil {
		return err
	}

	err = config.Update()
	if err != nil {
		return err
	}

	return nil
}