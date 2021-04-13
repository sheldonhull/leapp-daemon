package usecase

import (
  session2 "leapp_daemon/domain/session"
)

func CreateFederatedAwsSession(name string, accountNumber string, roleName string, roleArn string,
	                           idpArn string, region string, ssoUrl string, profile string) error {

  /*
	config, err := configuration.ReadConfiguration()
	if err != nil { return err }

	err = session2.CreateFederatedAwsSession(config, name, accountNumber, roleName, roleArn, idpArn, region, ssoUrl, profile)
	if err != nil { return err }

	err = config.Update()
	if err != nil { return err }
   */

	return nil
}

func GetFederatedAwsSession(id string) (*session2.FederatedAwsSession, error) {
	var sess *session2.FederatedAwsSession

	/*
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return sess, err
	}

	sess, err = session2.GetFederatedAwsSession(config, id)
	if err != nil {
		return sess, err
	}
	 */

	return sess, nil
}


func UpdateFederatedAwsSession(sessionId string, name string, accountNumber string, roleName string, roleArn string,
							   idpArn string, region string, ssoUrl string, profile string) error {

  /*
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return err
	}

	err = session2.UpdateFederatedAwsSession(config, sessionId, name, accountNumber, roleName, roleArn, idpArn, region, ssoUrl, profile)
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

func DeleteFederatedAwsSession(sessionId string) error {
  /*
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return err
	}

	err = session2.DeleteFederatedAwsSession(config, sessionId)
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

func StartFederatedAwsSession(sessionId string) error {
  /*
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return err
	}

	// Passing nil because, it will be the rotate method to check if we need the mfaToken or not
	err = session2.StartFederatedAwsSession(config, sessionId)
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

func StopFederatedAwsSession(sessionId string) error {
  /*
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return err
	}

	// Passing nil because, it will be the rotate method to check if we need the mfaToken or not
	err = session2.StopFederatedAwsSession(config, sessionId)
	if err != nil {
		return err
	}

	err = config.Update()
	if err != nil {
		return err
	}

	// TODO: we need profileName branch here to change the profile
	// sess, err := session.GetPlainAwsSession(config, sessionId)
	err = session_token.RemoveFromIniFile("default")
	if err != nil {
		return err
	}
   */

	return nil
}
