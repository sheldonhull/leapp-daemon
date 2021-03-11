package service

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"gopkg.in/ini.v1"
	"leapp_daemon/core/model"
	"leapp_daemon/core/util"
	"leapp_daemon/shared/constant"
	"leapp_daemon/shared/custom_error"
	"os"
	"sync"
	"time"
)

// The zero value is an unlocked mutex
var keychainMutex sync.Mutex

// The zero value is an unlocked mutex
var iniFileMutex sync.Mutex

func getStsStaticCredentialsClient(accessKeyId string, secretAccessKey string, region *string) (*sts.STS, error) {

	endpoint, err := util.GetRegionalEndpoint(region)
	if err != nil {
		return nil, err
	}

	stsConfig := &aws.Config{
			Credentials: credentials.NewStaticCredentials(accessKeyId, secretAccessKey, ""),
			Region: region,
			Endpoint: endpoint,
	}

	sess, err := session.NewSession(stsConfig)
	stsClient := sts.New(session.Must(sess, err))

	return stsClient, nil
}

func GenerateSessionToken(sess *model.PlainAwsSession, mfaToken *string) (*sts.Credentials, error) {
	accessKeyId, secretAccessKey, err := GetAccessKeys(sess.Account.Name)
	if err != nil {
		return nil, err
	}

	stsClient, err2 := getStsStaticCredentialsClient(accessKeyId, secretAccessKey, &sess.Account.Region)
	if err2 != nil {
		return nil, err2
	}

	durationSeconds := constant.SessionTokenDurationInSeconds
	var getSessionTokenInput sts.GetSessionTokenInput

	if *mfaToken == "" {
		getSessionTokenInput = sts.GetSessionTokenInput{
			DurationSeconds: &durationSeconds,
			SerialNumber:    nil,
			TokenCode:       nil,
		}
	} else {
		getSessionTokenInput = sts.GetSessionTokenInput{
			DurationSeconds: &durationSeconds,
			SerialNumber:    &sess.Account.MfaDevice,
			TokenCode:       mfaToken,
		}
	}

	getSessionTokenOutput, err3 := stsClient.GetSessionToken(&getSessionTokenInput)
	if err3 != nil {
		return nil, err3
	}

	return getSessionTokenOutput.Credentials, nil
}

func SaveSessionTokenInKeychain(accountName string, credentials *sts.Credentials) error {
	keychainMutex.Lock()
	defer keychainMutex.Unlock()

	credentialsJson, err := json.Marshal(credentials)
	if err != nil {
		return err
	}

	err = SaveSecret(string(credentialsJson),
		accountName + "-plain-aws-session-session-token")
	if err != nil {
		return err
	}

	err = SaveSecret(credentials.Expiration.Format(time.RFC3339),
		accountName + "-plain-aws-session-session-token-expiration")
	if err != nil {
		return err
	}

	return nil
}

func DoSessionTokenExist(accountName string) (bool, error) {
	doesSessionTokenExpirationExist, err := DoesSecretExist(accountName + "-plain-aws-session-session-token-expiration")
	if err != nil {
		return false, err
	}

	doesSessionTokenExist, err := DoesSecretExist(accountName + "-plain-aws-session-session-token")
	if err != nil {
		return false, err
	}

	return doesSessionTokenExpirationExist && doesSessionTokenExist, nil
}

func IsSessionTokenExpired(accountName string) (bool, error) {
	secret, err := RetrieveSecret(accountName + "-plain-aws-session-session-token-expiration")
	if err != nil {
		return false, err
	}

	t, err := time.Parse(time.RFC3339, secret)
	if err != nil {
		return false, err
	}

	return time.Now().After(t), nil
}

func SaveSessionTokenInIniFile(accessKeyId string, secretAccessKey string, sessionToken string, region string, profileName string) error {
	iniFileMutex.Lock()
	defer iniFileMutex.Unlock()

	homeDir, err := GetHomeDir()
	if err != nil {
		return err
	}

	credentialsFilePath := homeDir + "/" + constant.CredentialsFilePath

	if DoesFileExist(credentialsFilePath) {
		credentialsFile, err := ini.Load(credentialsFilePath)
		if err != nil {
			return err
		}

		section, err := credentialsFile.GetSection(profileName)
		if err != nil {
			return err
		}

		if section == nil {
			_, err = createNamedProfileSection(credentialsFile, profileName, accessKeyId, secretAccessKey,
				sessionToken, region)
			if err != nil {
				return err
			}

			err = appendToFile(credentialsFile, credentialsFilePath)
			if err != nil {
				return err
			}
		} else {
			credentialsFile.DeleteSection(profileName)

			_, err = createNamedProfileSection(credentialsFile, profileName, accessKeyId, secretAccessKey,
				sessionToken, region)
			if err != nil {
				return err
			}

			err = overwriteFile(credentialsFile, credentialsFilePath)
			if err != nil {
				return err
			}
		}
	} else {
		credentialsFile := ini.Empty()

		_, err = createNamedProfileSection(credentialsFile, profileName, accessKeyId, secretAccessKey,
			sessionToken, region)
		if err != nil {
			return err
		}

		err = overwriteFile(credentialsFile, credentialsFilePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetSessionToken(accountName string) (string, string, error) {
	sessionTokenSecretName := accountName + "-plain-aws-session-session-token"

	sessionToken, err := RetrieveSecret(sessionTokenSecretName)
	if err != nil {
		return "", "", custom_error.NewUnprocessableEntityError(err)
	}

	sessionTokenExpirationSecretName := accountName + "-plain-aws-session-session-token-expiration"

	sessionTokenExpiration, err := RetrieveSecret(sessionTokenExpirationSecretName)
	if err != nil {
		return "", "", custom_error.NewUnprocessableEntityError(err)
	}

	return sessionToken, sessionTokenExpiration, nil
}

func GetAccessKeys(accountName string) (string, string, error) {
	accessKeyIdSecretName := accountName + "-plain-aws-session-access-key-id"

	accessKeyId, err := RetrieveSecret(accessKeyIdSecretName)
	if err != nil {
		return "", "", custom_error.NewUnprocessableEntityError(err)
	}

	secretAccessKeySecretName := accountName + "-plain-aws-session-secret-access-key"

	secretAccessKey, err := RetrieveSecret(secretAccessKeySecretName)
	if err != nil {
		return "", "", custom_error.NewUnprocessableEntityError(err)
	}

	return accessKeyId, secretAccessKey, nil
}

func createNamedProfileSection(credentialsFile *ini.File, profileName string, accessKeyId string,
	secretAccessKey string, sessionToken string, region string) (*ini.Section, error) {

	section, err := credentialsFile.NewSection(profileName)
	if err != nil {
		return nil, err
	}

	_, err = section.NewKey("aws_access_key_id", accessKeyId)
	if err != nil {
		return nil, err
	}

	_, err = section.NewKey("aws_secret_access_key", secretAccessKey)
	if err != nil {
		return nil, err
	}

	_, err = section.NewKey("aws_session_token", sessionToken)
	if err != nil {
		return nil, err
	}

	if region != "" {
		_, err = section.NewKey("region", region)
		if err != nil {
			return nil, err
		}
	}

	return section, nil
}

func overwriteFile(file *ini.File, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	_, err = file.WriteTo(f)
	if err != nil {
		return err
	}

	return nil
}

func appendToFile(file *ini.File, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		return err
	}

	_, err = file.WriteTo(f)
	if err != nil {
		return err
	}

	return nil
}
