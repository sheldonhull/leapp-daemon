package session

import (
  "fmt"
  "leapp_daemon/infrastructure/http/http_error"
  "sync"
)

var facadeSingleton *facade
var facadeLock sync.Mutex
var plainAwsSessionsLock sync.Mutex

type PlainAwsSessionsObserver interface {
  UpdatePlainAwsSessions(plainAwsSessions []PlainAwsSession) error
}

type facade struct {
  plainAwsSessions []PlainAwsSession
  observers        []PlainAwsSessionsObserver
}

func GetPlainAwsSessionsFacade() *facade {
  facadeLock.Lock()
  defer facadeLock.Unlock()

  if facadeSingleton == nil {
    facadeSingleton = &facade {
      plainAwsSessions: make([]PlainAwsSession, 0),
    }
  }

  return facadeSingleton
}

func(fac *facade) Subscribe(observer PlainAwsSessionsObserver) {
  fac.observers = append(fac.observers, observer)
}

func(fac *facade) SetPlainAwsSessions(plainAwsSessions []PlainAwsSession) {
  fac.plainAwsSessions = plainAwsSessions
}

func(fac *facade) AddPlainAwsSession(plainAwsSession PlainAwsSession) error {
  plainAwsSessionsLock.Lock()
  defer plainAwsSessionsLock.Unlock()

  sessions := fac.plainAwsSessions

  for _, sess := range sessions {
    if plainAwsSession.Id == sess.Id {
      return http_error.NewConflictError(fmt.Errorf("a PlainAwsSession with id " + plainAwsSession.Id +
        " is already present"))
    }
  }

  sessions = append(sessions, plainAwsSession)
  fac.plainAwsSessions = sessions

  err := fac.notifyObservers()
  if err != nil {
    return err
  }

  return nil
}

func(fac *facade) notifyObservers() error {
  for _, observer := range fac.observers {
    err := observer.UpdatePlainAwsSessions(fac.plainAwsSessions)
    if err != nil {
      return err
    }
  }

  return nil
}
