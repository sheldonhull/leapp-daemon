package use_case

import (
  "leapp_daemon/domain/named_profile"
  "leapp_daemon/domain/session"
)

func ListAllSessions() (*map[string]interface{}, error) {

  sessions := session.GetGcpPlainSessionFacade().GetSessions()

	return &map[string]interface{} {
		"GcpSessions": sessions,
	}, nil
}

func ListAllNamedProfiles() ([]*named_profile.NamedProfile, error) {
	var namedProfiles []*named_profile.NamedProfile

	/*
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return namedProfiles, err
	}

	namedProfiles = config.NamedProfiles
	 */
	return namedProfiles, nil
}
