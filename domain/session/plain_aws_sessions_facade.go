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
  UpdatePlainAwsSessions(plainAwsSessions []PlainAwsSession) error
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

  sessions = append(sessions, plainAwsSession)
  fac.plainAwsSessions = sessions

  err := fac.notifyObservers()
  if err != nil {
    return err
  }

  return nil
}

func(fac *plainAwsSessionsFacade) notifyObservers() error {
  for _, observer := range fac.observers {
    err := observer.UpdatePlainAwsSessions(fac.plainAwsSessions)
    if err != nil {
      return err
    }
  }

  return nil
}
