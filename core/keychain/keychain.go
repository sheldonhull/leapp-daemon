package keychain

import (
	"github.com/zalando/go-keyring"
	"leapp_daemon/custom_error"
)

const ServiceName = "Leapp"

func SetSecret(secret string, label string) error {
	err := keyring.Set(ServiceName, label, secret)
	if err != nil {
		return custom_error.NewUnprocessableEntityError(err)
	}
	return nil
}

func GetSecret(label string) (string, error) {
	secret, err := keyring.Get(ServiceName, label)
	if err != nil {
		return "", custom_error.NewNotFoundError(err)
	}
	return secret, nil
}

func DoesSecretExist(label string) (bool, error) {
	_, err := keyring.Get(ServiceName, label)
	if err != nil {
		if err.Error() == "secret not found in keyring" {
			return false, nil
		}
		return false, custom_error.NewUnprocessableEntityError(err)
	}
	return true, nil
}
