package mock

import (
	"errors"
	"fmt"
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/http/http_error"
)

type GcpPlainSessionsFacadeMock struct {
	calls                    []string
	ExpErrorOnGetSessionById bool
	ExpErrorOnAddSession     bool
	ExpErrorOnEditSession    bool
	ExpErrorOnRemoveSession  bool
	ExpErrorOnStartSession   bool
	ExpErrorOnStopSession    bool
	ExpGetSessionById        session.GcpPlainSession
	ExpGetSessions           []session.GcpPlainSession
}

func NewGcpPlainSessionsFacadeMock() GcpPlainSessionsFacadeMock {
	return GcpPlainSessionsFacadeMock{calls: []string{}, ExpGetSessions: []session.GcpPlainSession{}}
}

func (facade *GcpPlainSessionsFacadeMock) GetCalls() []string {
	return facade.calls
}

func (facade *GcpPlainSessionsFacadeMock) GetSessions() []session.GcpPlainSession {
	facade.calls = append(facade.calls, "GetSessions()")
	return facade.ExpGetSessions

}

func (facade *GcpPlainSessionsFacadeMock) GetSessionById(sessionId string) (session.GcpPlainSession, error) {
	facade.calls = append(facade.calls, fmt.Sprintf("GetSessionById(%v)", sessionId))
	if facade.ExpErrorOnGetSessionById {
		return session.GcpPlainSession{}, http_error.NewNotFoundError(errors.New("session not found"))
	}
	return facade.ExpGetSessionById, nil
}

func (facade *GcpPlainSessionsFacadeMock) AddSession(session session.GcpPlainSession) error {
	facade.calls = append(facade.calls, fmt.Sprintf("AddSession(%v)", session.Name))
	if facade.ExpErrorOnAddSession {
		return http_error.NewConflictError(errors.New("session already exist"))
	}
	return nil
}

func (facade *GcpPlainSessionsFacadeMock) StartSession(sessionId string, startTime string) error {
	facade.calls = append(facade.calls, fmt.Sprintf("StartSession(%v, %v)", sessionId, startTime))
	if facade.ExpErrorOnStartSession {
		return http_error.NewInternalServerError(errors.New("unable to start the session"))
	}
	return nil
}

func (facade *GcpPlainSessionsFacadeMock) StopSession(sessionId string, stopTime string) error {
	facade.calls = append(facade.calls, fmt.Sprintf("StopSession(%v, %v)", sessionId, stopTime))
	if facade.ExpErrorOnStopSession {
		return http_error.NewInternalServerError(errors.New("unable to stop the session"))
	}
	return nil
}

func (facade *GcpPlainSessionsFacadeMock) RemoveSession(sessionId string) error {
	facade.calls = append(facade.calls, fmt.Sprintf("RemoveSession(%v)", sessionId))
	if facade.ExpErrorOnRemoveSession {
		return http_error.NewNotFoundError(errors.New("session not found"))
	}
	return nil
}

func (facade *GcpPlainSessionsFacadeMock) EditSession(sessionId string, name string, projectName string) error {
	facade.calls = append(facade.calls, fmt.Sprintf("EditSession(%v, %v, %v)", sessionId, name, projectName))
	if facade.ExpErrorOnEditSession {
		return http_error.NewConflictError(errors.New("unable to edit session, collision detected"))
	}

	return nil
}
