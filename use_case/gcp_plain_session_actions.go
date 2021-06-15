package use_case

import (
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/http/http_error"
)

type GcpPlainSessionActions struct {
	GcpApi                GcpApi
	Environment           Environment
	Keychain              Keychain
	GcpPlainSessionFacade GcpPlainSessionsFacade
	NamedProfilesActions  NamedProfilesActionsInterface
}

func (actions *GcpPlainSessionActions) GetSession(sessionId string) (session.GcpPlainSession, error) {
	return actions.GcpPlainSessionFacade.GetSessionById(sessionId)
}

func (actions *GcpPlainSessionActions) GetOAuthUrl() (string, error) {
	return actions.GcpApi.GetOauthUrl()
}

func (actions *GcpPlainSessionActions) CreateSession(name string, accountId string, projectName string, profileName string,
	oauthCode string) error {

	namedProfile, err := actions.NamedProfilesActions.GetOrCreateNamedProfile(profileName)
	if err != nil {
		return err
	}

	newSessionId := actions.Environment.GenerateUuid()
	credentialsLabel := newSessionId + "-gcp-plain-session-credentials"

	gcpSession := session.GcpPlainSession{
		Id:               newSessionId,
		Name:             name,
		AccountId:        accountId,
		ProjectName:      projectName,
		NamedProfileId:   namedProfile.Id,
		CredentialsLabel: credentialsLabel,
		Status:           session.NotActive,
		StartTime:        "",
		LastStopTime:     "",
	}

	token, err := actions.GcpApi.GetOauthToken(oauthCode)
	if err != nil {
		return err
	}

	credentials := actions.GcpApi.GetCredentials(token)

	err = actions.Keychain.SetSecret(credentials, credentialsLabel)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return actions.GcpPlainSessionFacade.AddSession(gcpSession)
}

func (actions *GcpPlainSessionActions) StartSession(sessionId string) error {

	facade := actions.GcpPlainSessionFacade
	for _, currentSession := range facade.GetSessions() {
		if currentSession.Status != session.NotActive && currentSession.Id != sessionId {
			err := facade.SetSessionStatus(currentSession.Id, session.NotActive)
			if err != nil {
				return err
			}
		}
	}
	return facade.SetSessionStatus(sessionId, session.Active)
}

func (actions *GcpPlainSessionActions) StopSession(sessionId string) error {

	return actions.GcpPlainSessionFacade.SetSessionStatus(sessionId, session.NotActive)
}

func (actions *GcpPlainSessionActions) DeleteSession(sessionId string) error {
	facade := actions.GcpPlainSessionFacade

	sessionToDelete, err := facade.GetSessionById(sessionId)
	if err != nil {
		return err
	}

	_ = actions.Keychain.DeleteSecret(sessionToDelete.CredentialsLabel)
	return facade.RemoveSession(sessionId)
}

func (actions *GcpPlainSessionActions) EditSession(sessionId string, name string, projectName string, profileName string) error {
	sessionsFacade := actions.GcpPlainSessionFacade

	namedProfile, err := actions.NamedProfilesActions.GetOrCreateNamedProfile(profileName)
	if err != nil {
		return err
	}

	return sessionsFacade.EditSession(sessionId, name, projectName, namedProfile.Id)
}
