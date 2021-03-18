package service

import (
	"leapp_daemon/core/configuration"
	"leapp_daemon/core/constant"
)

func ListAllSessions(query string, listType string) (*map[string]interface{}, error) {
	plainList := make([]*configuration.PlainAwsSession, 0)
	federatedList := make([]*configuration.FederatedAwsSession, 0)
	var err error = nil

	// Check and retrieve all sessions filtered by type or by query
	if listType == "" || listType == constant.SessionTypePlain {
		plainList, err = configuration.ListPlainAwsSession(query)
		if err != nil { return nil, err
		}
	}
	if listType == "" || listType == constant.SessionTypeFederated {
		federatedList, err = configuration.ListFederatedAwsSession(query)
		if err != nil { return nil, err
		}
	}

	return &map[string]interface{} {
		"PlainSessions": plainList,
		"FederatedSessions": federatedList,
	}, nil
}

func IsMfaRequiredForPlainAwsSession(id string) (bool, error) {
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return false, err
	}

	sess, err := configuration.GetById(config, id)
	if err != nil {
		return false, err
	}

	return sess.IsMfaRequired()
}

func StartPlainAwsSession(id string, mfaToken *string) error {
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return err
	}

	sess, err := configuration.GetById(config, id)
	if err != nil {
		return err
	}

	err = sess.Rotate(config, mfaToken)
	if err != nil {
		return err
	}

	return nil
}

func GetPlainAwsSession(id string) (*configuration.PlainAwsSession, error) {
	config, err := configuration.ReadConfiguration()
	if err != nil {
		return nil, err
	}

	sess, err := configuration.GetById(config, id)
	if err != nil {
		return nil, err
	}

	return sess, nil
}
