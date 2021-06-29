package use_case

import (
	"fmt"
	"leapp_daemon/domain/region"
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/http/http_error"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/google/uuid"
)

type TrustedAlibabaSessionService struct {
	Keychain Keychain
}

func (service *TrustedAlibabaSessionService) Create(parentId string, accountName string, accountNumber string, roleName string, region string) error {
	parentSession, err := GetAlibabaParentById(parentId)
	if err != nil {
		return err
	}

	sessions := session.GetTrustedAlibabaSessionsFacade().GetTrustedAlibabaSessions()

	for _, sess := range sessions {
		account := sess.Account
		if sess.ParentSession.GetId() == parentId && account.AccountNumber == accountNumber && account.Role.Name == roleName {
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
		Id:            uuidString,
		Status:        session.NotActive,
		StartTime:     "",
		ParentSession: parentSession,
		Account:       &trustedAlibabaAccount,
	}

	err = session.GetTrustedAlibabaSessionsFacade().SetTrustedAlibabaSessions(append(sessions, sess))
	if err != nil {
		return err
	}

	return nil
}

func (service *TrustedAlibabaSessionService) Get(id string) (*session.TrustedAlibabaSession, error) {
	return session.GetTrustedAlibabaSessionsFacade().GetTrustedAlibabaSessionById(id)
}

func (service *TrustedAlibabaSessionService) Update(id string, parentId string, accountName string, accountNumber string, roleName string, regionName string) error {
	parentSession, err := GetAlibabaParentById(parentId)
	if err != nil {
		return err
	}

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
		ParentSession: parentSession,
		Account:       &trustedAlibabaAccount,
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

func (service *TrustedAlibabaSessionService) Start(sessionId string) error {
	// call AssumeRole API
	sess, err := session.GetTrustedAlibabaSessionsFacade().GetTrustedAlibabaSessionById(sessionId)
	if err != nil {
		return err
	}
	region := sess.Account.Region
	label := sess.ParentSession.GetId() + "-" + sess.ParentSession.GetTypeString() + "-alibaba-session-access-key-id"
	accessKeyId, err := service.Keychain.GetSecret(label)
	if err != nil {
		return err
	}
	label = sess.ParentSession.GetId() + "-" + sess.ParentSession.GetTypeString() + "-alibaba-session-secret-access-key"
	accessKeySecret, err := service.Keychain.GetSecret(label)
	if err != nil {
		return err
	}

	var client *sts.Client
	if sess.ParentSession.GetTypeString() == "plain" {
		client, err = sts.NewClientWithAccessKey(region, accessKeyId, accessKeySecret)
		if err != nil {
			return err
		}
	} else {
		label = sess.ParentSession.GetId() + "-" + sess.ParentSession.GetTypeString() + "-alibaba-session-sts-token"
		stsToken, err := service.Keychain.GetSecret(label)
		if err != nil {
			return err
		}
		client, err = sts.NewClientWithStsToken(region, accessKeyId, accessKeySecret, stsToken)
		if err != nil {
			return err
		}
	}

	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"
	request.RoleArn = sess.Account.Role.Arn
	request.RoleSessionName = "leapp" // TODO: find better way
	response, err := client.AssumeRole(request)
	if err != nil {
		return err
	}

	// saves credentials into keychain
	err = service.Keychain.SetSecret(response.Credentials.AccessKeyId, sess.Id+"-trusted-alibaba-session-access-key-id")
	if err != nil {
		return http_error.NewInternalServerError(err)
	}
	err = service.Keychain.SetSecret(response.Credentials.AccessKeySecret, sess.Id+"-trusted-alibaba-session-secret-access-key")
	if err != nil {
		return http_error.NewInternalServerError(err)
	}
	err = service.Keychain.SetSecret(response.Credentials.SecurityToken, sess.Id+"-trusted-alibaba-session-sts-token")
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = session.GetTrustedAlibabaSessionsFacade().SetTrustedAlibabaSessionStatusToPending(sessionId)
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
	err := session.GetTrustedAlibabaSessionsFacade().SetTrustedAlibabaSessionStatusToInactive(sessionId)
	if err != nil {
		return err
	}

	return nil
}

func GetAlibabaParentById(parentId string) (session.ParentSession, error) {
	plain, err := session.GetPlainAlibabaSessionsFacade().GetSessionById(parentId)
	if err != nil {
		federated, err := session.GetFederatedAlibabaSessionsFacade().GetSessionById(parentId)
		if err != nil {
			return nil, http_error.NewNotFoundError(fmt.Errorf("no plain or federated session with id %s found", parentId))
		}
		return federated, nil
	}
	return plain, nil
}
