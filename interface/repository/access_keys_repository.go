package repository

import (
  "leapp_daemon/domain/access_keys"
  "leapp_daemon/infrastructure/http/http_error"
  "leapp_daemon/infrastructure/keychain"
)

type AccessKeyIdRepository struct {}

func(repository *AccessKeyIdRepository) Get(accountName string) (access_keys.AccessKeys, error) {
  var accessKeys access_keys.AccessKeys
  accessKeyIdSecretName := accountName + "-plain-aws-session-access-key-id"

  accessKeyId, err := keychain.GetSecret(accessKeyIdSecretName)
  if err != nil {
    return accessKeys, http_error.NewUnprocessableEntityError(err)
  }

  secretAccessKeySecretName := accountName + "-plain-aws-session-secret-access-key"

  secretAccessKey, err := keychain.GetSecret(secretAccessKeySecretName)
  if err != nil {
    return accessKeys, http_error.NewUnprocessableEntityError(err)
  }

  accessKeys.SetAccessKeyId(accessKeyId)
  accessKeys.SetSecretAccessKey(secretAccessKey)

  return accessKeys, nil
}

func(repository *AccessKeyIdRepository) Store(accessKeys access_keys.AccessKeys) error {
  return nil
}
