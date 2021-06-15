package use_case

import (
  "encoding/json"
  "github.com/google/uuid"
  "leapp_daemon/domain/session"
  "leapp_daemon/infrastructure/aws/sts_client"
  "leapp_daemon/infrastructure/http/http_error"
  "strings"
  "time"
)

type PlainAwsSessionService struct {
  Keychain             Keychain
  NamedProfilesActions NamedProfilesActions
}

func (service *PlainAwsSessionService) Create(alias string, awsAccessKeyId string, awsSecretAccessKey string,
  mfaDevice string, region string, profileName string) error {

  namedProfile, err := service.NamedProfilesActions.GetOrCreateNamedProfile(profileName)
  if err != nil {
    return err
  }

  plainAwsAccount := session.PlainAwsAccount{
    MfaDevice:              mfaDevice,
    Region:                 region,
    NamedProfileId:         namedProfile.Id,
    SessionTokenExpiration: "",
  }

  // TODO: extract UUID generation logic
  uuidString := uuid.New().String()
  uuidString = strings.Replace(uuidString, "-", "", -1)

  sess := session.PlainAwsSession{
    Id:           uuidString,
    Alias:        alias,
    Status:       session.NotActive,
    StartTime:    "",
    LastStopTime: "",
    Account:      &plainAwsAccount,
  }

  err = session.GetPlainAwsSessionsFacade().AddPlainAwsSession(sess)
  if err != nil {
    return err
  }

  // TODO: use access keys repository instead of direct keychain abstraction
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

func (service *PlainAwsSessionService) GetPlainAwsSession(id string) (*session.PlainAwsSession, error) {
  var sess *session.PlainAwsSession
  sess, err := session.GetPlainAwsSessionsFacade().GetPlainAwsSessionById(id)
  return sess, err
}

func (service *PlainAwsSessionService) UpdatePlainAwsSession(sessionId string, name string, accountNumber string, region string, user string,
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

func (service *PlainAwsSessionService) StartPlainAwsSession(sessionId string) error {
  plainAwsSession, err := session.GetPlainAwsSessionsFacade().GetPlainAwsSessionById(sessionId)
  if err != nil {
    return err
  }

  doesSessionTokenExist, err := service.Keychain.DoesSecretExist(plainAwsSession.Id + "-plain-aws-session-session-token")
  if err != nil {
    return err
  }

  if doesSessionTokenExist {
    sessionTokenExpiration := plainAwsSession.Account.SessionTokenExpiration

    if sessionTokenExpiration != "" {
      currentTime := time.Now()
      sessionTokenExpirationTime, err := time.Parse(time.RFC3339, sessionTokenExpiration)
      if err != nil {
        return err
      }

      if currentTime.After(sessionTokenExpirationTime) {
        err = service.generateSessionToken(*plainAwsSession)
        if err != nil {
          return err
        }
      }
    } else {
      err = service.generateSessionToken(*plainAwsSession)
      if err != nil {
        return err
      }
    }
  } else {
    err = service.generateSessionToken(*plainAwsSession)
    if err != nil {
      return err
    }
  }

  err = session.GetPlainAwsSessionsFacade().SetPlainAwsSessionStatusToPending(sessionId)
  if err != nil {
    return err
  }

  err = session.GetPlainAwsSessionsFacade().SetPlainAwsSessionStatusToActive(sessionId)
  if err != nil {
    return err
  }

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

// TODO: encapsulate this logic inside a session token generation interface
func (service *PlainAwsSessionService) generateSessionToken(plainAwsSession session.PlainAwsSession) error {
  accessKeyIdSecretName := plainAwsSession.Id + "-plain-aws-session-access-key-id"

  accessKeyId, err := service.Keychain.GetSecret(accessKeyIdSecretName)
  if err != nil {
    return http_error.NewUnprocessableEntityError(err)
  }

  secretAccessKeySecretName := plainAwsSession.Id + "-plain-aws-session-secret-access-key"

  secretAccessKey, err := service.Keychain.GetSecret(secretAccessKeySecretName)
  if err != nil {
    return http_error.NewUnprocessableEntityError(err)
  }

  credentials, err := sts_client.GenerateAccessToken(plainAwsSession.Account.Region,
    plainAwsSession.Account.MfaDevice, nil, accessKeyId, secretAccessKey)
  if err != nil {
    return err
  }

  credentialsJson, err := json.Marshal(credentials)
  if err != nil {
    return err
  }

  err = service.Keychain.SetSecret(string(credentialsJson),
    plainAwsSession.Id+"-plain-aws-session-session-token")
  if err != nil {
    return err
  }

  err = session.GetPlainAwsSessionsFacade().SetPlainAwsSessionSessionTokenExpiration(plainAwsSession.Id, *credentials.Expiration)
  if err != nil {
    return err
  }

  return nil
}
