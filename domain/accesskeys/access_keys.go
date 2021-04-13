package accesskeys

import (
	"leapp_daemon/infrastructure/http/http_error"
	"leapp_daemon/infrastructure/keychain"
)

type AccessKeysRepository interface {
  Get(accountName string) (AccessKeys, error)
  Store(accessKeys AccessKeys) error
}

type AccessKeys struct {
  accessKeyId string
  secretAccessKey string
}

func(accessKeys *AccessKeys) SetAccessKeyId(accessKeyId string) {
  accessKeys.accessKeyId = accessKeyId
}

func(accessKeys *AccessKeys) SetSecretAccessKey(secretAccessKey string) {
  accessKeys.secretAccessKey = secretAccessKey
}

func Get(accountName string) (string, string, error) {
	accessKeyIdSecretName := accountName + "-plain-aws-session-access-key-id"

	accessKeyId, err := keychain.GetSecret(accessKeyIdSecretName)
	if err != nil {
		return "", "", http_error.NewUnprocessableEntityError(err)
	}

	secretAccessKeySecretName := accountName + "-plain-aws-session-secret-access-key"

	secretAccessKey, err := keychain.GetSecret(secretAccessKeySecretName)
	if err != nil {
		return "", "", http_error.NewUnprocessableEntityError(err)
	}

	return accessKeyId, secretAccessKey, nil
}
