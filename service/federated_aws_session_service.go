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
