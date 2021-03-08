package session

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"leapp_daemon/custom_errors"
	"leapp_daemon/service"
	"leapp_daemon/service/domain"
	"strings"
)

func GetPlainAwsSession(id string) (*domain.PlainAwsAccountSession, error) {
	configuration, err := service.ReadConfiguration()
	if err != nil {
		return nil, err
	}

	sessions := configuration.PlainAwsAccountSessions
	for index, _ := range sessions {
		if sessions[index].Id == id {
			return &sessions[index], nil
		}
	}

	return nil, custom_errors.NewBadRequestError(errors.New("No session found with id:" + id))
}

func ListPlainAwsSession(query string) ([]domain.PlainAwsAccountSession, error) {
	configuration, err := service.ReadConfiguration()
	if err != nil {
		return nil, err
	}

	filteredList := make([]domain.PlainAwsAccountSession, 0)

	if query == "" {
		return append(filteredList, configuration.PlainAwsAccountSessions...), nil
	} else {
		for _, session := range configuration.PlainAwsAccountSessions {
			if strings.Contains(session.Id, query) ||
			   strings.Contains(session.Account.Name, query) ||
			   strings.Contains(session.Account.MfaDevice, query) ||
			   strings.Contains(session.Account.User, query) ||
			   strings.Contains(session.Account.Region, query) ||
			   strings.Contains(session.Account.AccountNumber, query) {

				filteredList = append(filteredList, session)
			}
		}

		return filteredList, nil
	}
}

func CreatePlainAwsSession(name string, accountNumber string, region string, user string, mfaDevice string) error {

	configuration, err := service.ReadConfiguration()
	if err != nil {
		return err
	}

	sessions := configuration.PlainAwsAccountSessions

	for _, session := range sessions {
		account := session.Account
		if account.AccountNumber == accountNumber && account.User == user {
			err = custom_errors.NewBadRequestError(errors.New("an account with the same account number and user is already present"))
			return err
		}
	}

	plainAwsAccount := domain.PlainAwsAccount {
		AccountNumber: accountNumber,
		Name:          name,
		Region:        region,
		User:          user,
		MfaDevice:     mfaDevice,
	}

	uuidString := uuid.New().String()
	uuidString = strings.Replace(uuidString, "-", "", -1)

	session := domain.PlainAwsAccountSession {
		Id:           uuidString,
		Active:       false,
		Loading:      false,
		LastStopDate: "",
		Account:      plainAwsAccount,
	}

	configuration.PlainAwsAccountSessions = append(configuration.PlainAwsAccountSessions, session)

	err = service.UpdateConfiguration(configuration, false)
	if err != nil {
		return err
	}

	return nil
}

func EditPlainAwsSession(id string, name string, accountNumber string, region string, user string, mfaDevice string) error {

	configuration, err := service.ReadConfiguration()
	if err != nil {
		return err
	}

	sessions := configuration.PlainAwsAccountSessions

	found := false
	for index, _ := range sessions {
		if sessions[index].Id == id {
			sessions[index].Account = domain.PlainAwsAccount {
				AccountNumber: accountNumber,
				Name:          name,
				Region:        region,
				User:          user,
				MfaDevice:     mfaDevice,
			}
			found = true
		}
	}

	if found == false {
		err = custom_errors.NewBadRequestError(errors.New("Session not found for Id: " + id))
		return err
	}

	configuration.PlainAwsAccountSessions = sessions

	err = service.UpdateConfiguration(configuration, false)
	if err != nil {
		return err
	}

	return nil
}

func DeletePlainAwsSession(id string) error {
	configuration, err := service.ReadConfiguration()
	if err != nil {
		return err
	}

	sessions := configuration.PlainAwsAccountSessions

	found := false
	for index, _ := range sessions {
		if sessions[index].Id == id {
			sessions = append(sessions[:index], sessions[index+1:]...)
			found = true
		}
	}

	if found == false {
		err = custom_errors.NewBadRequestError(errors.New("Session not found for Id: " + id))
		return err
	}

	configuration.PlainAwsAccountSessions = sessions

	err = service.UpdateConfiguration(configuration, false)
	if err != nil {
		return err
	}

	return nil
}
