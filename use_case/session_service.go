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

func ListAllNamedProfiles() ([]named_profile.NamedProfile, error) {
  return named_profile.GetNamedProfilesFacade().GetNamedProfiles(), nil
}
