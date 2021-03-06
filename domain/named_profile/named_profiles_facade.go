package named_profile

import (
	"fmt"
	"leapp_daemon/infrastructure/http/http_error"
	"sync"
)

var namedProfilesLock sync.Mutex

type NamedProfilesObserver interface {
	UpdateNamedProfiles(oldNamedProfiles []NamedProfile, newNamedProfiles []NamedProfile) error
}

type NamedProfilesFacade struct {
	namedProfiles []NamedProfile
	observers     []NamedProfilesObserver
}

func NewNamedProfilesFacade() *NamedProfilesFacade {
	return &NamedProfilesFacade{
		namedProfiles: make([]NamedProfile, 0),
	}
}

func (fac *NamedProfilesFacade) Subscribe(observer NamedProfilesObserver) {
	fac.observers = append(fac.observers, observer)
}

func (fac *NamedProfilesFacade) GetNamedProfiles() []NamedProfile {
	return fac.namedProfiles
}

func (fac *NamedProfilesFacade) SetNamedProfiles(namedProfiles []NamedProfile) {
	fac.namedProfiles = namedProfiles
}

func (fac *NamedProfilesFacade) AddNamedProfile(namedProfile NamedProfile) error {
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

func (fac *NamedProfilesFacade) GetNamedProfileByName(name string) (NamedProfile, error) {
	for _, namedProfile := range fac.namedProfiles {
		if namedProfile.Name == name {
			return namedProfile, nil
		}
	}
	return NamedProfile{}, http_error.NewNotFoundError(fmt.Errorf("named profile with name %v not found", name))
}

func (fac *NamedProfilesFacade) GetNamedProfileById(id string) (NamedProfile, error) {
	for _, namedProfile := range fac.namedProfiles {
		if namedProfile.Id == id {
			return namedProfile, nil
		}
	}
	return NamedProfile{}, http_error.NewNotFoundError(fmt.Errorf("named profile with id %v not found", id))
}

func (fac *NamedProfilesFacade) updateState(newState []NamedProfile) error {
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
