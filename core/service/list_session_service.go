package service

import (
	"leapp_daemon/core/configuration"
	"leapp_daemon/shared/constant"
)

func ListSessions(query string, listType string) (*map[string]interface{}, error) {
	plainList := make([]*configuration.PlainAwsSession, 0)
	federatedList := make([]*configuration.FederatedAwsSession, 0)
	var err error = nil

	// Check and retrieve all sessions filtered by type or by query
	if listType == "" || listType == constant.SessionTypePlain {
		plainList, err = configuration.List(query)
		if err != nil { return nil, err
		}
	}
	if listType == "" || listType == constant.SessionTypeFederated {
		federatedList, err = ListFederatedAwsSession(query)
		if err != nil { return nil, err
		}
	}

	return &map[string]interface{} {
		"PlainSessions": plainList,
		"FederatedSessions": federatedList,
	}, nil
}
