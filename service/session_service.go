package service

import (
	"leapp_daemon/core/configuration"
	"leapp_daemon/core/constant"
	"leapp_daemon/core/session"
	"strings"
)

func ListAllSessions(query string, listType string) (*map[string]interface{}, error) {
	plainList := make([]*session.PlainAwsSession, 0)
	federatedList := make([]*session.FederatedAwsSession, 0)
	var err error = nil

	// Check and retrieve all sessions filtered by type or by query
	if listType == "" || listType == constant.SessionTypePlain {
		plainList, err = ListPlainAwsSession(query)
		if err != nil { return nil, err
		}
	}
	if listType == "" || listType == constant.SessionTypeFederated {
		federatedList, err = session.ListFederatedAwsSession(query)
		if err != nil { return nil, err
		}
	}

	return &map[string]interface{} {
		"PlainSessions": plainList,
		"FederatedSessions": federatedList,
	}, nil
}

func ListPlainAwsSession(query string) ([]*session.PlainAwsSession, error) {
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return nil, err
	}

	filteredList := make([]*session.PlainAwsSession, 0)

	if query == "" {
		return append(filteredList, config.PlainAwsSessions...), nil
	} else {
		for _, sess := range config.PlainAwsSessions {
			if strings.Contains(sess.Id, query) ||
				strings.Contains(sess.Account.Name, query) ||
				strings.Contains(sess.Account.MfaDevice, query) ||
				strings.Contains(sess.Account.User, query) ||
				strings.Contains(sess.Account.Region, query) ||
				strings.Contains(sess.Account.AccountNumber, query) {

				filteredList = append(filteredList, sess)
			}
		}

		return filteredList, nil
	}
}
