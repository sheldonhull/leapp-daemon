package use_case

import (
	"leapp_daemon/domain/named_profile"
)

type NamedProfilesActions struct {
	Environment Environment
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
		namedProfile = named_profile.NamedProfile{
			Id:   actions.Environment.GenerateUuid(),
			Name: profileName,
		}
		err = facade.AddNamedProfile(namedProfile)
		if err != nil {
			return named_profile.NamedProfile{}, err
		}
	}

	return namedProfile, nil
}
