package session

import (
	"leapp_daemon/core/model"
	"leapp_daemon/shared/const"
)

func ListSessions(query string, listType string) (*map[string]interface{}, error) {
	plainList := make([]model.PlainAwsSession, 0)
	federatedList := make([]model.FederatedAwsSession, 0)
	var err2 error = nil

	// Check and retrieve all sessions filtered by type or by query
	if listType == "" || listType == _const.SessionTypePlain {
		plainList, err2 = ListPlainAwsSession(query)
		if err2 != nil { return nil, err2 }
	}
	if listType == "" || listType == _const.SessionTypeFederated {
		federatedList, err2 = ListFederatedAwsSession(query)
		if err2 != nil { return nil, err2 }
	}

	return &map[string]interface{} {
		"PlainSessions": plainList,
		"FederatedSessions": federatedList,
	}, nil
}
