package session

import (
	"fmt"
	"leapp_daemon/infrastructure/http/http_error"
	"sync"
)

var plainAlibabaSessionsFacadeSingleton *plainAlibabaSessionsFacade
var plainAlibabaSessionsFacadeLock sync.Mutex
var plainAlibabaSessionsLock sync.Mutex

type PlainAlibabaSessionsObserver interface {
	UpdatePlainAlibabaSessions(oldPlainAlibabaSessions []PlainAlibabaSession, newPlainAlibabaSessions []PlainAlibabaSession) error
}

type plainAlibabaSessionsFacade struct {
	plainAlibabaSessions []PlainAlibabaSession
	observers            []PlainAlibabaSessionsObserver
}

func GetPlainAlibabaSessionsFacade() *plainAlibabaSessionsFacade {
	plainAlibabaSessionsFacadeLock.Lock()
	defer plainAlibabaSessionsFacadeLock.Unlock()

	if plainAlibabaSessionsFacadeSingleton == nil {
		plainAlibabaSessionsFacadeSingleton = &plainAlibabaSessionsFacade{
			plainAlibabaSessions: make([]PlainAlibabaSession, 0),
		}
	}

	return plainAlibabaSessionsFacadeSingleton
}

func (fac *plainAlibabaSessionsFacade) Subscribe(observer PlainAlibabaSessionsObserver) {
	fac.observers = append(fac.observers, observer)
}

func (fac *plainAlibabaSessionsFacade) GetSessions() []PlainAlibabaSession {
	return fac.plainAlibabaSessions
}

func (fac *plainAlibabaSessionsFacade) SetSessions(newPlainAlibabaSessions []PlainAlibabaSession) error {
	fac.plainAlibabaSessions = newPlainAlibabaSessions

	err := fac.updateState(newPlainAlibabaSessions)
	if err != nil {
		return err
	}
	return nil
}

func (fac *plainAlibabaSessionsFacade) UpdateSession(newSession PlainAlibabaSession) error {
	allSessions := fac.GetSessions()
	for i, plainAlibabaSession := range allSessions {
		if plainAlibabaSession.Id == newSession.Id {
			allSessions[i] = newSession
		}
	}
	err := fac.SetSessions(allSessions)
	return err
}

func (fac *plainAlibabaSessionsFacade) AddSession(plainAlibabaSession PlainAlibabaSession) error {
	plainAlibabaSessionsLock.Lock()
	defer plainAlibabaSessionsLock.Unlock()

	oldPlainAlibabaSessions := fac.GetSessions()
	newPlainAlibabaSessions := make([]PlainAlibabaSession, 0)

	for i := range oldPlainAlibabaSessions {
		newPlainAlibabaSession := oldPlainAlibabaSessions[i]
		newPlainAlibabaSessionAccount := *oldPlainAlibabaSessions[i].Account
		newPlainAlibabaSession.Account = &newPlainAlibabaSessionAccount
		newPlainAlibabaSessions = append(newPlainAlibabaSessions, newPlainAlibabaSession)
	}

	for _, sess := range newPlainAlibabaSessions {
		if plainAlibabaSession.Id == sess.Id {
			return http_error.NewConflictError(fmt.Errorf("a PlainAlibabaSession with id " + plainAlibabaSession.Id +
				" is already present"))
		}

		if plainAlibabaSession.Alias == sess.Alias {
			return http_error.NewUnprocessableEntityError(fmt.Errorf("a session with the same alias " +
				"is already present"))
		}
	}

	newPlainAlibabaSessions = append(newPlainAlibabaSessions, plainAlibabaSession)

	err := fac.updateState(newPlainAlibabaSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *plainAlibabaSessionsFacade) RemoveSession(id string) error {
	plainAlibabaSessionsLock.Lock()
	defer plainAlibabaSessionsLock.Unlock()

	oldPlainAlibabaSessions := fac.GetSessions()
	newPlainAlibabaSessions := make([]PlainAlibabaSession, 0)

	for i := range oldPlainAlibabaSessions {
		newPlainAlibabaSession := oldPlainAlibabaSessions[i]
		newPlainAlibabaSessionAccount := *oldPlainAlibabaSessions[i].Account
		newPlainAlibabaSession.Account = &newPlainAlibabaSessionAccount
		newPlainAlibabaSessions = append(newPlainAlibabaSessions, newPlainAlibabaSession)
	}

	for i, sess := range newPlainAlibabaSessions {
		if sess.Id == id {
			newPlainAlibabaSessions = append(newPlainAlibabaSessions[:i], newPlainAlibabaSessions[i+1:]...)
			break
		}
	}

	if len(fac.GetSessions()) == len(newPlainAlibabaSessions) {
		return http_error.NewNotFoundError(fmt.Errorf("plain Alibaba session with id %s not found", id))
	}

	err := fac.updateState(newPlainAlibabaSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *plainAlibabaSessionsFacade) GetSessionById(id string) (*PlainAlibabaSession, error) {
	for _, plainAlibabaSession := range fac.GetSessions() {
		if plainAlibabaSession.Id == id {
			return &plainAlibabaSession, nil
		}
	}
	return nil, http_error.NewNotFoundError(fmt.Errorf("plain Alibaba session with id %s not found", id))
}

func (fac *plainAlibabaSessionsFacade) SetStatusToPending(id string) error {
	plainAlibabaSessionsLock.Lock()
	defer plainAlibabaSessionsLock.Unlock()

	plainAlibabaSession, err := fac.GetSessionById(id)
	if err != nil {
		return err
	}

	if !(plainAlibabaSession.Status == NotActive) {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("plain Alibaba session with id " + id + "cannot be started because it's in pending or active state"))
	}

	oldPlainAlibabaSessions := fac.GetSessions()
	newPlainAlibabaSessions := make([]PlainAlibabaSession, 0)

	for i := range oldPlainAlibabaSessions {
		newPlainAlibabaSession := oldPlainAlibabaSessions[i]
		newPlainAlibabaSessionAccount := *oldPlainAlibabaSessions[i].Account
		newPlainAlibabaSession.Account = &newPlainAlibabaSessionAccount
		newPlainAlibabaSessions = append(newPlainAlibabaSessions, newPlainAlibabaSession)
	}

	for i, session := range newPlainAlibabaSessions {
		if session.Id == id {
			newPlainAlibabaSessions[i].Status = Pending
		}
	}

	err = fac.updateState(newPlainAlibabaSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *plainAlibabaSessionsFacade) SetStatusToActive(id string) error {
	plainAlibabaSessionsLock.Lock()
	defer plainAlibabaSessionsLock.Unlock()

	plainAlibabaSession, err := fac.GetSessionById(id)
	if err != nil {
		return err
	}

	if !(plainAlibabaSession.Status == Pending) {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("plain Alibaba session with id " + id + "cannot be started because it's not in pending state"))
	}

	oldPlainAlibabaSessions := fac.GetSessions()
	newPlainAlibabaSessions := make([]PlainAlibabaSession, 0)

	for i := range oldPlainAlibabaSessions {
		newPlainAlibabaSession := oldPlainAlibabaSessions[i]
		newPlainAlibabaSessionAccount := *oldPlainAlibabaSessions[i].Account
		newPlainAlibabaSession.Account = &newPlainAlibabaSessionAccount
		newPlainAlibabaSessions = append(newPlainAlibabaSessions, newPlainAlibabaSession)
	}

	for i, session := range newPlainAlibabaSessions {
		if session.Id == id {
			newPlainAlibabaSessions[i].Status = Active
		}
	}

	err = fac.updateState(newPlainAlibabaSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *plainAlibabaSessionsFacade) SetStatusToInactive(id string) error {
	plainAlibabaSessionsLock.Lock()
	defer plainAlibabaSessionsLock.Unlock()

	plainAlibabaSession, err := fac.GetSessionById(id)
	if err != nil {
		return err
	}
	if plainAlibabaSession.Status != Active {
		fmt.Println(plainAlibabaSession.Status)
		return http_error.NewUnprocessableEntityError(fmt.Errorf("plain Alibaba session with id " + id + "cannot be started because it's not in active state"))
	}

	oldPlainAlibabaSessions := fac.GetSessions()
	newPlainAlibabaSessions := make([]PlainAlibabaSession, 0)

	for i := range oldPlainAlibabaSessions {
		newPlainAlibabaSession := oldPlainAlibabaSessions[i]
		newPlainAlibabaSessionAccount := *oldPlainAlibabaSessions[i].Account
		newPlainAlibabaSession.Account = &newPlainAlibabaSessionAccount
		newPlainAlibabaSessions = append(newPlainAlibabaSessions, newPlainAlibabaSession)
	}

	for i, session := range newPlainAlibabaSessions {
		if session.Id == id {
			newPlainAlibabaSessions[i].Status = NotActive
		}
	}

	err = fac.updateState(newPlainAlibabaSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *plainAlibabaSessionsFacade) updateState(newState []PlainAlibabaSession) error {
	oldPlainAlibabaSessions := fac.GetSessions()
	fac.plainAlibabaSessions = newState

	for _, observer := range fac.observers {
		err := observer.UpdatePlainAlibabaSessions(oldPlainAlibabaSessions, newState)
		if err != nil {
			return err
		}
	}

	return nil
}