package use_case

import (
	"encoding/json"
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/aws/sts_client"
	"leapp_daemon/infrastructure/http/http_error"
	"time"
)

type AwsPlainSessionActions struct {
	Environment            Environment
	Keychain               Keychain
	NamedProfilesActions   NamedProfilesActionsInterface
	AwsPlainSessionsFacade AwsPlainSessionsFacade
}

func (actions *AwsPlainSessionActions) Create(alias string, awsAccessKeyId string, awsSecretAccessKey string,
	mfaDevice string, region string, profileName string) error {

	namedProfile, err := actions.NamedProfilesActions.GetOrCreateNamedProfile(profileName)
	if err != nil {
		return err
	}

	awsPlainAccount := session.AwsPlainAccount{
		MfaDevice:              mfaDevice,
		Region:                 region,
		NamedProfileId:         namedProfile.Id,
		SessionTokenExpiration: "",
	}

	sess := session.AwsPlainSession{
		Id:           actions.Environment.GenerateUuid(),
		Alias:        alias,
		Status:       session.NotActive,
		StartTime:    "",
		LastStopTime: "",
		Account:      &awsPlainAccount,
	}

	err = actions.AwsPlainSessionsFacade.AddSession(sess)
	if err != nil {
		return err
	}

	// TODO: use access keys repository instead of direct keychain abstraction
	err = actions.Keychain.SetSecret(awsAccessKeyId, sess.Id+"-plain-aws-session-access-key-id")
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = actions.Keychain.SetSecret(awsSecretAccessKey, sess.Id+"-plain-aws-session-secret-access-key")
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (actions *AwsPlainSessionActions) GetAwsPlainSession(id string) (*session.AwsPlainSession, error) {
	var sess *session.AwsPlainSession
	sess, err := actions.AwsPlainSessionsFacade.GetSessionById(id)
	return sess, err
}

func (actions *AwsPlainSessionActions) UpdateAwsPlainSession(sessionId string, name string, accountNumber string, region string, user string,
	awsAccessKeyId string, awsSecretAccessKey string, mfaDevice string, profile string) error {

	/*
		config, err := configuration.ReadConfiguration()
		if err != nil {
			return err
		}

		err = session.UpdateAwsPlainSession(config, sessionId, name, accountNumber, region, user, awsAccessKeyId, awsSecretAccessKey, mfaDevice, profile)
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

func DeleteAwsPlainSession(sessionId string) error {
	/*
		config, err := configuration.ReadConfiguration()
		if err != nil {
			return err
		}

		err = session.DeleteAwsPlainSession(config, sessionId)
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

func (actions *AwsPlainSessionActions) StartAwsPlainSession(sessionId string) error {
	awsPlainSession, err := actions.AwsPlainSessionsFacade.GetSessionById(sessionId)
	if err != nil {
		return err
	}

	doesSessionTokenExist, err := actions.Keychain.DoesSecretExist(awsPlainSession.Id + "-plain-aws-session-session-token")
	if err != nil {
		return err
	}

	if doesSessionTokenExist {
		sessionTokenExpiration := awsPlainSession.Account.SessionTokenExpiration

		if sessionTokenExpiration != "" {
			currentTime := time.Now()
			sessionTokenExpirationTime, err := time.Parse(time.RFC3339, sessionTokenExpiration)
			if err != nil {
				return err
			}

			if currentTime.After(sessionTokenExpirationTime) {
				err = actions.generateSessionToken(*awsPlainSession)
				if err != nil {
					return err
				}
			}
		} else {
			err = actions.generateSessionToken(*awsPlainSession)
			if err != nil {
				return err
			}
		}
	} else {
		err = actions.generateSessionToken(*awsPlainSession)
		if err != nil {
			return err
		}
	}

	err = actions.AwsPlainSessionsFacade.SetSessionStatusToPending(sessionId)
	if err != nil {
		return err
	}

	err = actions.AwsPlainSessionsFacade.SetSessionStatusToActive(sessionId)
	if err != nil {
		return err
	}

	return nil
}

func StopAwsPlainSession(sessionId string) error {
	/*
			config, err := configuration.ReadConfiguration()
			if err != nil {
				return err
			}

			// Passing nil because, it will be the rotate method to check if we need the mfaToken or not
			err = session.StopAwsPlainSession(config, sessionId)
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

// TODO: encapsulate this logic inside a session token generation interface
func (actions *AwsPlainSessionActions) generateSessionToken(awsPlainSession session.AwsPlainSession) error {
	accessKeyIdSecretName := awsPlainSession.Id + "-plain-aws-session-access-key-id"

	accessKeyId, err := actions.Keychain.GetSecret(accessKeyIdSecretName)
	if err != nil {
		return http_error.NewUnprocessableEntityError(err)
	}

	secretAccessKeySecretName := awsPlainSession.Id + "-plain-aws-session-secret-access-key"

	secretAccessKey, err := actions.Keychain.GetSecret(secretAccessKeySecretName)
	if err != nil {
		return http_error.NewUnprocessableEntityError(err)
	}

	credentials, err := sts_client.GenerateAccessToken(awsPlainSession.Account.Region,
		awsPlainSession.Account.MfaDevice, nil, accessKeyId, secretAccessKey)
	if err != nil {
		return err
	}

	credentialsJson, err := json.Marshal(credentials)
	if err != nil {
		return err
	}

	err = actions.Keychain.SetSecret(string(credentialsJson),
		awsPlainSession.Id+"-plain-aws-session-session-token")
	if err != nil {
		return err
	}

	err = actions.AwsPlainSessionsFacade.SetSessionTokenExpiration(awsPlainSession.Id, *credentials.Expiration)
	if err != nil {
		return err
	}

	return nil
}
