package use_case

import (
  "github.com/google/uuid"
  "leapp_daemon/domain/named_profile"
  "leapp_daemon/domain/session"
  "leapp_daemon/infrastructure/http/http_error"
  "strings"
)

type Keychain interface {
  DoesSecretExist(label string) (bool, error)
  GetSecret(label string) (string, error)
  SetSecret(secret string, label string) error
}

type PlainAwsSessionService struct {
  Keychain Keychain
}

func(service *PlainAwsSessionService) Create(alias string, awsAccessKeyId string, awsSecretAccessKey string,
  mfaDevice string, region string, profileName string) error {

  namedProfile := named_profile.GetNamedProfilesFacade().GetNamedProfileByName(profileName)

  if namedProfile == nil {
    // TODO: extract UUID generation logic
    uuidString := uuid.New().String()
    uuidString = strings.Replace(uuidString, "-", "", -1)

    namedProfile = &named_profile.NamedProfile{
      Id:   uuidString,
      Name: profileName,
    }

    err := named_profile.GetNamedProfilesFacade().AddNamedProfile(*namedProfile)
    if err != nil {
      return err
    }
  }

  plainAwsAccount := session.PlainAwsAccount{
    MfaDevice: mfaDevice,
    Region: region,
    NamedProfileId: namedProfile.Id,
    SessionTokenExpiration: "",
  }

  // TODO: extract UUID generation logic
  uuidString := uuid.New().String()
  uuidString = strings.Replace(uuidString, "-", "", -1)

  sess := session.PlainAwsSession{
    Id: uuidString,
    Alias: alias,
    Status: session.NotActive,
    StartTime: "",
    LastStopTime: "",
    Account: &plainAwsAccount,
  }

  err := session.GetPlainAwsSessionsFacade().AddPlainAwsSession(sess)
  if err != nil {
    return err
  }

  err = service.Keychain.SetSecret(awsAccessKeyId, sess.Id+"-plain-aws-session-access-key-id")
  if err != nil {
    return http_error.NewInternalServerError(err)
  }

  err = service.Keychain.SetSecret(awsSecretAccessKey, sess.Id+"-plain-aws-session-secret-access-key")
  if err != nil {
    return http_error.NewInternalServerError(err)
  }

  return nil
}

func CreatePlainAwsSession(name string, accountNumber string, region string, user string,
	awsAccessKeyId string, awsSecretAccessKey string, mfaDevice string, profile string) error {

  /*
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

func GetPlainAwsSession(id string) (*session.PlainAwsSession, error) {
	var sess *session.PlainAwsSession
	sess, err := session.GetPlainAwsSessionsFacade().GetPlainAwsSessionById(id)
  return sess, err
}

func UpdatePlainAwsSession(sessionId string, name string, accountNumber string, region string, user string,
	awsAccessKeyId string, awsSecretAccessKey string, mfaDevice string, profile string) error {

  /*
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return err
	}

	err = session.UpdatePlainAwsSession(config, sessionId, name, accountNumber, region, user, awsAccessKeyId, awsSecretAccessKey, mfaDevice, profile)
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

	err = session.DeletePlainAwsSession(config, sessionId)
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
	err = session.StartPlainAwsSession(config, sessionId, nil)
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
	err = session.StopPlainAwsSession(config, sessionId)
	if err != nil {
		return err
	}

	err = config.Update()
	if err != nil {
		return err
	}

  // sess, err := session.GetPlainAwsSession(config, sessionId)
	err = session_token.RemoveFromIniFile("default")
	if err != nil {
		return err
	}
   */

	return nil
}
