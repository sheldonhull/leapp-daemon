package use_case

import (
  "fmt"
  "gopkg.in/ini.v1"
  "leapp_daemon/domain/constant"
  "leapp_daemon/domain/named_profile"
  "leapp_daemon/domain/session"
  "leapp_daemon/infrastructure/http/http_error"
  "leapp_daemon/infrastructure/logging"
  "os"
)

type AwsCredentialsApplier struct {
  FileSystem FileSystem
  Keychain Keychain
}

func(awsCredentialsApplier *AwsCredentialsApplier) UpdatePlainAwsSessions(oldPlainAwsSessions []session.PlainAwsSession, newPlainAwsSessions []session.PlainAwsSession) error {
  for i, oldSess := range oldPlainAwsSessions {
    if i < len(newPlainAwsSessions) {
      newSess := newPlainAwsSessions[i]

      if oldSess.Status == session.NotActive && newSess.Status == session.Pending {

        homeDir, err := awsCredentialsApplier.FileSystem.GetHomeDir()
        if err != nil {
          return err
        }

        credentialsFilePath := homeDir + "/" + constant.CredentialsFilePath
        profileName := named_profile.GetNamedProfilesFacade().GetNamedProfileById(newSess.Account.NamedProfileId).Name
        region := newSess.Account.Region

        accessKeyId, secretAccessKey, err := awsCredentialsApplier.getAccessKeys(newSess.Id)
        if err != nil {
          return err
        }

        sessionToken, err := awsCredentialsApplier.getSessionToken(newSess.Id)
        if err != nil {
          return err
        }

        if awsCredentialsApplier.FileSystem.DoesFileExist(credentialsFilePath) {
          credentialsFile, err := ini.Load(credentialsFilePath)
          if err != nil {
            return err
          }

          section, err := credentialsFile.GetSection(profileName)
          if err != nil && err.Error() != fmt.Sprintf("section %q does not exist", profileName){
            return err
          }

          if section == nil {
            _, err = awsCredentialsApplier.createNamedProfileSection(credentialsFile, profileName, accessKeyId,
              secretAccessKey, sessionToken, region)
            if err != nil {
              return err
            }

            err = awsCredentialsApplier.appendToFile(credentialsFile, credentialsFilePath)
            if err != nil {
              return err
            }
          } else {
            credentialsFile.DeleteSection(profileName)

            _, err = awsCredentialsApplier.createNamedProfileSection(credentialsFile, profileName, accessKeyId, secretAccessKey,
              sessionToken, region)
            if err != nil {
              return err
            }

            err = awsCredentialsApplier.overwriteFile(credentialsFile, credentialsFilePath)
            if err != nil {
              return err
            }
          }
        } else {
          credentialsFile := ini.Empty()

          _, err = awsCredentialsApplier.createNamedProfileSection(credentialsFile, profileName, accessKeyId, secretAccessKey,
            sessionToken, region)
          if err != nil {
            return err
          }

          err = awsCredentialsApplier.overwriteFile(credentialsFile, credentialsFilePath)
          if err != nil {
            return err
          }
        }
      }
    }
  }

  return nil
}

func(awsCredentialsApplier *AwsCredentialsApplier) getAccessKeys(sessionId string) (accessKeyId string, secretAccessKey string, error error) {
  accessKeyId = ""
  secretAccessKey = ""

  accessKeyIdSecretName := sessionId + "-plain-aws-session-access-key-id"
  accessKeyId, err := awsCredentialsApplier.Keychain.GetSecret(accessKeyIdSecretName)
  if err != nil {
    return accessKeyId, secretAccessKey, http_error.NewUnprocessableEntityError(err)
  }

  secretAccessKeySecretName := sessionId + "-plain-aws-session-secret-access-key"
  secretAccessKey, err = awsCredentialsApplier.Keychain.GetSecret(secretAccessKeySecretName)
  if err != nil {
    return accessKeyId, secretAccessKey, http_error.NewUnprocessableEntityError(err)
  }

  return accessKeyId, secretAccessKey, nil
}

func(awsCredentialsApplier *AwsCredentialsApplier) getSessionToken(sessionId string) (sessionToken string, error error) {
  sessionToken = ""

  sessionTokenSecretName := sessionId + "-plain-aws-session-session-token"
  sessionToken, err := awsCredentialsApplier.Keychain.GetSecret(sessionTokenSecretName)
  if err != nil {
    return sessionToken, http_error.NewUnprocessableEntityError(err)
  }

  return sessionToken, nil
}

func(awsCredentialsApplier *AwsCredentialsApplier) createNamedProfileSection(credentialsFile *ini.File, profileName string, accessKeyId string,
  secretAccessKey string, sessionToken string, region string) (*ini.Section, error) {

  section, err := credentialsFile.NewSection(profileName)
  if err != nil {
    return nil, http_error.NewInternalServerError(err)
  }

  _, err = section.NewKey("aws_access_key_id", accessKeyId)
  if err != nil {
    return nil, http_error.NewInternalServerError(err)
  }

  _, err = section.NewKey("aws_secret_access_key", secretAccessKey)
  if err != nil {
    return nil, http_error.NewInternalServerError(err)
  }

  _, err = section.NewKey("aws_session_token", sessionToken)
  if err != nil {
    return nil, http_error.NewInternalServerError(err)
  }

  if region != "" {
    _, err = section.NewKey("region", region)
    if err != nil {
      return nil, http_error.NewInternalServerError(err)
    }
  }

  return section, nil
}

func(awsCredentialsApplier *AwsCredentialsApplier) appendToFile(file *ini.File, path string) error {
  f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0600)
  if err != nil {
    return http_error.NewNotFoundError(err)
  }

  _, err = file.WriteTo(f)
  if err != nil {
    return http_error.NewUnprocessableEntityError(err)
  }

  return nil
}

func(awsCredentialsApplier *AwsCredentialsApplier) overwriteFile(file *ini.File, path string) error {
  logging.Entry().Error("flag 3")

  f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
  if err != nil {
    return http_error.NewNotFoundError(err)
  }

  _, err = file.WriteTo(f)
  if err != nil {
    return http_error.NewUnprocessableEntityError(err)
  }

  return nil
}
