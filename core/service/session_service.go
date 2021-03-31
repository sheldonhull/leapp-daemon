package service

import (
	"leapp_daemon/core/configuration"
	"leapp_daemon/core/constant"
	"leapp_daemon/core/session"
)

func ListAllSessions(query string, listType string) (*map[string]interface{}, error) {
	plainList := make([]*session.PlainAwsSession, 0)
	federatedList := make([]*session.FederatedAwsSession, 0)

	config, err := configuration.ReadConfiguration()
	if err != nil {
		return &map[string]interface{} {}, err
	}

	// Check and retrieve all sessions filtered by type or by query
	if listType == "" || listType == constant.SessionTypePlain {
		plainList, err = session.ListPlainAwsSession(config, query)
		if err != nil { return nil, err
		}
	}
	if listType == "" || listType == constant.SessionTypeFederated {
		federatedList, err = session.ListFederatedAwsSession(config, query)
		if err != nil { return nil, err
		}
	}

	return &map[string]interface{} {
		"PlainSessions": plainList,
		"FederatedSessions": federatedList,
	}, nil
}

func ListAllSNamedProfiles() ([]*session.NamedProfile, error) {
	var namedProfiles []*session.NamedProfile

	config, err := configuration.ReadConfiguration()
	if err != nil {
		return namedProfiles, err
	}

	namedProfiles = config.NamedProfiles
	return namedProfiles, nil
}
