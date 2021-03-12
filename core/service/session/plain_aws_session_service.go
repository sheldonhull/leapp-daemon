package session

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"leapp_daemon/core/model"
	"leapp_daemon/core/service"
	errors2 "leapp_daemon/shared/custom_error"
	"leapp_daemon/shared/logging"
	"strings"
)

func GetPlainAwsSession(id string) (*model.PlainAwsSession, error) {
	configuration, err := service.ReadConfiguration()
	if err != nil {
		return nil, err
	}

	sessions := configuration.PlainAwsSessions
	for index, _ := range sessions {
		if sessions[index].Id == id {
			return &sessions[index], nil
		}
	}

	return nil, errors2.NewBadRequestError(errors.New("No session found with id:" + id))
}

func ListPlainAwsSession(query string) ([]model.PlainAwsSession, error) {
	configuration, err := service.ReadConfiguration()
	if err != nil {
		return nil, err
	}

	filteredList := make([]model.PlainAwsSession, 0)

	if query == "" {
		return append(filteredList, configuration.PlainAwsSessions...), nil
	} else {
		for _, session := range configuration.PlainAwsSessions {
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

func CreatePlainAwsSession(name string, accountNumber string, region string, user string,
	                       awsAccessKeyId string, awsSecretAccessKey string, mfaDevice string) error {
	configuration, err := service.ReadConfiguration()
	if err != nil {
		return err
	}

	sessions := configuration.PlainAwsSessions

	for _, session := range sessions {
		account := session.Account
		if account.AccountNumber == accountNumber && account.User == user {
			err = errors2.NewBadRequestError(errors.New("an account with the same account number and user is already present"))
			return err
		}
	}

	plainAwsAccount := model.PlainAwsAccount{
		AccountNumber: accountNumber,
		Name:          name,
		Region:        region,
		User:          user,
		AwsAccessKeyId: awsAccessKeyId,
		AwsSecretAccessKey: awsSecretAccessKey,
		MfaDevice:     mfaDevice,
	}

	uuidString := uuid.New().String()
	uuidString = strings.Replace(uuidString, "-", "", -1)

	session := model.PlainAwsSession{
		Id:           uuidString,
		Active:       false,
		Loading:      false,
		StartTime: "",
		Account:      plainAwsAccount,
	}

	configuration.PlainAwsSessions = append(configuration.PlainAwsSessions, session)

	err = service.UpdateConfiguration(configuration, false)
	if err != nil {
		return err
	}

	return nil
}

func EditPlainAwsSession(id string, name string, accountNumber string, region string,
	                     user string, awsAccessKeyId string, awsSecretAccessKey string, mfaDevice string) error {

	configuration, err := service.ReadConfiguration()
	if err != nil {
		return err
	}

	sessions := configuration.PlainAwsSessions

	found := false
	for index, _ := range sessions {
		if sessions[index].Id == id {
			sessions[index].Account = model.PlainAwsAccount{
				AccountNumber: accountNumber,
				Name:          name,
				Region:        region,
				User:          user,
				AwsAccessKeyId: awsAccessKeyId,
				AwsSecretAccessKey: awsSecretAccessKey,
				MfaDevice:     mfaDevice,
			}
			found = true
		}
	}

	if found == false {
		err = errors2.NewBadRequestError(errors.New("Session not found for Id: " + id))
		return err
	}

	configuration.PlainAwsSessions = sessions

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

	sessions := configuration.PlainAwsSessions

	found := false
	for index, _ := range sessions {
		if sessions[index].Id == id {
			sessions = append(sessions[:index], sessions[index+1:]...)
			found = true
		}
	}

	if found == false {
		err = errors2.NewBadRequestError(errors.New("Session not found for Id: " + id))
		return err
	}

	configuration.PlainAwsSessions = sessions

	err = service.UpdateConfiguration(configuration, false)
	if err != nil {
		return err
	}

	return nil
}

func IsMfaRequiredForPlainAwsSession(id string) (bool, error) {
	configuration, err := service.ReadConfiguration()
	if err != nil {
		return false, err
	}

	sess, err := getPlainAwsSessionById(configuration, id)
	if err != nil {
		return false, err
	}

	return sess.Account.MfaDevice != "", nil
}

func StartPlainAwsSession(id string, mfaToken string) error {
	configuration, err := service.ReadConfiguration()
	if err != nil {
		return err
	}

	sess, err := getPlainAwsSessionById(configuration, id)
	if err != nil {
		return err
	}

	err = startPlainAwsSession(sess, configuration, &mfaToken)
	if err != nil {
		return err
	}

	return nil
}

func startPlainAwsSession(sess *model.PlainAwsSession, configuration *model.Configuration, mfaToken *string) error {
	doSessionTokenExist, err := service.DoSessionTokenExist(sess.Account.Name)
	if err != nil {
		return err
	}

	if doSessionTokenExist {
		isSessionTokenExpired, err := service.IsSessionTokenExpired(sess.Account.Name)
		if err != nil {
			return err
		}

		if isSessionTokenExpired {
			// generate creds
			logging.Info("session token no more valid")

			credentials, err := service.GenerateSessionToken(sess, mfaToken)
			if err != nil {
				return err
			}

			err = service.SaveSessionTokenInKeychain(sess.Account.Name, credentials)
			if err != nil {
				return err
			}

			err = service.SaveSessionTokenInIniFile(*credentials.AccessKeyId, *credentials.SecretAccessKey,
				*credentials.SessionToken, sess.Account.Region, "default")
			if err != nil {
				return err
			}

			sess.Active = true
			err = service.UpdateConfiguration(configuration, false)
			if err != nil {
				return err
			}
		} else {
			// re-use creds
			logging.Info("session token still valid")

			accessKeyId, secretAccessKey, err := service.GetAccessKeys(sess.Account.Name)
			if err != nil {
				return err
			}
			sessionToken, _, err := service.GetSessionToken(sess.Account.Name)

			err = service.SaveSessionTokenInIniFile(accessKeyId, secretAccessKey, sessionToken, sess.Account.Region,
				"default")
			if err != nil {
				return err
			}
		}
	} else {
		credentials, err := service.GenerateSessionToken(sess, mfaToken)
		if err != nil {
			return err
		}

		err = service.SaveSessionTokenInKeychain(sess.Account.Name, credentials)
		if err != nil {
			return err
		}

		err = service.SaveSessionTokenInIniFile(*credentials.AccessKeyId, *credentials.SecretAccessKey,
			*credentials.SessionToken, sess.Account.Region, "default")
		if err != nil {
			return err
		}

		sess.Active = true
		err = service.UpdateConfiguration(configuration, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func getPlainAwsSessionById(configuration *model.Configuration, id string) (*model.PlainAwsSession, error) {
	sessions := configuration.PlainAwsSessions
	var sess *model.PlainAwsSession

	for index, _ := range sessions {
		if sessions[index].Id == id {
			sess = &sessions[index]
		}
	}

	if sess == nil {
		err := errors2.NewBadRequestError(errors.New("Session not found for Id: " + id))
		return nil, err
	}

	return sess, nil
}
