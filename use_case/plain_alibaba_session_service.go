package use_case

import (
	"fmt"
	"leapp_daemon/domain/named_profile"
	"leapp_daemon/domain/region"
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/http/http_error"
	"strings"

	"github.com/google/uuid"
)

type AlibabaKeychain interface {
	DoesSecretExist(label string) (bool, error)
	GetSecret(label string) (string, error)
	SetSecret(secret string, label string) error
}

type PlainAlibabaSessionService struct {
	Keychain AlibabaKeychain
}

func (service *PlainAlibabaSessionService) Create(alias string, alibabaAccessKeyId string, alibabaSecretAccessKey string, regionName string, profileName string) error {

	namedProfile := named_profile.GetNamedProfilesFacade().GetNamedProfileByName(profileName)

	if namedProfile == nil {
		// TODO: extract UUID generation logic
		uuidString := uuid.New().String()
		uuidString = strings.Replace(uuidString, "-", "", -1)

		namedProfile = &named_profile.NamedProfile{
			Id:   uuidString,
			Name: profileName,
		}

		err := named_profile.GetNamedProfilesFacade().AddNamedProfile(*namedProfile)
		if err != nil {
			return err
		}
	}

	isRegionValid := region.IsAlibabaRegionValid(regionName)
	if !isRegionValid {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("Region " + regionName + " not valid"))
	}

	plainAlibabaAccount := session.PlainAlibabaAccount{
		Region:         regionName,
		NamedProfileId: namedProfile.Id,
	}

	// TODO: extract UUID generation logic
	uuidString := uuid.New().String()
	uuidString = strings.Replace(uuidString, "-", "", -1)

	sess := session.PlainAlibabaSession{
		Id:      uuidString,
		Alias:   alias,
		Status:  session.NotActive,
		Account: &plainAlibabaAccount,
	}

	err := session.GetPlainAlibabaSessionsFacade().AddPlainAlibabaSession(sess)
	if err != nil {
		return err
	}

	err = service.Keychain.SetSecret(alibabaAccessKeyId, sess.Id+"-plain-alibaba-session-access-key-id")
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = service.Keychain.SetSecret(alibabaSecretAccessKey, sess.Id+"-plain-alibaba-session-secret-access-key")
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (service *PlainAlibabaSessionService) GetPlainAlibabaSession(id string) (*session.PlainAlibabaSession, error) {
	var sess *session.PlainAlibabaSession
	sess, err := session.GetPlainAlibabaSessionsFacade().GetPlainAlibabaSessionById(id)
	return sess, err
}

func (service *PlainAlibabaSessionService) UpdatePlainAlibabaSession(sessionId string, name string, region string, user string,
	alibabaAccessKeyId string, alibabaSecretAccessKey string, profile string) error {

	/*
		config, err := configuration.ReadConfiguration()
		if err != nil {
			return err
		}

		err = session.UpdatePlainAlibabaSession(config, sessionId, name, region, user, alibabaAccessKeyId, alibabaSecretAccessKey, mfaDevice, profile)
		if err != nil {
			return err
		}

		err = config.Update()
		if err != nil {
			return err
		}
	*/

	return nil
}

func DeletePlainAlibabaSession(sessionId string) error {
	/*
		config, err := configuration.ReadConfiguration()
		if err != nil {
			return err
		}

		err = session.DeletePlainAlibabaSession(config, sessionId)
		if err != nil {
			return err
		}

		err = config.Update()
		if err != nil {
			return err
		}
	*/

	return nil
}

func (service *PlainAlibabaSessionService) StartPlainAlibabaSession(sessionId string) error {

	err := session.GetPlainAlibabaSessionsFacade().SetPlainAlibabaSessionStatusToPending(sessionId)
	if err != nil {
		return err
	}

	err = session.GetPlainAlibabaSessionsFacade().SetPlainAlibabaSessionStatusToActive(sessionId)
	if err != nil {
		return err
	}

	return nil
}

func StopPlainAlibabaSession(sessionId string) error {
	/*
			config, err := configuration.ReadConfiguration()
			if err != nil {
				return err
			}

			// Passing nil because, it will be the rotate method to check if we need the mfaToken or not
			err = session.StopPlainAlibabaSession(config, sessionId)
			if err != nil {
				return err
			}

			err = config.Update()
			if err != nil {
				return err
			}

		  // sess, err := session.GetPlainAlibabaSession(config, sessionId)
			err = session_token.RemoveFromIniFile("default")
			if err != nil {
				return err
			}
	*/

	return nil
}
