package session

import (
	"leapp_daemon/core/model"
	"leapp_daemon/shared/constant"
)

func ListSessions(query string, listType string) (*map[string]interface{}, error) {
	plainList := make([]model.PlainAwsSession, 0)
	federatedList := make([]model.FederatedAwsSession, 0)
	var err2 error = nil

	// Check and retrieve all sessions filtered by type or by query
	if listType == "" || listType == constant.SessionTypePlain {
		plainList, err2 = ListPlainAwsSession(query)
		if err2 != nil { return nil, err2 }
	}
	if listType == "" || listType == constant.SessionTypeFederated {
		federatedList, err2 = ListFederatedAwsSession(query)
		if err2 != nil { return nil, err2 }
	}

	return &map[string]interface{} {
		"PlainSessions": plainList,
		"FederatedSessions": federatedList,
	}, nil
}
