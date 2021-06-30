package use_case

import (
	"encoding/json"
	"leapp_daemon/domain/aws"
	"leapp_daemon/domain/aws/aws_iam_user"
	"leapp_daemon/infrastructure/http/http_error"
	"time"
)

type AwsIamUserSessionActions struct {
	Environment              Environment
	Keychain                 Keychain
	StsApi                   StsApi
	NamedProfilesActions     NamedProfilesActionsInterface
	AwsIamUserSessionsFacade AwsIamUserSessionsFacade
}

func (actions *AwsIamUserSessionActions) GetSession(sessionId string) (aws_iam_user.AwsIamUserSession, error) {
	return actions.AwsIamUserSessionsFacade.GetSessionById(sessionId)
}

func (actions *AwsIamUserSessionActions) CreateSession(name string, awsAccessKeyId string, awsSecretKey string,
	mfaDevice string, region string, profileName string) error {

	newSessionId := actions.Environment.GenerateUuid()
	accessKeyIdLabel := newSessionId + "-aws-iam-user-session-access-key-id"
	secretKeyLabel := newSessionId + "-aws-iam-user-session-secret-key"
	sessionTokenExpirationLabel := newSessionId + "-aws-iam-user-session-token-expiration"
	namedProfile, err := actions.NamedProfilesActions.GetOrCreateNamedProfile(profileName)
	if err != nil {
		return err
	}

	sess := aws_iam_user.AwsIamUserSession{
		Id:                     newSessionId,
		Name:                   name,
		Region:                 region,
		NamedProfileId:         namedProfile.Id,
		MfaDevice:              mfaDevice,
		AccessKeyIdLabel:       accessKeyIdLabel,
		SecretKeyLabel:         secretKeyLabel,
		SessionTokenLabel:      sessionTokenExpirationLabel,
		Status:                 aws.NotActive,
		SessionTokenExpiration: "",
		StartTime:              "",
		LastStopTime:           "",
	}

	err = actions.Keychain.SetSecret(awsAccessKeyId, accessKeyIdLabel)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = actions.Keychain.SetSecret(awsSecretKey, secretKeyLabel)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return actions.AwsIamUserSessionsFacade.AddSession(sess)
}

func (actions *AwsIamUserSessionActions) StartSession(sessionId string) error {
	facade := actions.AwsIamUserSessionsFacade

	sessionToStart, err := facade.GetSessionById(sessionId)
	if err != nil {
		return err
	}

	err = facade.StartingSession(sessionId)
	if err != nil {
		return err
	}

	currentTime := actions.Environment.GetTime()
	err = actions.refreshSessionTokenIfNeeded(sessionToStart, currentTime)
	if err != nil {
		goto StartSessionFailed
	}

	err = facade.StartSession(sessionId, currentTime)
	if err != nil {
		goto StartSessionFailed
	}

	return nil

StartSessionFailed:
	facade.StopSession(sessionId, currentTime)
	return err
}

func (actions *AwsIamUserSessionActions) StopSession(sessionId string) error {
	return actions.AwsIamUserSessionsFacade.StopSession(sessionId, actions.Environment.GetTime())
}

func (actions *AwsIamUserSessionActions) DeleteSession(sessionId string) error {
	facade := actions.AwsIamUserSessionsFacade

	sessionToDelete, err := facade.GetSessionById(sessionId)
	if err != nil {
		return err
	}

	_ = actions.Keychain.DeleteSecret(sessionToDelete.AccessKeyIdLabel)
	_ = actions.Keychain.DeleteSecret(sessionToDelete.SecretKeyLabel)
	_ = actions.Keychain.DeleteSecret(sessionToDelete.SessionTokenLabel)
	return facade.RemoveSession(sessionId)
}

func (actions *AwsIamUserSessionActions) EditAwsIamUserSession(sessionId string, sessionName string, accountNumber string,
	region string, user string, awsAccessKeyId string, awsSecretKey string, mfaDevice string, namedProfile string) error {
	/*
		return actions.AwsIamUserSessionsFacade.EditSession(sessionId, sessionName, region, mfaDevice, namedProfile)

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

func (actions *AwsIamUserSessionActions) refreshSessionTokenIfNeeded(session aws_iam_user.AwsIamUserSession, currentTime string) error {
	if !actions.isSessionTokenValid(session.SessionTokenLabel, session.SessionTokenExpiration, currentTime) {
		err := actions.refreshSessionToken(session)
		if err != nil {
			return err
		}
	}

	return nil
}

func (actions *AwsIamUserSessionActions) isSessionTokenValid(sessionTokenLabel string, sessionTokenExpiration string, currentTime string) bool {
	isSessionTokenStoredIntoKeychain, err := actions.Keychain.DoesSecretExist(sessionTokenLabel)
	if err != nil || !isSessionTokenStoredIntoKeychain {
		return false
	}

	if sessionTokenExpiration == "" {
		return false
	}

	sessionCurrentTime, err := time.Parse(time.RFC3339, currentTime)
	if err != nil {
		return false
	}

	sessionTokenExpirationTime, err := time.Parse(time.RFC3339, sessionTokenExpiration)
	if err != nil {
		return false
	}

	if sessionCurrentTime.After(sessionTokenExpirationTime) {
		return false
	}

	return true
}

func (actions *AwsIamUserSessionActions) refreshSessionToken(session aws_iam_user.AwsIamUserSession) error {
	accessKeyId, err := actions.Keychain.GetSecret(session.AccessKeyIdLabel)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	secretKey, err := actions.Keychain.GetSecret(session.SecretKeyLabel)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	credentials, err := actions.StsApi.GenerateNewSessionToken(accessKeyId, secretKey, session.Region, session.MfaDevice, nil)
	if err != nil {
		return err
	}

	credentialsJson, err := json.Marshal(credentials)
	if err != nil {
		return err
	}

	err = actions.Keychain.SetSecret(string(credentialsJson), session.SessionTokenLabel)
	if err != nil {
		return err
	}

	return actions.AwsIamUserSessionsFacade.SetSessionTokenExpiration(session.Id, credentials.Expiration.Format(time.RFC3339))
}
