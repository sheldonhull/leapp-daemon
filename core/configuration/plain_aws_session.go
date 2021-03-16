package configuration

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"leapp_daemon/constant"
	"leapp_daemon/core/session_token"
	"leapp_daemon/custom_error"
	"leapp_daemon/logging"
	"strings"
	"time"
)

type PlainAwsSession struct {
	Id           string
	Active       bool
	Loading      bool
	StartTime    string
	Account      *PlainAwsAccount
}

type PlainAwsAccount struct {
	AccountNumber       string
	Name                string
	Region              string
	User                string
	AwsAccessKeyId      string
	AwsSecretAccessKey  string
	MfaDevice           string
}

func(sess *PlainAwsSession) Rotate(configuration *Configuration, mfaToken *string) error {
	if sess.Active {
		isRotationIntervalExpired, err := sess.IsRotationIntervalExpired()
		if err != nil {
			return err
		}

		if isRotationIntervalExpired {
			isMfaTokenRequired, err := sess.IsMfaRequired()
			if err != nil { return nil }

			if isMfaTokenRequired {
				// TODO: need to implement a way to ask for token.
				//  We will probably need WebSocket for this. We will exploit the Hub to wait for client response.
			} else {
				println("Rotating session with id", sess.Id)
				err = sess.rotatePlainAwsSessionCredentials(configuration, nil)
				if err != nil { return nil }
			}
		}
	}

	return nil
}

func(sess *PlainAwsSession) IsRotationIntervalExpired() (bool, error) {
	startTime, _ := time.Parse(time.RFC3339, sess.StartTime)
	secondsPassedFromStart := time.Now().Sub(startTime).Seconds()
	return int64(secondsPassedFromStart) > constant.RotationIntervalInSeconds, nil
}

func(sess *PlainAwsSession) IsMfaRequired() (bool, error) {
	return sess.Account.MfaDevice != "", nil
}

func(sess *PlainAwsSession) rotatePlainAwsSessionCredentials(config *Configuration, mfaToken *string) error {
	doSessionTokenExist, err := session_token.DoExist(sess.Account.Name)
	if err != nil {
		return err
	}

	if doSessionTokenExist {
		isSessionTokenExpired, err := session_token.IsExpired(sess.Account.Name)
		if err != nil {
			return err
		}

		if isSessionTokenExpired {
			logging.Entry().Error("Plain AWS session token no more valid")

			credentials, err := session_token.Generate(sess.Account.Name, sess.Account.Region, sess.Account.MfaDevice, mfaToken)
			if err != nil {
				return err
			}

			err = session_token.SaveInKeychain(sess.Account.Name, credentials)
			if err != nil {
				return err
			}

			err = session_token.SaveInIniFile(*credentials.AccessKeyId, *credentials.SecretAccessKey,
				*credentials.SessionToken, sess.Account.Region, "default")
			if err != nil {
				return err
			}

			sess.Active = true
			sess.StartTime = time.Now().Format(time.RFC3339)

			err = UpdateConfiguration(config, false)
			if err != nil {
				return err
			}
		} else {
			logging.Entry().Error("Plain AWS session token still valid")

			sessionTokenJson, _, err := session_token.Get(sess.Account.Name)

			data := struct {
				AccessKeyId string
				SecretAccessKey string
				SessionToken string
			} {}

			err = json.Unmarshal([]byte(sessionTokenJson), &data)
			if err != nil { return err }

			err = session_token.SaveInIniFile(data.AccessKeyId, data.SecretAccessKey, data.SessionToken,
				sess.Account.Region, "default")
			if err != nil { return err }

			sess.Active = true
			sess.StartTime = time.Now().Format(time.RFC3339)

			err = UpdateConfiguration(config, false)
			if err != nil {
				return err
			}
		}
	} else {
		credentials, err := session_token.Generate(sess.Account.Name, sess.Account.Region, sess.Account.MfaDevice, mfaToken)
		if err != nil {
			return err
		}

		err = session_token.SaveInKeychain(sess.Account.Name, credentials)
		if err != nil {
			return err
		}

		err = session_token.SaveInIniFile(*credentials.AccessKeyId, *credentials.SecretAccessKey,
			*credentials.SessionToken, sess.Account.Region, "default")
		if err != nil {
			return err
		}

		sess.Active = true
		sess.StartTime = time.Now().Format(time.RFC3339)

		err = UpdateConfiguration(config, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreatePlainAwsSession(name string, accountNumber string, region string, user string,
	awsAccessKeyId string, awsSecretAccessKey string, mfaDevice string) error {
	config, err := ReadConfiguration()
	if err != nil {
		return err
	}

	sessions := config.PlainAwsSessions

	for _, session := range sessions {
		account := session.Account
		if account.AccountNumber == accountNumber && account.User == user {
			err = custom_error.NewBadRequestError(errors.New("an account with the same account number and user is already present"))
			return err
		}
	}

	plainAwsAccount := PlainAwsAccount{
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

	session := PlainAwsSession{
		Id:           uuidString,
		Active:       false,
		Loading:      false,
		StartTime: "",
		Account:      &plainAwsAccount,
	}

	config.PlainAwsSessions = append(config.PlainAwsSessions, &session)

	err = UpdateConfiguration(config, false)
	if err != nil {
		return err
	}

	return nil
}

func GetById(config *Configuration, id string) (*PlainAwsSession, error) {
	sessions := config.PlainAwsSessions
	var sess *PlainAwsSession

	for index := range sessions {
		if sessions[index].Id == id {
			sess = sessions[index]
			return sess, nil
		}
	}

	err := custom_error.NewBadRequestError(errors.New("Plain AWS session not found for Id: " + id))
	return nil, err
}

func List(query string) ([]*PlainAwsSession, error) {
	config, err := ReadConfiguration()
	if err != nil {
		return nil, err
	}

	filteredList := make([]*PlainAwsSession, 0)

	if query == "" {
		return append(filteredList, config.PlainAwsSessions...), nil
	} else {
		for _, session := range config.PlainAwsSessions {
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

func UpdatePlainAwsSession(id string, name string, accountNumber string, region string,
	user string, awsAccessKeyId string, awsSecretAccessKey string, mfaDevice string) error {

	config, err := ReadConfiguration()
	if err != nil {
		return err
	}

	sessions := config.PlainAwsSessions

	found := false
	for index := range sessions {
		if sessions[index].Id == id {
			sessions[index].Account = &PlainAwsAccount{
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
		err = custom_error.NewBadRequestError(errors.New("Plain AWS session not found for Id: " + id))
		return err
	}

	config.PlainAwsSessions = sessions

	err = UpdateConfiguration(config, false)
	if err != nil {
		return err
	}

	return nil
}

func DeletePlainAwsSession(id string) error {
	config, err := ReadConfiguration()
	if err != nil {
		return err
	}

	sessions := config.PlainAwsSessions

	found := false
	for index := range sessions {
		if sessions[index].Id == id {
			sessions = append(sessions[:index], sessions[index+1:]...)
			found = true
		}
	}

	if found == false {
		err = custom_error.NewBadRequestError(errors.New("Plain AWS session not found for Id: " + id))
		return err
	}

	config.PlainAwsSessions = sessions

	err = UpdateConfiguration(config, false)
	if err != nil {
		return err
	}

	return nil
}