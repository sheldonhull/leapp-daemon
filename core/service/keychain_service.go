package service

import (
	"github.com/zalando/go-keyring"
	"leapp_daemon/shared/constant"
)

func SaveSecret(secret string, label string) error {
	err := keyring.Set(constant.KeychainService, label, secret)
	if err != nil {
		return err
	}
	return nil
}

func RetrieveSecret(label string) (string, error) {
	secret, err := keyring.Get(constant.KeychainService, label)
	if err != nil {
		return "", err
	}
	return secret, nil
}

func DoesSecretExist(label string) (bool, error) {
	_, err := keyring.Get(constant.KeychainService, label)
	if err != nil {
		if err.Error() == "secret not found in keyring" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}