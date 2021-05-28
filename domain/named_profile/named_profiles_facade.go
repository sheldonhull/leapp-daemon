package named_profile

import (
  "fmt"
  "leapp_daemon/infrastructure/http/http_error"
  "sync"
)

var namedProfilesFacadeSingleton *namedProfilesFacade
var namedProfilesFacadeLock sync.Mutex
var namedProfilesLock sync.Mutex

type NamedProfilesObserver interface {
  UpdateNamedProfiles(namedProfiles []NamedProfile) error
}

type namedProfilesFacade struct {
  namedProfiles []NamedProfile
  observers     []NamedProfilesObserver
}

func GetNamedProfilesFacade() *namedProfilesFacade {
  namedProfilesFacadeLock.Lock()
  defer namedProfilesFacadeLock.Unlock()

  if namedProfilesFacadeSingleton == nil {
    namedProfilesFacadeSingleton = &namedProfilesFacade{
      namedProfiles: make([]NamedProfile, 0),
    }
  }

  return namedProfilesFacadeSingleton
}

func(fac *namedProfilesFacade) Subscribe(observer NamedProfilesObserver) {
  fac.observers = append(fac.observers, observer)
}

func(fac *namedProfilesFacade) GetNamedProfiles() []NamedProfile {
  return fac.namedProfiles
}

func(fac *namedProfilesFacade) SetNamedProfiles(namedProfiles []NamedProfile) {
  fac.namedProfiles = namedProfiles
}

func(fac *namedProfilesFacade) AddNamedProfile(namedProfile NamedProfile) error {
  namedProfilesLock.Lock()
  defer namedProfilesLock.Unlock()

  for _, np := range fac.namedProfiles {
    if namedProfile.Id == np.Id {
      return http_error.NewConflictError(fmt.Errorf("a NamedProfile with id " + namedProfile.Id +
        " is already present"))
    }
    if namedProfile.Name == np.Name {
      return http_error.NewConflictError(fmt.Errorf("a NamedProfile with name " + namedProfile.Name +
        " is already present"))
    }
  }

  fac.namedProfiles = append(fac.namedProfiles, namedProfile)

  err := fac.notifyObservers()
  if err != nil {
    return err
  }

  return nil
}

func(fac *namedProfilesFacade) GetNamedProfileByName(name string) *NamedProfile {
  for _, np := range fac.namedProfiles {
    if np.Name == name {
      return &np
    }
  }
  return nil
}

func(fac *namedProfilesFacade) notifyObservers() error {
  for _, observer := range fac.observers {
    err := observer.UpdateNamedProfiles(fac.namedProfiles)
    if err != nil {
      return err
    }
  }

  return nil
}
