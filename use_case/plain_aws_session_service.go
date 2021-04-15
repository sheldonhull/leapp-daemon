package use_case

import (
  "fmt"
  "github.com/google/uuid"
  "leapp_daemon/domain/named_profile"
  "leapp_daemon/domain/session"
  "leapp_daemon/infrastructure/http/http_error"
  "strings"
)

type PlainAwsSessionService struct {
  PlainAwsSessionContainer session.PlainAwsSessionContainer
  NamedProfileContainer    named_profile.NamedProfileContainer
}

func(service *PlainAwsSessionService) Create(name string, accountNumber string, region string, user string,
  awsAccessKeyId string, awsSecretAccessKey string, mfaDevice string, profileName string) error {

  sessions, err := service.PlainAwsSessionContainer.GetAllPlainAwsSessions()
  if err != nil {
    return err
  }

  for _, sess := range sessions {
    account := sess.Account
    if account.AccountNumber == accountNumber && account.User == user {
      err := http_error.NewUnprocessableEntityError(fmt.Errorf("a session with the same account number and user is already present"))
      return err
    }
  }

  plainAwsAccount := session.PlainAwsAccount{
    AccountNumber: accountNumber,
    Name:          name,
    Region:        region,
    User:          user,
    AwsAccessKeyId: awsAccessKeyId,
    AwsSecretAccessKey: awsSecretAccessKey,
    MfaDevice:     mfaDevice,
  }

  // TODO: extract UUID generation logic
  uuidString := uuid.New().String()
  uuidString = strings.Replace(uuidString, "-", "", -1)

  if profileName == "" {
    profileName = "default"
  }

  doesNamedProfileExist := service.NamedProfileContainer.DoesNamedProfileExist(profileName)

  var namedProfile named_profile.NamedProfile

  if !doesNamedProfileExist {
    namedProfile = named_profile.NamedProfile{
      Id:   strings.Replace(uuid.New().String(), "-", "", -1),
      Name: profileName,
    }

    err = service.NamedProfileContainer.AddNamedProfile(namedProfile)
    if err != nil {
      return err
    }
  } else {
    namedProfile, err = service.NamedProfileContainer.FindNamedProfileByName(profileName)
    if err != nil {
      return err
    }
  }

  sess := session.PlainAwsSession{
    Id:        uuidString,
    Status:    session.NotActive,
    StartTime: "",
    Account:   &plainAwsAccount,
    Profile:   namedProfile.Id,
  }

  err = service.PlainAwsSessionContainer.AddPlainAwsSession(sess)
  if err != nil {
    return err
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

	/*
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return sess, err
	}

	sess, err = session.GetPlainAwsSession(config, id)
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
