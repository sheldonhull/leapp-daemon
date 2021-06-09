package use_case

import (
	"github.com/google/uuid"
	"leapp_daemon/domain/named_profile"
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/http/http_error"
	"strings"
)

type GcpPlainSessionActions struct {
	Keychain Keychain
	GcpApi   GcpApi
}

func (actions *GcpPlainSessionActions) CreateSession(name string, accountId string, projectName string, profileName string,
	oauthCode string) error {

	// TODO: check this logic
	if profileName == "" {
		profileName = "default"
	}

	namedProfile := named_profile.GetNamedProfilesFacade().GetNamedProfileByName(profileName)

	// TODO: move to named_profiles_actions
	if namedProfile == nil {
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

	// TODO: move to external logic
	newSessionId := strings.Replace(uuid.New().String(), "-", "", -1)

	gcpSession := session.GcpPlainSession{
		Id:             newSessionId,
		Name:           name,
		AccountId:      accountId,
		ProjectName:    projectName,
		NamedProfileId: namedProfile.Id,
		Status:         session.NotActive,
		StartTime:      "",
		LastStopTime:   "",
	}

	token, err := actions.GcpApi.GetOauthToken(oauthCode)
	if err != nil {
		return err
	}

	credentials := actions.GcpApi.GetCredentials(token)

	// TODO: Move to a dedicated GCP Keychain interface
	err = actions.Keychain.SetSecret(credentials, newSessionId+"-gcp-plain-session-credentials")
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return session.GetGcpPlainSessionFacade().AddSession(gcpSession)
}

func (actions *GcpPlainSessionActions) GetSession(sessionId string) (session.GcpPlainSession, error) {
	return session.GetGcpPlainSessionFacade().GetSessionById(sessionId)
}

func (actions *GcpPlainSessionActions) GetOAuthUrl() (string, error) {
	return actions.GcpApi.GetOauthUrl()
}
