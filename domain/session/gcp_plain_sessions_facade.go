package session

import (
	"fmt"
	"leapp_daemon/infrastructure/http/http_error"
	"sync"
)

var sessionsFacadeSingleton *gcpPlainSessionsFacade
var sessionsFacadeLock sync.Mutex
var sessionsLock sync.Mutex

type GcpPlainSessionsObserver interface {
	UpdateGcpPlainSessions(oldSessions []GcpPlainSession, newSessions []GcpPlainSession) error
}

type gcpPlainSessionsFacade struct {
	gcpPlainSessions []GcpPlainSession
	observers        []GcpPlainSessionsObserver
}

func GetGcpPlainSessionFacade() *gcpPlainSessionsFacade {
	sessionsFacadeLock.Lock()
	defer sessionsFacadeLock.Unlock()

	if sessionsFacadeSingleton == nil {
		sessionsFacadeSingleton = &gcpPlainSessionsFacade{
			gcpPlainSessions: make([]GcpPlainSession, 0),
		}
	}
	return sessionsFacadeSingleton
}

func (facade *gcpPlainSessionsFacade) Subscribe(observer GcpPlainSessionsObserver) {
	facade.observers = append(facade.observers, observer)
}

func (facade *gcpPlainSessionsFacade) GetSessions() []GcpPlainSession {
	return facade.gcpPlainSessions
}

func (facade *gcpPlainSessionsFacade) GetSessionById(id string) (GcpPlainSession, error) {
	for _, session := range facade.GetSessions() {
		if session.Id == id {
			return session, nil
		}
	}

	return GcpPlainSession{}, http_error.NewNotFoundError(fmt.Errorf("gcp plain session with id %s not found", id))
}

func (facade *gcpPlainSessionsFacade) SetSessions(sessions []GcpPlainSession) {
	sessionsLock.Lock()
	defer sessionsLock.Unlock()

	facade.gcpPlainSessions = sessions
}

func (facade *gcpPlainSessionsFacade) AddSession(session GcpPlainSession) error {
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

	err := facade.updateState(newSessions)
	if err != nil {
		return err
	}

	return nil
}

func (facade *gcpPlainSessionsFacade) RemoveSession(id string) error {
	sessionsLock.Lock()
	defer sessionsLock.Unlock()

	currentSessions := facade.GetSessions()
	newSessions := make([]GcpPlainSession, 0)

	for _, session := range currentSessions {
		if session.Id != id {
			newSessions = append(newSessions, session)
		}
	}

	if len(currentSessions) == len(newSessions) {
		return http_error.NewNotFoundError(fmt.Errorf("plain gcp session with id %s not found", id))
	}

	err := facade.updateState(newSessions)
	if err != nil {
		return err
	}

	return nil
}

func (facade *gcpPlainSessionsFacade) SetSessionStatus(id string, status Status) error {
	sessionsLock.Lock()
	defer sessionsLock.Unlock()

	currentSessions := facade.GetSessions()
	newSessions := make([]GcpPlainSession, 0)

	sessionFound := false
	for _, session := range currentSessions {
		if session.Id == id {
			session.Status = status
			sessionFound = true
		}
		newSessions = append(newSessions, session)
	}

	if !sessionFound {
		return http_error.NewNotFoundError(fmt.Errorf("plain gcp session with id %s not found", id))
	}
	return facade.updateState(newSessions)
}

func (facade *gcpPlainSessionsFacade) updateState(newSessions []GcpPlainSession) error {
	oldSessions := facade.GetSessions()
	facade.gcpPlainSessions = newSessions

	for _, observer := range facade.observers {
		err := observer.UpdateGcpPlainSessions(oldSessions, newSessions)
		if err != nil {
			return err
		}
	}

	return nil
}
