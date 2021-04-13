package use_case

import (
  session2 "leapp_daemon/domain/session"
)

func ListAllSessions(query string, listType string) (*map[string]interface{}, error) {
	plainList := make([]*session2.PlainAwsSession, 0)
	federatedList := make([]*session2.FederatedAwsSession, 0)

	/*
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return &map[string]interface{} {}, err
	}

	// Check and retrieve all sessions filtered by type or by query
	if listType == "" || listType == constant.SessionTypePlain {
		plainList, err = session2.ListPlainAwsSession(config, query)
		if err != nil { return nil, err
		}
	}
	if listType == "" || listType == constant.SessionTypeFederated {
		federatedList, err = session2.ListFederatedAwsSession(config, query)
		if err != nil { return nil, err
		}
	}
	 */

	return &map[string]interface{} {
		"PlainSessions": plainList,
		"FederatedSessions": federatedList,
	}, nil
}

func ListAllNamedProfiles() ([]*session2.NamedProfile, error) {
	var namedProfiles []*session2.NamedProfile

	/*
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return namedProfiles, err
	}

	namedProfiles = config.NamedProfiles
	 */
	return namedProfiles, nil
}
