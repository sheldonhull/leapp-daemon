package access_keys

import (
	"leapp_daemon/core/keychain"
	"leapp_daemon/shared/custom_error"
)

func Get(accountName string) (string, string, error) {
	accessKeyIdSecretName := accountName + "-plain-aws-session-access-key-id"

	accessKeyId, err := keychain.RetrieveSecret(accessKeyIdSecretName)
	if err != nil {
		return "", "", custom_error.NewUnprocessableEntityError(err)
	}

	secretAccessKeySecretName := accountName + "-plain-aws-session-secret-access-key"

	secretAccessKey, err := keychain.RetrieveSecret(secretAccessKeySecretName)
	if err != nil {
		return "", "", custom_error.NewUnprocessableEntityError(err)
	}

	return accessKeyId, secretAccessKey, nil
}
