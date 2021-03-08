package session

import (
	"leapp_daemon/service/domain"
)

func ListSessions(query string, listType string) (*map[string]interface{}, error) {
	plainList := make([]domain.PlainAwsAccountSession, 0)
	federatedList := make([]domain.FederatedAwsAccountSession, 0)
	var err2 error = nil

	// Check and retrieve all sessions filtered by type or by query
	if listType == "" || listType == domain.SessionTypePlain {
		plainList, err2 = ListPlainAwsSession(query)
		if err2 != nil { return nil, err2 }
	}
	if listType == "" || listType == domain.SessionTypeFederated {
		federatedList, err2 = ListFederatedAwsSession(query)
		if err2 != nil { return nil, err2 }
	}

	return &map[string]interface{} {
		"PlainSessions": plainList,
		"FederatedSessions": federatedList,
	}, nil
}
