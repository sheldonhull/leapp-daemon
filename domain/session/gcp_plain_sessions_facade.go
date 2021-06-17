package session

import (
	"fmt"
	"leapp_daemon/infrastructure/http/http_error"
	"sync"
)

var sessionsLock sync.Mutex

type GcpPlainSessionsObserver interface {
	UpdateGcpPlainSessions(oldSessions []GcpPlainSession, newSessions []GcpPlainSession)
}

type GcpPlainSessionsFacade struct {
	gcpPlainSessions []GcpPlainSession
	observers        []GcpPlainSessionsObserver
}

func NewGcpPlainSessionsFacade() *GcpPlainSessionsFacade {
	return &GcpPlainSessionsFacade{
		gcpPlainSessions: make([]GcpPlainSession, 0),
	}
}

func (facade *GcpPlainSessionsFacade) Subscribe(observer GcpPlainSessionsObserver) {
	facade.observers = append(facade.observers, observer)
}

func (facade *GcpPlainSessionsFacade) GetSessions() []GcpPlainSession {
	return facade.gcpPlainSessions
}

func (facade *GcpPlainSessionsFacade) GetSessionById(sessionId string) (GcpPlainSession, error) {
	for _, session := range facade.GetSessions() {
		if session.Id == sessionId {
			return session, nil
		}
	}

	return GcpPlainSession{}, http_error.NewNotFoundError(fmt.Errorf("gcp plain session with id %s not found", sessionId))
}

func (facade *GcpPlainSessionsFacade) SetSessions(sessions []GcpPlainSession) {
	sessionsLock.Lock()
	defer sessionsLock.Unlock()

	facade.gcpPlainSessions = sessions
}

func (facade *GcpPlainSessionsFacade) AddSession(session GcpPlainSession) error {
	sessionsLock.Lock()
	defer sessionsLock.Unlock()

	currentSessions := facade.GetSessions()

	for _, sess := range currentSessions {
		if session.Id == sess.Id {
			return http_error.NewConflictError(fmt.Errorf("a session with id %v is already present", session.Id))
		}

		if session.Name == sess.Name {
			return http_error.NewConflictError(fmt.Errorf("a session named %v is already present", sess.Name))
		}
	}

	newSessions := append(currentSessions, session)

	facade.updateState(newSessions)
	return nil
}

func (facade *GcpPlainSessionsFacade) RemoveSession(sessionId string) error {
	sessionsLock.Lock()
	defer sessionsLock.Unlock()

	currentSessions := facade.GetSessions()
	newSessions := make([]GcpPlainSession, 0)

	for _, session := range currentSessions {
		if session.Id != sessionId {
			newSessions = append(newSessions, session)
		}
	}

	if len(currentSessions) == len(newSessions) {
		return http_error.NewNotFoundError(fmt.Errorf("plain gcp session with id %s not found", sessionId))
	}

	facade.updateState(newSessions)
	return nil
}

func (facade *GcpPlainSessionsFacade) SetSessionStatus(sessionId string, status Status) error {
	sessionsLock.Lock()
	defer sessionsLock.Unlock()

	sessionToUpdate, err := facade.GetSessionById(sessionId)
	if err != nil {
		return err
	}

	sessionToUpdate.Status = status
	return facade.replaceSession(sessionId, sessionToUpdate)
}

func (facade *GcpPlainSessionsFacade) EditSession(sessionId string, name string, projectName string, profileId string) error {
	sessionsLock.Lock()
	defer sessionsLock.Unlock()

	sessionToEdit, err := facade.GetSessionById(sessionId)
	if err != nil {
		return err
	}

	sessionToEdit.Name = name
	sessionToEdit.ProjectName = projectName
	sessionToEdit.NamedProfileId = profileId

	return facade.replaceSession(sessionId, sessionToEdit)
}

func (facade *GcpPlainSessionsFacade) replaceSession(sessionId string, newSession GcpPlainSession) error {
	newSessions := make([]GcpPlainSession, 0)
	for _, session := range facade.GetSessions() {
		if session.Id == sessionId {
			newSessions = append(newSessions, newSession)
		} else {
			newSessions = append(newSessions, session)
		}
	}

	facade.updateState(newSessions)
	return nil
}

func (facade *GcpPlainSessionsFacade) updateState(newSessions []GcpPlainSession) {
	oldSessions := facade.GetSessions()
	facade.gcpPlainSessions = newSessions

	for _, observer := range facade.observers {
		observer.UpdateGcpPlainSessions(oldSessions, newSessions)
	}
}
