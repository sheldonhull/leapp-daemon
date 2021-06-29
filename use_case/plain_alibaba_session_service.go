package use_case

import (
	"fmt"
	"leapp_daemon/domain/constant"
	"leapp_daemon/domain/named_profile"
	"leapp_daemon/domain/region"
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/http/http_error"
	"strings"

	"github.com/google/uuid"
)

type PlainAlibabaSessionService struct {
	Keychain Keychain
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

	err := session.GetPlainAlibabaSessionsFacade().AddSession(sess)
	if err != nil {
		return err
	}

	err = service.Keychain.SetSecret(alibabaAccessKeyId, sess.Id+constant.PlainAlibabaKeyIdSuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = service.Keychain.SetSecret(alibabaSecretAccessKey, sess.Id+constant.PlainAlibabaSecretAccessKeySuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (service *PlainAlibabaSessionService) Get(id string) (*session.PlainAlibabaSession, error) {
	var sess *session.PlainAlibabaSession
	sess, err := session.GetPlainAlibabaSessionsFacade().GetSessionById(id)
	return sess, err
}

func (service *PlainAlibabaSessionService) Update(id string, alias string, regionName string,
	alibabaAccessKeyId string, alibabaSecretAccessKey string, profileName string) error {

	isRegionValid := region.IsAlibabaRegionValid(regionName)
	if !isRegionValid {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("Region " + regionName + " not valid"))
	}

	oldSess, err := session.GetPlainAlibabaSessionsFacade().GetSessionById(id)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	plainAlibabaAccount := session.PlainAlibabaAccount{
		Region:         regionName,
		NamedProfileId: oldSess.Account.NamedProfileId,
	}

	sess := session.PlainAlibabaSession{
		Id:      id,
		Alias:   alias,
		Status:  session.NotActive,
		Account: &plainAlibabaAccount,
	}

	err = session.GetPlainAlibabaSessionsFacade().UpdateSession(sess)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	oldNamedProfile := named_profile.GetNamedProfilesFacade().GetNamedProfileById(oldSess.Account.NamedProfileId)
	oldNamedProfile.Name = profileName
	named_profile.GetNamedProfilesFacade().UpdateNamedProfileName(oldNamedProfile)

	err = service.Keychain.SetSecret(alibabaAccessKeyId, sess.Id+constant.PlainAlibabaKeyIdSuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = service.Keychain.SetSecret(alibabaSecretAccessKey, sess.Id+constant.PlainAlibabaSecretAccessKeySuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (service *PlainAlibabaSessionService) Delete(sessionId string) error {

	err := session.GetPlainAlibabaSessionsFacade().RemoveSession(sessionId)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (service *PlainAlibabaSessionService) Start(sessionId string) error {

	err := session.GetPlainAlibabaSessionsFacade().SetStatusToPending(sessionId)
	if err != nil {
		return err
	}

	err = session.GetPlainAlibabaSessionsFacade().SetStatusToActive(sessionId)
	if err != nil {
		return err
	}

	return nil
}

func (service *PlainAlibabaSessionService) Stop(sessionId string) error {

	err := session.GetPlainAlibabaSessionsFacade().SetStatusToInactive(sessionId)
	if err != nil {
		return err
	}

	return nil
}
