package service

import (
	"github.com/zalando/go-keyring"
	"leapp_daemon/shared/const"
)

// TODO: do not use log.Fatal, since it will call os.Exit(1)
func SaveSecret(secret string, label string) error {
	err := keyring.Set(_const.KeychainService, label, secret)
	if err != nil {
		//log.Fatal(err)
		return err
	}
	return nil
}

// TODO: do not use log.Fatal, since it will call os.Exit(1)
func RetrieveSecret(label string) (string, error) {
	secret, err := keyring.Get(_const.KeychainService, label)
	if err != nil {
		//log.Fatal(err)
		return "", err
	}
	return secret, nil
}

func DoesSecretExist(label string) (bool, error) {
	_, err := keyring.Get(_const.KeychainService, label)
	if err != nil {
		if err.Error() == "secret not found in keyring" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}