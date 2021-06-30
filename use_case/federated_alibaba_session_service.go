package use_case

import (
	"fmt"
	"leapp_daemon/domain/constant"
	"leapp_daemon/domain/named_profile"
	"leapp_daemon/domain/region"
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/http/http_error"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/google/uuid"
)

type FederatedAlibabaSessionService struct {
	Keychain Keychain
}

// TODO: mettere da qualche parte questa funzione
func SAMLAuth(region string, idpArn string, roleArn string, assertion string) (key string, secret string, token string, err error) {
	// I'm using this since NewClient() method returns a panic saying literally "not support yet"
	// This method actually never use the credentials so I placed 2 placeholders
	client, _ := sts.NewClientWithAccessKey(region, "", "")

	request := sts.CreateAssumeRoleWithSAMLRequest()
	request.Scheme = "https"
	request.SAMLProviderArn = idpArn
	request.RoleArn = roleArn
	request.SAMLAssertion = assertion
	response, err := client.AssumeRoleWithSAML(request)
	if err != nil {
		return "", "", "", err
	}
	key = response.Credentials.AccessKeyId
	secret = response.Credentials.AccessKeySecret
	token = response.Credentials.SecurityToken
	return
}

func (service *FederatedAlibabaSessionService) Create(name string, accountNumber string, roleName string, roleArn string,
	idpArn string, regionName string, ssoUrl string, profileName string) error {

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

	federatedAlibabaRole := session.FederatedAlibabaRole{
		Name: roleName,
		Arn:  roleArn,
	}

	federatedAlibabaAccount := session.FederatedAlibabaAccount{
		AccountNumber: accountNumber,
		Name:          name,
		Role:          &federatedAlibabaRole,
		IdpArn:        idpArn,
		Region:        regionName,
		/*SsoUrl:        ssoUrl,*/
		NamedProfileId: namedProfile.Id,
	}

	// TODO: extract UUID generation logic
	uuidString := uuid.New().String()
	uuidString = strings.Replace(uuidString, "-", "", -1)

	sess := session.FederatedAlibabaSession{
		Id:      uuidString,
		Status:  session.NotActive,
		Account: &federatedAlibabaAccount,
		Profile: profileName,
	}

	err := session.GetFederatedAlibabaSessionsFacade().AddSession(sess)
	if err != nil {
		return err
	}

	alibabaAccessKeyId, alibabaSecretAccessKey, alibabaStsToken, err := SAMLAuth(regionName, idpArn, roleArn, ssoUrl)
	if err != nil {
		return err
	}

	err = service.Keychain.SetSecret(alibabaAccessKeyId, sess.Id+constant.FederatedAlibabaKeyIdSuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = service.Keychain.SetSecret(alibabaSecretAccessKey, sess.Id+constant.FederatedAlibabaSecretAccessKeySuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = service.Keychain.SetSecret(alibabaStsToken, sess.Id+constant.FederatedAlibabaStsTokenSuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (service *FederatedAlibabaSessionService) Get(id string) (*session.FederatedAlibabaSession, error) {
	return session.GetFederatedAlibabaSessionsFacade().GetSessionById(id)
}

func (service *FederatedAlibabaSessionService) Update(id string, name string, accountNumber string, roleName string, roleArn string,
	idpArn string, regionName string, ssoUrl string, profileName string) error {

	isRegionValid := region.IsAlibabaRegionValid(regionName)
	if !isRegionValid {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("Region " + regionName + " not valid"))
	}

	oldSess, err := session.GetFederatedAlibabaSessionsFacade().GetSessionById(id)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	federatedAlibabaRole := session.FederatedAlibabaRole{
		Name: roleName,
		Arn:  roleArn,
	}

	federatedAlibabaAccount := session.FederatedAlibabaAccount{
		AccountNumber: accountNumber,
		Name:          name,
		Role:          &federatedAlibabaRole,
		IdpArn:        idpArn,
		Region:        regionName,
		/*SsoUrl:        ssoUrl,*/
		NamedProfileId: oldSess.Account.NamedProfileId,
	}

	sess := session.FederatedAlibabaSession{
		Id:      id,
		Status:  session.NotActive,
		Account: &federatedAlibabaAccount,
		Profile: profileName,
	}

	err = session.GetFederatedAlibabaSessionsFacade().UpdateSession(sess)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	oldNamedProfile := named_profile.GetNamedProfilesFacade().GetNamedProfileById(oldSess.Account.NamedProfileId)
	oldNamedProfile.Name = profileName
	named_profile.GetNamedProfilesFacade().UpdateNamedProfileName(oldNamedProfile)

	alibabaAccessKeyId, alibabaSecretAccessKey, alibabaStsToken, err := SAMLAuth(regionName, idpArn, roleArn, ssoUrl)
	if err != nil {
		return err
	}

	err = service.Keychain.SetSecret(alibabaAccessKeyId, id+constant.FederatedAlibabaKeyIdSuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = service.Keychain.SetSecret(alibabaSecretAccessKey, id+constant.FederatedAlibabaSecretAccessKeySuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = service.Keychain.SetSecret(alibabaStsToken, id+constant.FederatedAlibabaStsTokenSuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (service *FederatedAlibabaSessionService) Delete(sessionId string) error {
	sess, err := session.GetFederatedAlibabaSessionsFacade().GetSessionById(sessionId)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	if sess.Status != session.NotActive {
		err = service.Stop(sessionId)
		if err != nil {
			return err
		}
	}

	oldNamedProfile := named_profile.GetNamedProfilesFacade().GetNamedProfileById(sess.Account.NamedProfileId)
	named_profile.GetNamedProfilesFacade().DeleteNamedProfile(oldNamedProfile.Id)

	err = session.GetFederatedAlibabaSessionsFacade().RemoveSession(sessionId)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = service.Keychain.DeleteSecret(sessionId + constant.FederatedAlibabaKeyIdSuffix)
	if err != nil {
		return err
	}

	err = service.Keychain.DeleteSecret(sessionId + constant.FederatedAlibabaSecretAccessKeySuffix)
	if err != nil {
		return err
	}

	err = service.Keychain.DeleteSecret(sessionId + constant.FederatedAlibabaStsTokenSuffix)
	if err != nil {
		return err
	}
	return nil
}

func (service *FederatedAlibabaSessionService) Start(sessionId string) error {

	err := session.GetFederatedAlibabaSessionsFacade().SetStatusToPending(sessionId)
	if err != nil {
		return err
	}

	err = session.GetFederatedAlibabaSessionsFacade().SetStatusToActive(sessionId)
	if err != nil {
		return err
	}

	return nil
}

func (service *FederatedAlibabaSessionService) Stop(sessionId string) error {

	err := session.GetFederatedAlibabaSessionsFacade().SetStatusToInactive(sessionId)
	if err != nil {
		return err
	}

	return nil
}
