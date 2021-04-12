package session

import (
	"github.com/pkg/errors"
	"leapp_daemon/core/uuid"
)

type NamedProfile struct {
	Id   string
	Name string
}

func CreateNamedProfile(sessionContainer Container, name string) (string, error) {

	if name == "" {
		name = "default"
	}

	namedProfiles, err := sessionContainer.GetNamedProfiles()
	if err != nil {
		return "", err
	}

	for _, namedProfile := range namedProfiles {
		if namedProfile.Name == name {
			return namedProfile.Id, nil
		}
	}

	uuidString := uuid.New()

	newNamedProfile := NamedProfile{
		Id:   uuidString,
		Name: name,
	}

	err = sessionContainer.SetNamedProfiles(append(namedProfiles, &newNamedProfile))
	if err != nil {
		return "", err
	}

	return uuidString, nil
}

func EditNamedProfile(sessionContainer Container, namedProfileId string, newName string) (string, error) {
	if newName == "" {
		newName = "default"
	}

	namedProfiles, err := sessionContainer.GetNamedProfiles()
	if err != nil {
		return "", err
	}

	for _, namedProfile := range namedProfiles {
		if namedProfile.Id == namedProfileId {
			namedProfile.Name = newName
			return namedProfileId, nil
		}
	}

	return "", errors.New("No named profile exists with Id: " + namedProfileId + ". Unable to edit profile's name")
}
