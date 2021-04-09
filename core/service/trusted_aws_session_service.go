package service

import (
	"fmt"
	"github.com/google/uuid"
	"leapp_daemon/core/configuration"
	"leapp_daemon/core/session"
	"leapp_daemon/custom_error"
	"strings"
)

func CreateTrustedAwsSession(parentId string, accountName string, accountNumber string, roleName string, region string) error {
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return err
	}

	err = CheckParentExist(parentId, config)
	if err != nil {
		return err
	}

	sessions, err := config.GetTrustedAwsSessions()
	if err != nil {
		return err
	}

	for _, sess := range sessions {
		account := sess.Account
		if sess.ParentId == parentId && account.AccountNumber == accountNumber && account.Role.Name == roleName {
			err := custom_error.NewConflictError(fmt.Errorf("a session with the same parent, account number and role name already exists"))
			return err
		}
	}

	trustedAwsAccount := session.TrustedAwsAccount{
		AccountNumber: accountNumber,
		Name:          accountName,
		Role: &session.TrustedAwsRole{
			Name: roleName,
			Arn:  fmt.Sprintf("arn:aws:iam::%s:role/%s", accountNumber, roleName),
		},
		Region: region,
	}

	// TODO check uuid format
	uuidString := uuid.New().String()
	uuidString = strings.Replace(uuidString, "-", "", -1)

	sess := session.TrustedAwsSession{
		Id:        uuidString,
		Status:    session.NotActive,
		StartTime: "",
		ParentId:  parentId,
		Account:   &trustedAwsAccount,
	}

	err = config.SetTrustedAwsSessions(append(sessions, &sess))
	if err != nil {
		return err
	}

	err = config.Update()
	if err != nil {
		return err
	}

	return nil
}

func GetTrustedAwsSession(id string) (*session.TrustedAwsSession, error) {
	var sess *session.TrustedAwsSession

	config, err := configuration.ReadConfiguration()
	if err != nil {
		return sess, err
	}

	sessions, err := config.GetTrustedAwsSessions()
	if err != nil {
		return nil, err
	}

	for _, s := range sessions {
		if s.Id == id {
			return s, nil
		}
	}

	return nil, custom_error.NewNotFoundError(fmt.Errorf("no session found with id %s", id))
}

func UpdateTrustedAwsSession(id string, parentId string, accountName string, accountNumber string, roleName string, region string) error {
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return err
	}

	foundId := false
	sessions, err := config.GetTrustedAwsSessions()
	if err != nil {
		return err
	}

	for _, s := range sessions {
		if s.Id == id {
			foundId = true

			if parentId != "" {
				err := CheckParentExist(parentId, config)
				if err != nil {
					return err
				}
				s.ParentId = parentId
			}

			if accountNumber != "" {
				s.Account.AccountNumber = accountNumber
			}

			if accountName != "" {
				s.Account.Name = accountName
			}
			if roleName != "" {
				s.Account.Role.Name = roleName
			}
			if region != "" {
				s.Account.Region = region
			}

			if accountNumber != "" || roleName != "" {
				s.Account.Role.Arn = fmt.Sprintf("arn:aws:iam::%s:role/%s", s.Account.AccountNumber, s.Account.Role.Name)
			}

			break
		}
	}

	if !foundId {
		return custom_error.NewNotFoundError(fmt.Errorf("no session found with id %s", id))
	}

	err = config.Update()
	if err != nil {
		return err
	}

	return nil
}

func DeleteTrustedAwsSession(id string) error {
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return err
	}

	sessions, err := config.GetTrustedAwsSessions()
	if err != nil {
		return err
	}

	found := false
	for index := range sessions {
		if sessions[index].Id == id {
			sessions = append(sessions[:index], sessions[index+1:]...)
			found = true
			break
		}
	}

	if found == false {
		err = custom_error.NewNotFoundError(fmt.Errorf("trusted aws session with id %s not found", id))
		return err
	}

	err = config.SetTrustedAwsSessions(sessions)
	if err != nil {
		return err
	}

	err = config.Update()
	if err != nil {
		return err
	}

	return nil
}

func CheckParentExist(parentId string, config *configuration.Configuration) error {
	foundId := false
	plains, err := config.GetPlainAwsSessions()
	if err != nil {
		return err
	}
	for _, sess := range plains {
		if sess.Id == parentId {
			foundId = true
		}
	}

	feds, err := config.GetFederatedAwsSessions()
	if err != nil {
		return err
	}
	for _, sess := range feds {
		if sess.Id == parentId {
			foundId = true
		}
	}

	if !foundId {
		err := custom_error.NewNotFoundError(fmt.Errorf("no plain or federated session with id %s found", parentId))
		return err
	}

	return nil
}
