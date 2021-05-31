package session

import (
  "fmt"
  "leapp_daemon/infrastructure/http/http_error"
  "sync"
)

var plainAwsSessionsFacadeSingleton *plainAwsSessionsFacade
var plainAwsSessionsFacadeLock sync.Mutex
var plainAwsSessionsLock sync.Mutex

type PlainAwsSessionsObserver interface {
  UpdatePlainAwsSessions(oldPlainAwsSessions []PlainAwsSession, newPlainAwsSessions []PlainAwsSession) error
}

type plainAwsSessionsFacade struct {
  plainAwsSessions []PlainAwsSession
  observers        []PlainAwsSessionsObserver
}

func GetPlainAwsSessionsFacade() *plainAwsSessionsFacade {
  plainAwsSessionsFacadeLock.Lock()
  defer plainAwsSessionsFacadeLock.Unlock()

  if plainAwsSessionsFacadeSingleton == nil {
    plainAwsSessionsFacadeSingleton = &plainAwsSessionsFacade{
      plainAwsSessions: make([]PlainAwsSession, 0),
    }
  }

  return plainAwsSessionsFacadeSingleton
}

func(fac *plainAwsSessionsFacade) Subscribe(observer PlainAwsSessionsObserver) {
  fac.observers = append(fac.observers, observer)
}

func(fac *plainAwsSessionsFacade) GetPlainAwsSessions() []PlainAwsSession {
  return fac.plainAwsSessions
}

func(fac *plainAwsSessionsFacade) SetPlainAwsSessions(plainAwsSessions []PlainAwsSession) {
  fac.plainAwsSessions = plainAwsSessions
}

func(fac *plainAwsSessionsFacade) AddPlainAwsSession(plainAwsSession PlainAwsSession) error {
  plainAwsSessionsLock.Lock()
  defer plainAwsSessionsLock.Unlock()

  sessions := fac.plainAwsSessions

  for _, sess := range sessions {
    if plainAwsSession.Id == sess.Id {
      return http_error.NewConflictError(fmt.Errorf("a PlainAwsSession with id " + plainAwsSession.Id +
        " is already present"))
    }

    if plainAwsSession.Alias == sess.Alias {
      return http_error.NewUnprocessableEntityError(fmt.Errorf("a session with the same alias " +
        "is already present"))
    }
  }

  fac.plainAwsSessions = append(sessions, plainAwsSession)

  err := fac.updateState(fac.plainAwsSessions)
  if err != nil {
    return err
  }

  return nil
}

func(fac *plainAwsSessionsFacade) RemovePlainAwsSession(id string) error {
  plainAwsSessionsLock.Lock()
  defer plainAwsSessionsLock.Unlock()

  sessions := fac.plainAwsSessions

  for i, sess := range sessions {
    if sess.Id == id {
      sessions = append(sessions[:i], sessions[i+1:]...)
      break
    }
  }

  if len(fac.plainAwsSessions) == len(sessions) {
    return http_error.NewNotFoundError(fmt.Errorf("plain aws session with id %s not found", id))
  }

  fac.plainAwsSessions = sessions

  err := fac.updateState(fac.plainAwsSessions)
  if err != nil {
    return err
  }

  return nil
}

func(fac *plainAwsSessionsFacade) GetPlainAwsSessionById(id string) (*PlainAwsSession, error) {
  for _, plainAwsSession := range fac.plainAwsSessions {
    if plainAwsSession.Id == id {
      return &plainAwsSession, nil
    }
  }
  return nil, http_error.NewNotFoundError(fmt.Errorf("plain aws session with id %s not found", id))
}

func(fac *plainAwsSessionsFacade) SetPlainAwsSessionStatusToPending(id string) error {
  plainAwsSession, err := fac.GetPlainAwsSessionById(id)
  if err != nil {
    return err
  }

  if !(plainAwsSession.Status == NotActive) {
    return http_error.NewUnprocessableEntityError(fmt.Errorf("plain aws session with id " + id + "cannot be started because it's in pending or active state"))
  }

  for i, session := range fac.plainAwsSessions {
    if session.Id == id {
      fac.plainAwsSessions[i].Status = Pending
    }
  }

  err = fac.updateState(fac.plainAwsSessions)
  if err != nil {
    return err
  }

  return nil
}

func(fac *plainAwsSessionsFacade) SetPlainAwsSessionStatusToActive(id string) error {
  plainAwsSession, err := fac.GetPlainAwsSessionById(id)
  if err != nil {
    return err
  }

  if !(plainAwsSession.Status == Pending) {
    return http_error.NewUnprocessableEntityError(fmt.Errorf("plain aws session with id " + id + "cannot be started because it's not in pending state"))
  }

  for i, session := range fac.plainAwsSessions {
    if session.Id == id {
      fac.plainAwsSessions[i].Status = Active
    }
  }

  err = fac.updateState(fac.plainAwsSessions)
  if err != nil {
    return err
  }

  return nil
}

func(fac *plainAwsSessionsFacade) updateState(newState []PlainAwsSession) error {
  oldPlainAwsSessions := fac.plainAwsSessions
  for _, observer := range fac.observers {
    err := observer.UpdatePlainAwsSessions(oldPlainAwsSessions, newState)
    if err != nil {
      return err
    }
  }
  return nil
}
