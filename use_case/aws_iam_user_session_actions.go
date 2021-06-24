package use_case

import (
	"encoding/json"
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/aws/sts_client"
	"leapp_daemon/infrastructure/http/http_error"
	"time"
)

type AwsIamUserSessionActions struct {
	Environment              Environment
	Keychain                 Keychain
	NamedProfilesActions     NamedProfilesActionsInterface
	AwsIamUserSessionsFacade AwsIamUserSessionsFacade
}

func (actions *AwsIamUserSessionActions) Create(name string, awsAccessKeyId string, awsSecretAccessKey string,
	mfaDevice string, region string, profileName string) error {

	namedProfile, err := actions.NamedProfilesActions.GetOrCreateNamedProfile(profileName)
	if err != nil {
		return err
	}

	sess := session.AwsIamUserSession{
		Id:                     actions.Environment.GenerateUuid(),
		Name:                   name,
		Status:                 session.NotActive,
		StartTime:              "",
		LastStopTime:           "",
		MfaDevice:              mfaDevice,
		Region:                 region,
		NamedProfileId:         namedProfile.Id,
		SessionTokenExpiration: "",
	}

	err = actions.AwsIamUserSessionsFacade.AddSession(sess)
	if err != nil {
		return err
	}

	// TODO: use access keys repository instead of direct keychain abstraction
	err = actions.Keychain.SetSecret(awsAccessKeyId, sess.Id+"-aws-iam-user-session-access-key-id")
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = actions.Keychain.SetSecret(awsSecretAccessKey, sess.Id+"-aws-iam-user-session-secret-access-key")
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (actions *AwsIamUserSessionActions) GetAwsIamUserSession(sessionId string) (session.AwsIamUserSession, error) {
	return actions.AwsIamUserSessionsFacade.GetSessionById(sessionId)
}

func (actions *AwsIamUserSessionActions) UpdateAwsIamUserSession(sessionId string, name string, accountNumber string, region string, user string,
	awsAccessKeyId string, awsSecretAccessKey string, mfaDevice string, profile string) error {

	/*
		config, err := configuration.ReadConfiguration()
		if err != nil {
			return err
		}

		err = session.UpdateAwsIamUserSession(config, sessionId, name, accountNumber, region, user, awsAccessKeyId, awsSecretAccessKey, mfaDevice, profile)
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

func DeleteAwsIamUserSession(sessionId string) error {
	/*
		config, err := configuration.ReadConfiguration()
		if err != nil {
			return err
		}

		err = session.DeleteAwsIamUserSession(config, sessionId)
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

func (actions *AwsIamUserSessionActions) StartAwsIamUserSession(sessionId string) error {
	awsIamUserSession, err := actions.AwsIamUserSessionsFacade.GetSessionById(sessionId)
	if err != nil {
		return err
	}

	doesSessionTokenExist, err := actions.Keychain.DoesSecretExist(awsIamUserSession.Id + "-aws-iam-user-session-session-token")
	if err != nil {
		return err
	}

	err = actions.AwsIamUserSessionsFacade.StartingSession(sessionId)
	if err != nil {
		return err
	}

	if doesSessionTokenExist {
		sessionTokenExpiration := awsIamUserSession.SessionTokenExpiration

		if sessionTokenExpiration != "" {
			currentTime := time.Now()
			sessionTokenExpirationTime, err := time.Parse(time.RFC3339, sessionTokenExpiration)
			if err != nil {
				return err
			}

			if currentTime.After(sessionTokenExpirationTime) {
				err = actions.generateSessionToken(awsIamUserSession)
				if err != nil {
					return err
				}
			}
		} else {
			err = actions.generateSessionToken(awsIamUserSession)
			if err != nil {
				return err
			}
		}
	} else {
		err = actions.generateSessionToken(awsIamUserSession)
		if err != nil {
			return err
		}
	}

	err = actions.AwsIamUserSessionsFacade.StartSession(sessionId, actions.Environment.GetTime())
	if err != nil {
		return err
	}

	return nil
}

func StopAwsIamUserSession(sessionId string) error {
	/*
			config, err := configuration.ReadConfiguration()
			if err != nil {
				return err
			}

			// Passing nil because, it will be the rotate method to check if we need the mfaToken or not
			err = session.StopAwsIamUserSession(config, sessionId)
			if err != nil {
				return err
			}

			err = config.Update()
			if err != nil {
				return err
			}

		  // sess, err := session.GetAwsIamUserSession(config, sessionId)
			err = session_token.RemoveFromIniFile("default")
			if err != nil {
				return err
			}
	*/

	return nil
}

// TODO: encapsulate this logic inside a session token generation interface
func (actions *AwsIamUserSessionActions) generateSessionToken(awsIamUserSession session.AwsIamUserSession) error {
	accessKeyIdSecretName := awsIamUserSession.Id + "-aws-iam-user-session-access-key-id"

	accessKeyId, err := actions.Keychain.GetSecret(accessKeyIdSecretName)
	if err != nil {
		return http_error.NewUnprocessableEntityError(err)
	}

	secretAccessKeySecretName := awsIamUserSession.Id + "-aws-iam-user-session-secret-access-key"

	secretAccessKey, err := actions.Keychain.GetSecret(secretAccessKeySecretName)
	if err != nil {
		return http_error.NewUnprocessableEntityError(err)
	}

	credentials, err := sts_client.GenerateAccessToken(awsIamUserSession.Region,
		awsIamUserSession.MfaDevice, nil, accessKeyId, secretAccessKey)
	if err != nil {
		return err
	}

	credentialsJson, err := json.Marshal(credentials)
	if err != nil {
		return err
	}

	err = actions.Keychain.SetSecret(string(credentialsJson),
		awsIamUserSession.Id+"-aws-iam-user-session-session-token")
	if err != nil {
		return err
	}

	err = actions.AwsIamUserSessionsFacade.SetSessionTokenExpiration(awsIamUserSession.Id, credentials.Expiration.Format(time.RFC3339))
	if err != nil {
		return err
	}

	return nil
}
