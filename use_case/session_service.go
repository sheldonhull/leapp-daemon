package use_case

import (
  "leapp_daemon/domain/named_profile"
  session2 "leapp_daemon/domain/session"
)

func ListAllSessions(query string, listType string) (*map[string]interface{}, error) {
	plainAwsList := make([]*session2.PlainAwsSession, 0)
	federatedAwsList := make([]*session2.FederatedAwsSession, 0)

	plainAlibabaList := make([]*session2.PlainAlibabaSession, 0)

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
		"PlainAwsSessions": plainAwsList,
		"FederatedAwsSessions": federatedAwsList,
		"PlainAlibabaSession": plainAlibabaList,
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
