package service

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"leapp_daemon/core/configuration"
	"leapp_daemon/custom_error"
	"strings"
)

// TODO: move into configuration package

func GetFederatedAwsSession(id string) (*configuration.FederatedAwsSession, error) {
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return nil, err
	}

	sessions := config.FederatedAwsSessions
	for index, _ := range sessions {
		if sessions[index].Id == id {
			return sessions[index], nil
		}
	}

	return nil, custom_error.NewBadRequestError(errors.New("No session found with id:" + id))
}

func ListFederatedAwsSession(query string) ([]*configuration.FederatedAwsSession, error) {
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return nil, err
	}

	filteredList := make([]*configuration.FederatedAwsSession, 0)

	if query == "" {
		return append(filteredList, config.FederatedAwsSessions...), nil
	} else {
		for _, session := range config.FederatedAwsSessions {
			if  strings.Contains(session.Id, query) ||
			    strings.Contains(session.Account.Name, query) ||
				strings.Contains(session.Account.IdpArn, query) ||
				strings.Contains(session.Account.SsoUrl, query) ||
				strings.Contains(session.Account.Region, query) ||
				strings.Contains(session.Account.AccountNumber, query) {
				// TODO: add also role filters
				filteredList = append(filteredList, session)
			}
		}

		return filteredList, nil
	}
}

func CreateFederatedAwsSession(name string, accountNumber string, roleName string, roleArn string, idpArn string,
	region string, ssoUrl string) error {
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return err
	}

	sessions := config.FederatedAwsSessions

	for _, session := range sessions {
		account := session.Account
		if account.AccountNumber == accountNumber && account.Role.Name == roleName {
			err = custom_error.NewBadRequestError(errors.New("an account with the same account number and " +
				"role name is already present"))
			return err
		}
	}

	role := configuration.FederatedAwsRole{
		Name: roleName,
		Arn:  roleArn,
	}

	federatedAwsAccount := configuration.FederatedAwsAccount{
		AccountNumber: accountNumber,
		Name:          name,
		Role:          &role,
		IdpArn:        idpArn,
		Region:        region,
		SsoUrl:        ssoUrl,
	}

	uuidString := uuid.New().String()
	uuidString = strings.Replace(uuidString, "-", "", -1)

	session := configuration.FederatedAwsSession{
		Id:           uuidString,
		Active:       false,
		Loading:      false,
		StartTime: "",
		Account:      &federatedAwsAccount,
	}

	config.FederatedAwsSessions = append(config.FederatedAwsSessions, &session)

	err = configuration.UpdateConfiguration(config, false)
	if err != nil {
		return err
	}

	return nil
}