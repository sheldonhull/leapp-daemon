package session

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"leapp_daemon/core/aws_session_token"
	"leapp_daemon/core/configuration"
	"leapp_daemon/core/constant"
	"leapp_daemon/core/websocket"
	"leapp_daemon/custom_error"
	"leapp_daemon/logging"
	"log"
	"strings"
	"time"
)

type Status int

const (
	NotActive Status = iota
	Pending
	Active
)

type PlainAwsSession struct {
	Id        string
	Status    Status
	StartTime string
	Account   *PlainAwsAccount
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

type AwsSessionToken struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
}

func(sess *PlainAwsSession) IsMfaRequired() (bool, error) {
	return sess.Account.MfaDevice != "", nil
}

func(sess *PlainAwsSession) IsRotationIntervalExpired() (bool, error) {
	startTime, _ := time.Parse(time.RFC3339, sess.StartTime)
	secondsPassedFromStart := time.Now().Sub(startTime).Seconds()
	return int64(secondsPassedFromStart) > constant.RotationIntervalInSeconds, nil
}

func(sess *PlainAwsSession) RotateCredentials(mfaToken *string) error {
	if sess.Status == Active {
		isRotationIntervalExpired, err := sess.IsRotationIntervalExpired()
		if err != nil {
			return err
		}

		if isRotationIntervalExpired {
			println("Rotating session with id", sess.Id)
			err = sess.RotatePlainAwsSessionCredentials(mfaToken)
			if err != nil {
				return nil
			}
		}
	}
	return nil
}

func(sess *PlainAwsSession) RotatePlainAwsSessionCredentials(mfaToken *string) error {
	doSessionTokenExist, err := aws_session_token.DoExist(sess.Account.Name)
	if err != nil {
		return err
	}

	if doSessionTokenExist {
		isSessionTokenExpired, err := aws_session_token.IsExpired(sess.Account.Name)
		if err != nil {
			return err
		}

		if isSessionTokenExpired {
			logging.Entry().Error("Plain AWS session token no more valid")

			isMfaTokenRequired, err := sess.IsMfaRequired()
			if err != nil {
				return err
			}

			if isMfaTokenRequired && mfaToken == nil {
				sess.Status = Pending

				err = sendMfaRequestMessage(sess)
				if err != nil {
					return err
				}

				return nil
			}

			credentials, err := aws_session_token.Generate(sess.Account.Name, sess.Account.Region, sess.Account.MfaDevice, mfaToken)
			if err != nil {
				return err
			}

			err = aws_session_token.SaveInKeychain(sess.Account.Name, credentials)
			if err != nil {
				return err
			}

			err = aws_session_token.SaveInIniFile(*credentials.AccessKeyId, *credentials.SecretAccessKey,
				*credentials.SessionToken, sess.Account.Region, "default")
			if err != nil {
				return err
			}

			sess.Status = Active
			sess.StartTime = time.Now().Format(time.RFC3339)
		} else {
			logging.Entry().Error("Plain AWS session token still valid")

			data, err := sess.unmarshallSessionToken()
			if err != nil {
				return err
			}

			err = aws_session_token.SaveInIniFile(data.AccessKeyId,
				data.SecretAccessKey,
				data.SessionToken,
				sess.Account.Region,
				"default")
			if err != nil {
				return err
			}

			sess.Status = Active
			sess.StartTime = time.Now().Format(time.RFC3339)
		}

		return nil
	} else {
		credentials, err := aws_session_token.Generate(sess.Account.Name, sess.Account.Region, sess.Account.MfaDevice, mfaToken)
		if err != nil {
			return err
		}

		err = aws_session_token.SaveInKeychain(sess.Account.Name, credentials)
		if err != nil {
			return err
		}

		err = aws_session_token.SaveInIniFile(*credentials.AccessKeyId, *credentials.SecretAccessKey,
			*credentials.SessionToken, sess.Account.Region, "default")
		if err != nil {
			return err
		}

		sess.Status = Active
		sess.StartTime = time.Now().Format(time.RFC3339)

		return nil
	}
}

func (sess *PlainAwsSession) unmarshallSessionToken() (AwsSessionToken, error) {
	sessionTokenJson, _, err := aws_session_token.Get(sess.Account.Name)

	var data AwsSessionToken

	err = json.Unmarshal([]byte(sessionTokenJson), &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func CreatePlainAwsSession(sessionContainer Container, name string, accountNumber string,
	region string, user string, awsAccessKeyId string, awsSecretAccessKey string, mfaDevice string) error {

	sessions, err := sessionContainer.GetPlainAwsSessions()
	if err != nil {
		return err
	}

	for _, sess := range sessions {
		account := sess.Account
		if account.AccountNumber == accountNumber && account.User == user {
			err := custom_error.NewBadRequestError(errors.New("a session with the same account number and user is already present"))
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

	sess := PlainAwsSession{
		Id:        uuidString,
		Status:    NotActive,
		StartTime: "",
		Account:   &plainAwsAccount,
	}

	err = sessionContainer.SetPlainAwsSessions(append(sessions, &sess))
	if err != nil {
		return err
	}

	return nil
}

func DeletePlainAwsSession(id string) error {
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return err
	}

	sessions := config.PlainAwsSessions

	found := false
	for index := range sessions {
		if sessions[index].Id == id {
			sessions = append(sessions[:index], sessions[index+1:]...)
			found = true
			break
		}
	}

	if found == false {
		err = custom_error.NewBadRequestError(errors.New("Plain AWS session not found for Id: " + id))
		return err
	}

	config.PlainAwsSessions = sessions

	err = configuration.UpdateConfiguration(config, false)
	if err != nil {
		return err
	}

	return nil
}

func GetById(config *configuration.Configuration, id string) (*PlainAwsSession, error) {
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

func GetPlainAwsSession(id string) (*PlainAwsSession, error) {
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return nil, err
	}

	sess, err := GetById(config, id)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func IsMfaRequiredForPlainAwsSession(id string) (bool, error) {
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return false, err
	}

	sess, err := GetById(config, id)
	if err != nil {
		return false, err
	}

	return sess.IsMfaRequired()
}

func StartPlainAwsSession(id string, mfaToken *string) error {
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return err
	}

	sess, err := GetById(config, id)
	if err != nil {
		return err
	}

	println("Rotating session with id", sess.Id)
	err = sess.RotatePlainAwsSessionCredentials(mfaToken)
	if err != nil { return err }

	return nil
}

func UpdatePlainAwsSession(id string, name string, accountNumber string, region string,
	user string, awsAccessKeyId string, awsSecretAccessKey string, mfaDevice string) error {

	config, err := configuration.ReadConfiguration()
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

	err = configuration.UpdateConfiguration(config, false)
	if err != nil {
		return err
	}

	return nil
}

func sendMfaRequestMessage(sess *PlainAwsSession) error {
	messageData := websocket.MfaTokenRequestData{
		SessionId: sess.Id,
	}

	messageDataJson, err := json.Marshal(messageData)
	if err != nil {
		return err
	}

	msg := websocket.Message{
		MessageType: websocket.MfaTokenRequest,
		Data:        string(messageDataJson),
	}

	err = websocket.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Println("sent message", string(messageDataJson))

	return nil
}