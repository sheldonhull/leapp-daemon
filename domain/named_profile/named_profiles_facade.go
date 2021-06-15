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
  UpdateNamedProfiles(oldNamedProfiles []NamedProfile, newNamedProfiles []NamedProfile) error
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

func (fac *namedProfilesFacade) Subscribe(observer NamedProfilesObserver) {
  fac.observers = append(fac.observers, observer)
}

func (fac *namedProfilesFacade) GetNamedProfiles() []NamedProfile {
  return fac.namedProfiles
}

func (fac *namedProfilesFacade) SetNamedProfiles(namedProfiles []NamedProfile) {
  fac.namedProfiles = namedProfiles
}

func (fac *namedProfilesFacade) AddNamedProfile(namedProfile NamedProfile) error {
  namedProfilesLock.Lock()
  defer namedProfilesLock.Unlock()

  namedProfiles := fac.GetNamedProfiles()

  for _, np := range namedProfiles {
    if namedProfile.Id == np.Id {
      return http_error.NewConflictError(fmt.Errorf("a NamedProfile with id " + namedProfile.Id +
        " is already present"))
    }
    if namedProfile.Name == np.Name {
      return http_error.NewConflictError(fmt.Errorf("a NamedProfile with name " + namedProfile.Name +
        " is already present"))
    }
  }

  namedProfiles = append(namedProfiles, namedProfile)

  err := fac.updateState(namedProfiles)
  if err != nil {
    return err
  }

  return nil
}

func (fac *namedProfilesFacade) GetNamedProfileByName(name string) (NamedProfile, error) {
  for _, namedProfile := range fac.namedProfiles {
    if namedProfile.Name == name {
      return namedProfile, nil
    }
  }
  return NamedProfile{}, http_error.NewNotFoundError(fmt.Errorf("named profile with name %v not found", name))
}

func (fac *namedProfilesFacade) GetNamedProfileById(id string) (NamedProfile, error) {
  for _, namedProfile := range fac.namedProfiles {
    if namedProfile.Id == id {
      return namedProfile, nil
    }
  }
  return NamedProfile{}, http_error.NewNotFoundError(fmt.Errorf("named profile with id %v not found", id))
}

func (fac *namedProfilesFacade) updateState(newState []NamedProfile) error {
  oldNamedProfiles := fac.GetNamedProfiles()
  fac.SetNamedProfiles(newState)

  for _, observer := range fac.observers {
    err := observer.UpdateNamedProfiles(oldNamedProfiles, newState)
    if err != nil {
      return err
    }
  }

  return nil
}
