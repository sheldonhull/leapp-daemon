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
  UpdateGcpPlainSessions(oldSessions []GcpPlainSession, newSessions []GcpPlainSession)
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

  facade.updateState(newSessions)
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

  facade.updateState(newSessions)
  return nil
}

func (facade *gcpPlainSessionsFacade) SetSessionStatus(sessionId string, status Status) error {
  sessionsLock.Lock()
  defer sessionsLock.Unlock()

  sessionToUpdate, err := facade.GetSessionById(sessionId)
  if err != nil {
    return err
  }

  sessionToUpdate.Status = status
  return facade.replaceSession(sessionId, sessionToUpdate)
}

func (facade *gcpPlainSessionsFacade) EditSession(sessionId string, name string, projectName string, profileId string) error {
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

func (facade *gcpPlainSessionsFacade) replaceSession(sessionId string, newSession GcpPlainSession) error {
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

func (facade *gcpPlainSessionsFacade) updateState(newSessions []GcpPlainSession) {
  oldSessions := facade.GetSessions()
  facade.gcpPlainSessions = newSessions

  for _, observer := range facade.observers {
    observer.UpdateGcpPlainSessions(oldSessions, newSessions)
  }
}
