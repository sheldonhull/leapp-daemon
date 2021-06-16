package use_case

import (
	session2 "leapp_daemon/domain/session"
)

func CreateAwsFederatedSession(name string, accountNumber string, roleName string, roleArn string,
	idpArn string, region string, ssoUrl string, profile string) error {

	/*
		config, err := configuration.ReadConfiguration()
		if err != nil { return err }

		err = session2.CreateAwsFederatedSession(config, name, accountNumber, roleName, roleArn, idpArn, region, ssoUrl, profile)
		if err != nil { return err }

		err = config.Update()
		if err != nil { return err }
	*/

	return nil
}

func GetAwsFederatedSession(id string) (*session2.AwsFederatedSession, error) {
	var sess *session2.AwsFederatedSession

	/*
		config, err := configuration.ReadConfiguration()
		if err != nil {
			return sess, err
		}

		sess, err = session2.GetAwsFederatedSession(config, id)
		if err != nil {
			return sess, err
		}
	*/

	return sess, nil
}

func UpdateAwsFederatedSession(sessionId string, name string, accountNumber string, roleName string, roleArn string,
	idpArn string, region string, ssoUrl string, profile string) error {

	/*
		config, err := configuration.ReadConfiguration()
		if err != nil {
			return err
		}

		err = session2.UpdateAwsFederatedSession(config, sessionId, name, accountNumber, roleName, roleArn, idpArn, region, ssoUrl, profile)
		if err != nil {
			return err
		}

		err = config.Update()
		if err != nil {
			return err
		}
	*/

	return nil
}

func DeleteAwsFederatedSession(sessionId string) error {
	/*
		config, err := configuration.ReadConfiguration()
		if err != nil {
			return err
		}

		err = session2.DeleteAwsFederatedSession(config, sessionId)
		if err != nil {
			return err
		}

		err = config.Update()
		if err != nil {
			return err
		}
	*/

	return nil
}

func StartAwsFederatedSession(sessionId string) error {
	/*
		config, err := configuration.ReadConfiguration()
		if err != nil {
			return err
		}

		// Passing nil because, it will be the rotate method to check if we need the mfaToken or not
		err = session2.StartAwsFederatedSession(config, sessionId)
		if err != nil {
			return err
		}

		err = config.Update()
		if err != nil {
			return err
		}
	*/

	return nil
}

func StopAwsFederatedSession(sessionId string) error {
	/*
		config, err := configuration.ReadConfiguration()
		if err != nil {
			return err
		}

		// Passing nil because, it will be the rotate method to check if we need the mfaToken or not
		err = session2.StopAwsFederatedSession(config, sessionId)
		if err != nil {
			return err
		}

		err = config.Update()
		if err != nil {
			return err
		}

		// sess, err := session.GetAwsPlainSession(config, sessionId)
		err = session_token.RemoveFromIniFile("default")
		if err != nil {
			return err
		}
	*/

	return nil
}
