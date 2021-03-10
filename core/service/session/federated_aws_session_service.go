package session

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"leapp_daemon/core/model"
	"leapp_daemon/core/service"
	errors2 "leapp_daemon/shared/custom_error"
	"strings"
)


func GetFederatedAwsSession(id string) (*model.FederatedAwsSession, error) {
	configuration, err := service.ReadConfiguration()
	if err != nil {
		return nil, err
	}

	sessions := configuration.FederatedAwsSessions
	for index, _ := range sessions {
		if sessions[index].Id == id {
			return &sessions[index], nil
		}
	}

	return nil, errors2.NewBadRequestError(errors.New("No session found with id:" + id))
}

func ListFederatedAwsSession(query string) ([]model.FederatedAwsSession, error) {
	configuration, err := service.ReadConfiguration()
	if err != nil {
		return nil, err
	}

	filteredList := make([]model.FederatedAwsSession, 0)

	if query == "" {
		return append(filteredList, configuration.FederatedAwsSessions...), nil
	} else {
		for _, session := range configuration.FederatedAwsSessions {
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
	configuration, err := service.ReadConfiguration()
	if err != nil {
		return err
	}

	sessions := configuration.FederatedAwsSessions

	for _, session := range sessions {
		account := session.Account
		if account.AccountNumber == accountNumber && account.Role.Name == roleName {
			err = errors2.NewBadRequestError(errors.New("an account with the same account number and " +
				"role name is already present"))
			return err
		}
	}

	role := model.FederatedAwsRole{
		Name: roleName,
		Arn:  roleArn,
	}

	federatedAwsAccount := model.FederatedAwsAccount{
		AccountNumber: accountNumber,
		Name:          name,
		Role:          role,
		IdpArn:        idpArn,
		Region:        region,
		SsoUrl:        ssoUrl,
	}

	uuidString := uuid.New().String()
	uuidString = strings.Replace(uuidString, "-", "", -1)

	session := model.FederatedAwsSession{
		Id:           uuidString,
		Active:       false,
		Loading:      false,
		LastStopDate: "",
		Account:      federatedAwsAccount,
	}

	configuration.FederatedAwsSessions = append(configuration.FederatedAwsSessions, session)

	err = service.UpdateConfiguration(configuration, false)
	if err != nil {
		return err
	}

	return nil
}