package use_case

import (
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/interface/repository"
	"reflect"
)

type GcpCredentialsApplier struct {
	Keychain   Keychain
	Repository *repository.GcloudConfigurationRepository
}

func (applier *GcpCredentialsApplier) UpdateGcpPlainSessions(oldSessions []session.GcpPlainSession, newSessions []session.GcpPlainSession) {
	oldActiveSession, newActiveSession := applier.getActiveSessions(oldSessions, newSessions)

	if oldActiveSession != nil {
		if newActiveSession == nil {
			applier.deactivateSession(oldActiveSession)
		} else if oldActiveSession.Id != newActiveSession.Id {
			applier.deactivateSession(oldActiveSession)
			applier.activateSession(newActiveSession)
		} else {
			if !reflect.DeepEqual(oldActiveSession, newActiveSession) {
				applier.activateSession(newActiveSession)
			}
		}
	} else {
		if newActiveSession != nil {
			applier.activateSession(newActiveSession)
		}
	}
}

func (applier *GcpCredentialsApplier) getActiveSessions(oldSessions []session.GcpPlainSession, newSessions []session.GcpPlainSession) (*session.GcpPlainSession, *session.GcpPlainSession) {
	var oldActiveSession *session.GcpPlainSession
	for _, oldSession := range oldSessions {
		if oldSession.Status == session.Active {
			oldActiveSession = &oldSession
			break
		}
	}
	var newActiveSession *session.GcpPlainSession
	for _, newSession := range newSessions {
		if newSession.Status == session.Active {
			newActiveSession = &newSession
			break
		}
	}
	return oldActiveSession, newActiveSession
}

func (applier *GcpCredentialsApplier) activateSession(session *session.GcpPlainSession) {
	credentials, err := applier.Keychain.GetSecret(session.CredentialsLabel)
	if err != nil {
		logging.Entry().Error(err)
		return
	}

	err = applier.Repository.WriteCredentialsToDb(session.AccountId, credentials)
	if err != nil {
		logging.Entry().Error(err)
		return
	}

	err = applier.Repository.CreateConfiguration(session.Name, session.AccountId, session.ProjectName)
	if err != nil {
		logging.Entry().Error(err)
		return
	}

	err = applier.Repository.ActivateConfiguration(session.Name)
	if err != nil {
		logging.Entry().Error(err)
		return
	}

	err = applier.Repository.WriteDefaultCredentials(credentials)
	if err != nil {
		logging.Entry().Error(err)
		return
	}
}

func (applier *GcpCredentialsApplier) deactivateSession(session *session.GcpPlainSession) {
	err := applier.Repository.RemoveDefaultCredentials()
	if err != nil {
		logging.Entry().Error(err)
		return
	}

	err = applier.Repository.DeactivateConfiguration()
	if err != nil {
		logging.Entry().Error(err)
		return
	}

	err = applier.Repository.RemoveCredentialsFromDb(session.AccountId)
	if err != nil {
		logging.Entry().Error(err)
		return
	}

	err = applier.Repository.RemoveAccessTokensFromDb(session.AccountId)
	if err != nil {
		logging.Entry().Error(err)
		return
	}

	err = applier.Repository.RemoveConfiguration(session.Name)
	if err != nil {
		logging.Entry().Error(err)
		return
	}
}
