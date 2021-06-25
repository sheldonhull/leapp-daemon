package use_case

import (
	"fmt"
	"leapp_daemon/domain/region"
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/http/http_error"
	"strings"

	"github.com/google/uuid"
)

type TrustedAlibabaSessionService struct {
	Keychain Keychain
}

func (service *TrustedAlibabaSessionService) Create(parentId string, accountName string, accountNumber string, roleName string, region string) error {
	/* Questo lo controllo più avanti
	err = CheckParentExist(parentId)
	if err != nil {
		return err
	}
	*/

	sessions := session.GetTrustedAlibabaSessionsFacade().GetTrustedAlibabaSessions()

	for _, sess := range sessions {
		account := sess.Account
		if sess.ParentId == parentId && account.AccountNumber == accountNumber && account.Role.Name == roleName {
			err := http_error.NewConflictError(fmt.Errorf("a session with the same parent, account number and role name already exists"))
			return err
		}
	}

	trustedAlibabaAccount := session.TrustedAlibabaAccount{
		AccountNumber: accountNumber,
		Name:          accountName,
		Role: &session.TrustedAlibabaRole{
			Name: roleName,
			Arn:  fmt.Sprintf("acs:ram::%s:role/%s", accountNumber, roleName),
		},
		Region: region,
	}

	// TODO check uuid format
	uuidString := uuid.New().String()
	uuidString = strings.Replace(uuidString, "-", "", -1)

	sess := session.TrustedAlibabaSession{
		Id:        uuidString,
		Status:    session.NotActive,
		StartTime: "",
		ParentId:  parentId,
		Account:   &trustedAlibabaAccount,
	}

	err := session.GetTrustedAlibabaSessionsFacade().SetTrustedAlibabaSessions(append(sessions, sess))
	if err != nil {
		return err
	}

	return nil
}

func (service *TrustedAlibabaSessionService) Get(id string) (*session.TrustedAlibabaSession, error) {
	return session.GetTrustedAlibabaSessionsFacade().GetTrustedAlibabaSessionById(id)
}

func (service *TrustedAlibabaSessionService) Update(id string, parentId string, accountName string, accountNumber string, roleName string, regionName string) error {

	isRegionValid := region.IsAlibabaRegionValid(regionName)
	if !isRegionValid {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("Region " + regionName + " not valid"))
	}

	trustedAlibabaRole := session.TrustedAlibabaRole{
		Name: roleName,
		Arn:  fmt.Sprintf("acs:ram::%s:role/%s", accountNumber, roleName),
	}

	trustedAlibabaAccount := session.TrustedAlibabaAccount{
		AccountNumber: accountNumber,
		Name:          accountName,
		Role:          &trustedAlibabaRole,
		Region:        regionName,
	}

	sess := session.TrustedAlibabaSession{
		Id:     id,
		Status: session.NotActive,
		//StartTime string
		ParentId: parentId,
		Account:  &trustedAlibabaAccount,
		//Profile   string
	}

	session.GetTrustedAlibabaSessionsFacade().SetTrustedAlibabaSessionById(sess)
	return nil
}

func (service *TrustedAlibabaSessionService) Delete(id string) error {
	err := session.GetTrustedAlibabaSessionsFacade().RemoveTrustedAlibabaSession(id)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

/*
func CheckParentExist(parentId string, config *configuration.Configuration) error {
	/*
		foundId := false
		plains, err := config.GetPlainAlibabaSessions()
		if err != nil {
			return err
		}
		for _, sess := range plains {
			if sess.Id == parentId {
				foundId = true
			}
		}

		feds, err := config.GetFederatedAlibabaSessions()
		if err != nil {
			return err
		}
		for _, sess := range feds {
			if sess.Id == parentId {
				foundId = true
			}
		}

		if !foundId {
			err := http_error2.NewNotFoundError(fmt.Errorf("no plain or federated session with id %s found", parentId))
			return err
		}
*/

//	return nil
//}

func (service *TrustedAlibabaSessionService) Start(sessionId string) error {

	err := session.GetTrustedAlibabaSessionsFacade().SetTrustedAlibabaSessionStatusToPending(sessionId)
	if err != nil {
		return err
	}

	err = session.GetTrustedAlibabaSessionsFacade().SetTrustedAlibabaSessionStatusToActive(sessionId)
	if err != nil {
		return err
	}

	return nil
}

func (service *TrustedAlibabaSessionService) Stop(sessionId string) error {
	/*
		config, err := configuration.ReadConfiguration()
		if err != nil {
			return err
		}

		// Passing nil because, it will be the rotate method to check if we need the mfaToken or not
		err = session2.StopFederatedAlibabaSession(config, sessionId)
		if err != nil {
			return err
		}

		err = config.Update()
		if err != nil {
			return err
		}

		// sess, err := session.GetFederatedAlibabaSession(config, sessionId)
		err = session_token.RemoveFromIniFile("default")
		if err != nil {
			return err
		}
	*/

	return nil
}
