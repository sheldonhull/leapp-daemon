package keychain

import (
  "github.com/zalando/go-keyring"
  http_error2 "leapp_daemon/infrastructure/http/http_error"
)

const ServiceName = "Leapp"

func SetSecret(secret string, label string) error {
	err := keyring.Set(ServiceName, label, secret)
	if err != nil {
		return http_error2.NewUnprocessableEntityError(err)
	}
	return nil
}

func GetSecret(label string) (string, error) {
	secret, err := keyring.Get(ServiceName, label)
	if err != nil {
		return "", http_error2.NewNotFoundError(err)
	}
	return secret, nil
}

func DoesSecretExist(label string) (bool, error) {
	_, err := keyring.Get(ServiceName, label)
	if err != nil {
		if err.Error() == "secret not found in keyring" {
			return false, nil
		}
		return false, http_error2.NewUnprocessableEntityError(err)
	}
	return true, nil
}
