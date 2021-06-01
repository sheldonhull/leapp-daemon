package session

import (
  "fmt"
  "leapp_daemon/infrastructure/http/http_error"
  "sync"
  "time"
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

  oldPlainAwsSessions := fac.GetPlainAwsSessions()
  newPlainAwsSessions := make([]PlainAwsSession, 0)

  for i := range oldPlainAwsSessions {
    newPlainAwsSession := oldPlainAwsSessions[i]
    newPlainAwsSessionAccount := *oldPlainAwsSessions[i].Account
    newPlainAwsSession.Account = &newPlainAwsSessionAccount
    newPlainAwsSessions = append(newPlainAwsSessions, newPlainAwsSession)
  }

  for _, sess := range newPlainAwsSessions {
    if plainAwsSession.Id == sess.Id {
      return http_error.NewConflictError(fmt.Errorf("a PlainAwsSession with id " + plainAwsSession.Id +
        " is already present"))
    }

    if plainAwsSession.Alias == sess.Alias {
      return http_error.NewUnprocessableEntityError(fmt.Errorf("a session with the same alias " +
        "is already present"))
    }
  }

  newPlainAwsSessions = append(newPlainAwsSessions, plainAwsSession)

  err := fac.updateState(newPlainAwsSessions)
  if err != nil {
    return err
  }

  return nil
}

func(fac *plainAwsSessionsFacade) RemovePlainAwsSession(id string) error {
  plainAwsSessionsLock.Lock()
  defer plainAwsSessionsLock.Unlock()

  oldPlainAwsSessions := fac.GetPlainAwsSessions()
  newPlainAwsSessions := make([]PlainAwsSession, 0)

  for i := range oldPlainAwsSessions {
    newPlainAwsSession := oldPlainAwsSessions[i]
    newPlainAwsSessionAccount := *oldPlainAwsSessions[i].Account
    newPlainAwsSession.Account = &newPlainAwsSessionAccount
    newPlainAwsSessions = append(newPlainAwsSessions, newPlainAwsSession)
  }

  for i, sess := range newPlainAwsSessions {
    if sess.Id == id {
      newPlainAwsSessions = append(newPlainAwsSessions[:i], newPlainAwsSessions[i+1:]...)
      break
    }
  }

  if len(fac.GetPlainAwsSessions()) == len(newPlainAwsSessions) {
    return http_error.NewNotFoundError(fmt.Errorf("plain aws session with id %s not found", id))
  }

  err := fac.updateState(newPlainAwsSessions)
  if err != nil {
    return err
  }

  return nil
}

func(fac *plainAwsSessionsFacade) GetPlainAwsSessionById(id string) (*PlainAwsSession, error) {
  for _, plainAwsSession := range fac.GetPlainAwsSessions() {
    if plainAwsSession.Id == id {
      return &plainAwsSession, nil
    }
  }
  return nil, http_error.NewNotFoundError(fmt.Errorf("plain aws session with id %s not found", id))
}

func(fac *plainAwsSessionsFacade) SetPlainAwsSessionStatusToPending(id string) error {
  plainAwsSessionsLock.Lock()
  defer plainAwsSessionsLock.Unlock()

  plainAwsSession, err := fac.GetPlainAwsSessionById(id)
  if err != nil {
    return err
  }

  if !(plainAwsSession.Status == NotActive) {
    return http_error.NewUnprocessableEntityError(fmt.Errorf("plain aws session with id " + id + "cannot be started because it's in pending or active state"))
  }

  oldPlainAwsSessions := fac.GetPlainAwsSessions()
  newPlainAwsSessions := make([]PlainAwsSession, 0)

  for i := range oldPlainAwsSessions {
    newPlainAwsSession := oldPlainAwsSessions[i]
    newPlainAwsSessionAccount := *oldPlainAwsSessions[i].Account
    newPlainAwsSession.Account = &newPlainAwsSessionAccount
    newPlainAwsSessions = append(newPlainAwsSessions, newPlainAwsSession)
  }

  for i, session := range newPlainAwsSessions {
    if session.Id == id {
      newPlainAwsSessions[i].Status = Pending
    }
  }

  err = fac.updateState(newPlainAwsSessions)
  if err != nil {
    return err
  }

  return nil
}

func(fac *plainAwsSessionsFacade) SetPlainAwsSessionStatusToActive(id string) error {
  plainAwsSessionsLock.Lock()
  defer plainAwsSessionsLock.Unlock()

  plainAwsSession, err := fac.GetPlainAwsSessionById(id)
  if err != nil {
    return err
  }

  if !(plainAwsSession.Status == Pending) {
    return http_error.NewUnprocessableEntityError(fmt.Errorf("plain aws session with id " + id + "cannot be started because it's not in pending state"))
  }

  oldPlainAwsSessions := fac.GetPlainAwsSessions()
  newPlainAwsSessions := make([]PlainAwsSession, 0)

  for i := range oldPlainAwsSessions {
    newPlainAwsSession := oldPlainAwsSessions[i]
    newPlainAwsSessionAccount := *oldPlainAwsSessions[i].Account
    newPlainAwsSession.Account = &newPlainAwsSessionAccount
    newPlainAwsSessions = append(newPlainAwsSessions, newPlainAwsSession)
  }

  for i, session := range newPlainAwsSessions {
    if session.Id == id {
      newPlainAwsSessions[i].Status = Active
      newPlainAwsSessions[i].StartTime = time.Now().Format(time.RFC3339)
    }
  }

  err = fac.updateState(newPlainAwsSessions)
  if err != nil {
    return err
  }

  return nil
}

func(fac *plainAwsSessionsFacade) SetPlainAwsSessionSessionTokenExpiration(id string, sessionTokenExpiration time.Time) error {
  plainAwsSessionsLock.Lock()
  defer plainAwsSessionsLock.Unlock()

  oldPlainAwsSessions := fac.GetPlainAwsSessions()
  newPlainAwsSessions := make([]PlainAwsSession, 0)

  for i := range oldPlainAwsSessions {
    newPlainAwsSession := oldPlainAwsSessions[i]
    newPlainAwsSessionAccount := *oldPlainAwsSessions[i].Account
    newPlainAwsSession.Account = &newPlainAwsSessionAccount
    newPlainAwsSessions = append(newPlainAwsSessions, newPlainAwsSession)
  }

  for i, session := range newPlainAwsSessions {
    if session.Id == id {
      newPlainAwsSessions[i].Account.SessionTokenExpiration = sessionTokenExpiration.Format(time.RFC3339)

      err := fac.updateState(newPlainAwsSessions)
      if err != nil {
        return err
      }

      return nil
    }
  }

  return http_error.NewNotFoundError(fmt.Errorf("plain aws session with id %s not found", id))
}

func(fac *plainAwsSessionsFacade) updateState(newState []PlainAwsSession) error {
  oldPlainAwsSessions := fac.GetPlainAwsSessions()
  fac.SetPlainAwsSessions(newState)

  for _, observer := range fac.observers {
    err := observer.UpdatePlainAwsSessions(oldPlainAwsSessions, newState)
    if err != nil {
      return err
    }
  }

  return nil
}
