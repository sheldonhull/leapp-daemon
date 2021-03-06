package mock

import (
	"errors"
	"fmt"
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/http/http_error"
)

type GcpIamUserAccountOauthSessionsFacadeMock struct {
	calls                    []string
	ExpErrorOnGetSessionById bool
	ExpErrorOnAddSession     bool
	ExpErrorOnEditSession    bool
	ExpErrorOnRemoveSession  bool
	ExpErrorOnStartSession   bool
	ExpErrorOnStopSession    bool
	ExpGetSessionById        session.GcpIamUserAccountOauthSession
	ExpGetSessions           []session.GcpIamUserAccountOauthSession
}

func NewGcpIamUserAccountOauthSessionsFacadeMock() GcpIamUserAccountOauthSessionsFacadeMock {
	return GcpIamUserAccountOauthSessionsFacadeMock{calls: []string{}, ExpGetSessions: []session.GcpIamUserAccountOauthSession{}}
}

func (facade *GcpIamUserAccountOauthSessionsFacadeMock) GetCalls() []string {
	return facade.calls
}

func (facade *GcpIamUserAccountOauthSessionsFacadeMock) GetSessions() []session.GcpIamUserAccountOauthSession {
	facade.calls = append(facade.calls, "GetSessions()")
	return facade.ExpGetSessions

}

func (facade *GcpIamUserAccountOauthSessionsFacadeMock) GetSessionById(sessionId string) (session.GcpIamUserAccountOauthSession, error) {
	facade.calls = append(facade.calls, fmt.Sprintf("GetSessionById(%v)", sessionId))
	if facade.ExpErrorOnGetSessionById {
		return session.GcpIamUserAccountOauthSession{}, http_error.NewNotFoundError(errors.New("session not found"))
	}
	return facade.ExpGetSessionById, nil
}

func (facade *GcpIamUserAccountOauthSessionsFacadeMock) AddSession(session session.GcpIamUserAccountOauthSession) error {
	facade.calls = append(facade.calls, fmt.Sprintf("AddSession(%v)", session.Name))
	if facade.ExpErrorOnAddSession {
		return http_error.NewConflictError(errors.New("session already exist"))
	}
	return nil
}

func (facade *GcpIamUserAccountOauthSessionsFacadeMock) StartSession(sessionId string, startTime string) error {
	facade.calls = append(facade.calls, fmt.Sprintf("StartSession(%v, %v)", sessionId, startTime))
	if facade.ExpErrorOnStartSession {
		return http_error.NewInternalServerError(errors.New("unable to start the session"))
	}
	return nil
}

func (facade *GcpIamUserAccountOauthSessionsFacadeMock) StopSession(sessionId string, stopTime string) error {
	facade.calls = append(facade.calls, fmt.Sprintf("StopSession(%v, %v)", sessionId, stopTime))
	if facade.ExpErrorOnStopSession {
		return http_error.NewInternalServerError(errors.New("unable to stop the session"))
	}
	return nil
}

func (facade *GcpIamUserAccountOauthSessionsFacadeMock) RemoveSession(sessionId string) error {
	facade.calls = append(facade.calls, fmt.Sprintf("RemoveSession(%v)", sessionId))
	if facade.ExpErrorOnRemoveSession {
		return http_error.NewNotFoundError(errors.New("session not found"))
	}
	return nil
}

func (facade *GcpIamUserAccountOauthSessionsFacadeMock) EditSession(sessionId string, name string, projectName string) error {
	facade.calls = append(facade.calls, fmt.Sprintf("EditSession(%v, %v, %v)", sessionId, name, projectName))
	if facade.ExpErrorOnEditSession {
		return http_error.NewConflictError(errors.New("unable to edit session, collision detected"))
	}

	return nil
}
