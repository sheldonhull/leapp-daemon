package use_case

import (
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/interface/repository"
)

type AwsCredentialsApplier struct {
	Keychain                   Keychain
	NamedProfilesFacade        NamedProfilesFacade
	AwsConfigurationRepository AwsConfigurationRepository
}

func (applier *AwsCredentialsApplier) UpdateAwsIamUserSessions(oldSessions []session.AwsIamUserSession, newSessions []session.AwsIamUserSession) {

	activeCredentials := make([]repository.AwsTempCredentials, 0)
	for _, newSession := range newSessions {
		if newSession.Status != session.Active {
			continue
		}

		namedProfile, err := applier.NamedProfilesFacade.GetNamedProfileById(newSession.NamedProfileId)
		if err != nil {
			logging.Entry().Error(err)
			return
		}

		accessKeyId, err := applier.Keychain.GetSecret(newSession.AccessKeyIdLabel)
		if err != nil {
			logging.Entry().Error(err)
			return
		}

		secretKey, err := applier.Keychain.GetSecret(newSession.SecretKeyLabel)
		if err != nil {
			logging.Entry().Error(err)
			return
		}

		sessionToken, err := applier.Keychain.GetSecret(newSession.SessionTokenLabel)
		if err != nil {
			logging.Entry().Error(err)
			return
		}

		tempCredentials := repository.AwsTempCredentials{
			ProfileName:  namedProfile.Name,
			AccessKeyId:  accessKeyId,
			SecretKey:    secretKey,
			SessionToken: sessionToken,
			Region:       newSession.Region,
		}
		activeCredentials = append(activeCredentials, tempCredentials)
	}
	err := applier.AwsConfigurationRepository.WriteCredentials(activeCredentials)
	if err != nil {
		logging.Entry().Error(err)
	}
}
