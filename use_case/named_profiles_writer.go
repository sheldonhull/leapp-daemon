package use_case

import (
  "leapp_daemon/domain/configuration"
  "leapp_daemon/domain/named_profile"
)

type NamedProfilesWriter struct {
  ConfigurationRepository configuration.Repository
}

func(namedProfilesWriter *NamedProfilesWriter) UpdateNamedProfiles(oldNamedProfiles []named_profile.NamedProfile, newNamedProfiles []named_profile.NamedProfile) error {
  config, err := namedProfilesWriter.ConfigurationRepository.GetConfiguration()
  if err != nil {
    return err
  }

  config.NamedProfiles = newNamedProfiles
  err = namedProfilesWriter.ConfigurationRepository.UpdateConfiguration(config)

  return err
}
