package use_case

import (
  session2 "leapp_daemon/domain/session"
)

func CreatePlainAwsSession(name string, accountNumber string, region string, user string,
	awsAccessKeyId string, awsSecretAccessKey string, mfaDevice string, profile string) error {

  /*
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return err
	}

	err = session2.CreatePlainAwsSession(
		config,
		name,
		accountNumber,
		region,
		user,
		awsAccessKeyId,
		awsSecretAccessKey,
		mfaDevice,
		profile)
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

func GetPlainAwsSession(id string) (*session2.PlainAwsSession, error) {
	var sess *session2.PlainAwsSession

	/*
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return sess, err
	}

	sess, err = session2.GetPlainAwsSession(config, id)
	if err != nil {
		return sess, err
	}
	 */

	return sess, nil
}

func UpdatePlainAwsSession(sessionId string, name string, accountNumber string, region string, user string,
	awsAccessKeyId string, awsSecretAccessKey string, mfaDevice string, profile string) error {

  /*
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return err
	}

	err = session2.UpdatePlainAwsSession(config, sessionId, name, accountNumber, region, user, awsAccessKeyId, awsSecretAccessKey, mfaDevice, profile)
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

func DeletePlainAwsSession(sessionId string) error {
  /*
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return err
	}

	err = session2.DeletePlainAwsSession(config, sessionId)
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

func StartPlainAwsSession(sessionId string) error {
  /*
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return err
	}

	// Passing nil because, it will be the rotate method to check if we need the mfaToken or not
	err = session2.StartPlainAwsSession(config, sessionId, nil)
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

func StopPlainAwsSession(sessionId string) error {
  /*
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return err
	}

	// Passing nil because, it will be the rotate method to check if we need the mfaToken or not
	err = session2.StopPlainAwsSession(config, sessionId)
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
