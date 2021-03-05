package sessions

import (
	"leapp_daemon/services/domain"
)

func ListSessions(query string, listType string) (*map[string]interface{}, error) {
	plainList := make([]domain.PlainAwsAccountSession, 0)
	federatedList := make([]domain.FederatedAwsAccountSession, 0)
	var err2 error = nil

	// Check and retrieve all sessions filtered by type or by query
	if listType == "" || listType == domain.SESSION_TYPE_PLAIN {
		plainList, err2 = ListPlainAwsSession(query)
		if err2 != nil { return nil, err2 }
	}
	if listType == "" || listType == domain.SESSION_TYPE_FEDERATED {
		federatedList, err2 = ListFederatedAwsSession(query)
		if err2 != nil { return nil, err2 }
	}

	return &map[string]interface{} {
		"PlainSessions": plainList,
		"FederatedSessions": federatedList,
	}, nil
}
