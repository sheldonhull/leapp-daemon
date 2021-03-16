package session_token

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/sts"
	"gopkg.in/ini.v1"
	"leapp_daemon/constant"
	"leapp_daemon/core/access_keys"
	"leapp_daemon/core/aws_client"
	"leapp_daemon/core/file_system"
	"leapp_daemon/core/ini_file"
	"leapp_daemon/core/keychain"
	"leapp_daemon/custom_error"
	"sync"
	"time"
)

// The zero value is an unlocked mutex
var keychainMutex sync.Mutex

// The zero value is an unlocked mutex
var iniFileMutex sync.Mutex

func DoExist(accountName string) (bool, error) {
	doesSessionTokenExpirationExist, err := keychain.DoesSecretExist(accountName + "-plain-aws-session-session-token-expiration")
	if err != nil {
		return false, err
	}

	doesSessionTokenExist, err := keychain.DoesSecretExist(accountName + "-plain-aws-session-session-token")
	if err != nil {
		return false, err
	}

	return doesSessionTokenExpirationExist && doesSessionTokenExist, nil
}

func IsExpired(accountName string) (bool, error) {
	secret, err := keychain.RetrieveSecret(accountName + "-plain-aws-session-session-token-expiration")
	if err != nil {
		return false, err
	}

	t, err := time.Parse(time.RFC3339, secret)
	if err != nil {
		return false, err
	}

	return time.Now().After(t), nil
}

func Generate(accountName string, region string, mfaDevice string, mfaToken *string) (*sts.Credentials, error) {
	accessKeyId, secretAccessKey, err := access_keys.Get(accountName)
	if err != nil {
		return nil, err
	}

	stsClient, err := aws_client.GetStaticCredentialsClient(accessKeyId, secretAccessKey, &region)
	if err != nil {
		return nil, err
	}

	durationSeconds := constant.SessionTokenDurationInSeconds
	var getSessionTokenInput sts.GetSessionTokenInput

	if mfaToken == nil {
		getSessionTokenInput = sts.GetSessionTokenInput{
			DurationSeconds: &durationSeconds,
			SerialNumber:    nil,
			TokenCode:       nil,
		}
	} else {
		getSessionTokenInput = sts.GetSessionTokenInput{
			DurationSeconds: &durationSeconds,
			SerialNumber:    &mfaDevice,
			TokenCode:       mfaToken,
		}
	}

	getSessionTokenOutput, err := stsClient.GetSessionToken(&getSessionTokenInput)
	if err != nil {
		return nil, err
	}

	return getSessionTokenOutput.Credentials, nil
}

func SaveInKeychain(accountName string, credentials *sts.Credentials) error {
	keychainMutex.Lock()
	defer keychainMutex.Unlock()

	credentialsJson, err := json.Marshal(credentials)
	if err != nil {
		return err
	}

	err = keychain.SaveSecret(string(credentialsJson),
		accountName + "-plain-aws-session-session-token")
	if err != nil {
		return err
	}

	err = keychain.SaveSecret(credentials.Expiration.Format(time.RFC3339),
		accountName + "-plain-aws-session-session-token-expiration")
	if err != nil {
		return err
	}

	return nil
}

func SaveInIniFile(accessKeyId string, secretAccessKey string, sessionToken string, region string, profileName string) error {
	iniFileMutex.Lock()
	defer iniFileMutex.Unlock()

	homeDir, err := file_system.GetHomeDir()
	if err != nil {
		return err
	}

	credentialsFilePath := homeDir + "/" + constant.CredentialsFilePath

	if file_system.DoesFileExist(credentialsFilePath) {
		credentialsFile, err := ini.Load(credentialsFilePath)
		if err != nil {
			return err
		}

		section, err := credentialsFile.GetSection(profileName)
		if err != nil {
			return err
		}

		if section == nil {
			_, err = ini_file.CreateNamedProfileSection(credentialsFile, profileName, accessKeyId, secretAccessKey,
				sessionToken, region)
			if err != nil {
				return err
			}

			err = ini_file.AppendToFile(credentialsFile, credentialsFilePath)
			if err != nil {
				return err
			}
		} else {
			credentialsFile.DeleteSection(profileName)

			_, err = ini_file.CreateNamedProfileSection(credentialsFile, profileName, accessKeyId, secretAccessKey,
				sessionToken, region)
			if err != nil {
				return err
			}

			err = ini_file.OverwriteFile(credentialsFile, credentialsFilePath)
			if err != nil {
				return err
			}
		}
	} else {
		credentialsFile := ini.Empty()

		_, err = ini_file.CreateNamedProfileSection(credentialsFile, profileName, accessKeyId, secretAccessKey,
			sessionToken, region)
		if err != nil {
			return err
		}

		err = ini_file.OverwriteFile(credentialsFile, credentialsFilePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func Get(accountName string) (string, string, error) {
	sessionTokenSecretName := accountName + "-plain-aws-session-session-token"

	sessionToken, err := keychain.RetrieveSecret(sessionTokenSecretName)
	if err != nil {
		return "", "", custom_error.NewUnprocessableEntityError(err)
	}

	sessionTokenExpirationSecretName := accountName + "-plain-aws-session-session-token-expiration"

	sessionTokenExpiration, err := keychain.RetrieveSecret(sessionTokenExpirationSecretName)
	if err != nil {
		return "", "", custom_error.NewUnprocessableEntityError(err)
	}

	return sessionToken, sessionTokenExpiration, nil
}
