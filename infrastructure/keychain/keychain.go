// TODO: probably, it should be moved in interface layer, since it acts as an interface between a domain concept,
//  e.g. access keys, and an external service.

package keychain

import (
  "github.com/zalando/go-keyring"
  "leapp_daemon/infrastructure/http/http_error"
)

const ServiceName = "Leapp"

func SetSecret(secret string, label string) error {
	err := keyring.Set(ServiceName, label, secret)
	if err != nil {
		return http_error.NewUnprocessableEntityError(err)
	}
	return nil
}

func GetSecret(label string) (string, error) {
	secret, err := keyring.Get(ServiceName, label)
	if err != nil {
		return "", http_error.NewNotFoundError(err)
	}
	return secret, nil
}

func DoesSecretExist(label string) (bool, error) {
	_, err := keyring.Get(ServiceName, label)
	if err != nil {
		if err.Error() == "secret not found in keyring" {
			return false, nil
		}
		return false, http_error.NewUnprocessableEntityError(err)
	}
	return true, nil
}
