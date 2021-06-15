package use_case

import (
  "github.com/google/uuid"
  "leapp_daemon/domain/named_profile"
  "strings"
)

type NamedProfilesActions struct {
}

func (actions *NamedProfilesActions) GetNamedProfileById(profileId string) (named_profile.NamedProfile, error) {
  return named_profile.GetNamedProfilesFacade().GetNamedProfileById(profileId)
}

func (actions *NamedProfilesActions) GetOrCreateNamedProfile(profileName string) (named_profile.NamedProfile, error) {
  if profileName == "" {
    profileName = "default"
  }

  facade := named_profile.GetNamedProfilesFacade()
  namedProfile, err := facade.GetNamedProfileByName(profileName)
  if err != nil {
    uuidString := uuid.New().String()
    uuidString = strings.Replace(uuidString, "-", "", -1)

    namedProfile = named_profile.NamedProfile{
      Id:   uuidString,
      Name: profileName,
    }
    err = facade.AddNamedProfile(namedProfile)
    if err != nil {
      return named_profile.NamedProfile{}, err
    }
  }

  return namedProfile, nil
}
