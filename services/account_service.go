package services

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"leapp_daemon/custom_errors"
	"leapp_daemon/services/domain"
	"strings"
)

func CreateFederatedAwsAccount(name string, accountNumber string, roleName string, roleArn string, idpArn string,
	region string, ssoUrl string) error {
	configuration, err := ReadConfiguration()
	if err != nil {
		return err
	}

	sessions := configuration.FederatedAwsAccountSessions

	for _, session := range sessions {
		account := session.Account
		if account.AccountNumber == accountNumber && account.Role.Name == roleName {
			err = custom_errors.NewBadRequestError(errors.New("an account with the same account number and " +
				"role name is already present"))
			return err
		}
	}

	role := domain.FederatedAwsRole{
		Name: roleName,
		Arn:  roleArn,
	}

	federatedAwsAccount := domain.FederatedAwsAccount{
		AccountNumber: accountNumber,
		Name:          name,
		Role:          role,
		IdpArn:        idpArn,
		Region:        region,
		SsoUrl:        ssoUrl,
	}

	uuidString := uuid.New().String()
	uuidString = strings.Replace(uuidString, "-", "", -1)

	session := domain.FederatedAwsAccountSession{
		Id:           uuidString,
		Active:       false,
		Loading:      false,
		LastStopDate: "",
		Account:      federatedAwsAccount,
	}

	configuration.FederatedAwsAccountSessions = append(configuration.FederatedAwsAccountSessions, session)

	err = UpdateConfiguration(configuration, false)
	if err != nil {
		return err
	}

	return nil
}
