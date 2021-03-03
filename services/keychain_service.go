package services

import (
	"log"
	"github.com/zalando/go-keyring"
	"leapp_daemon/services/domain"
)

func SaveSecret(secret string, label string) error {
	err := keyring.Set(domain.KeychainService, label, secret)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func RetrieveSecret(label string) (string, error) {
	secret, err := keyring.Get(domain.KeychainService, label)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return secret, nil
}
