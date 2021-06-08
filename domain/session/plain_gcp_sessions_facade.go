package session

import (
	"fmt"
	"leapp_daemon/infrastructure/http/http_error"
	"sync"
)

var plainGcpSessionsFacadeSingleton *plainGcpSessionsFacade
var plainGcpSessionsFacadeLock sync.Mutex
var plainGcpSessionsLock sync.Mutex

type PlainGcpSessionsObserver interface {
	UpdatePlainGcpSessions(oldPlainGcpSessions []PlainGcpSession, newPlainGcpSessions []PlainGcpSession) error
}

type plainGcpSessionsFacade struct {
	plainGcpSessions []PlainGcpSession
	observers        []PlainGcpSessionsObserver
}

func GetPlainGcpSessionsFacade() *plainGcpSessionsFacade {
	plainGcpSessionsFacadeLock.Lock()
	defer plainGcpSessionsFacadeLock.Unlock()

	if plainGcpSessionsFacadeSingleton == nil {
		plainGcpSessionsFacadeSingleton = &plainGcpSessionsFacade{
			plainGcpSessions: make([]PlainGcpSession, 0),
		}
	}
	return plainGcpSessionsFacadeSingleton
}

func (facade *plainGcpSessionsFacade) Subscribe(observer PlainGcpSessionsObserver) {
	facade.observers = append(facade.observers, observer)
}

func (facade *plainGcpSessionsFacade) GetSessions() []PlainGcpSession {
	return facade.plainGcpSessions
}

// TODO: blast it if not needed
func (facade *plainGcpSessionsFacade) GetSessionById(id string) (PlainGcpSession, error) {
	for _, plainGcpSession := range facade.GetSessions() {
		if plainGcpSession.Id == id {
			return plainGcpSession, nil
		}
	}
	return PlainGcpSession{}, http_error.NewNotFoundError(fmt.Errorf("plain gcp session with id %s not found", id))
}

func (facade *plainGcpSessionsFacade) SetSessions(plainGcpSessions []PlainGcpSession) {
	plainGcpSessionsLock.Lock()
	defer plainGcpSessionsLock.Unlock()

	facade.plainGcpSessions = plainGcpSessions
}

func (facade *plainGcpSessionsFacade) AddSession(plainGcpSession PlainGcpSession) error {
	plainGcpSessionsLock.Lock()
	defer plainGcpSessionsLock.Unlock()

	currentSessions := facade.GetSessions()

	for _, sess := range currentSessions {
		if plainGcpSession.Id == sess.Id {
			return http_error.NewConflictError(fmt.Errorf("a PlainGcpSession with id " + plainGcpSession.Id +
				" is already present"))
		}

		if plainGcpSession.Alias == sess.Alias {
			return http_error.NewUnprocessableEntityError(fmt.Errorf("a session with the same alias " +
				"is already present"))
		}
	}

	newSessions := append(currentSessions, plainGcpSession)

	err := facade.updateState(newSessions)
	if err != nil {
		return err
	}

	return nil
}

func (facade *plainGcpSessionsFacade) RemoveSession(id string) error {
	plainGcpSessionsLock.Lock()
	defer plainGcpSessionsLock.Unlock()

	currentSessions := facade.GetSessions()
	newSessions := make([]PlainGcpSession, 0)

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

func (facade *plainGcpSessionsFacade) SetSessionStatus(id string, status Status) error {
	plainGcpSessionsLock.Lock()
	defer plainGcpSessionsLock.Unlock()

	currentSessions := facade.GetSessions()
	newSessions := make([]PlainGcpSession, 0)

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

func (facade *plainGcpSessionsFacade) updateState(newSessions []PlainGcpSession) error {
	oldSessions := facade.GetSessions()
	facade.plainGcpSessions = newSessions

	for _, observer := range facade.observers {
		err := observer.UpdatePlainGcpSessions(oldSessions, newSessions)
		if err != nil {
			return err
		}
	}

	return nil
}
