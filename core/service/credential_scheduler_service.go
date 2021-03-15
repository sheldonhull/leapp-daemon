package service

import (
	"leapp_daemon/core/model"
	"leapp_daemon/shared/constant"
	"time"
)

func CheckAllSessions() error {
	configuration, err := ReadConfiguration()
	if err != nil {
		return err
	}

	for i, _ := range configuration.PlainAwsSessions {
		session := configuration.PlainAwsSessions[i]
		if session.Active {
			checkTimingPlainAwsSession(&session)
		}
		configuration.PlainAwsSessions[i] = session
	}

	for i, _ := range configuration.FederatedAwsSessions {
		session := configuration.FederatedAwsSessions[i]
		if session.Active {
			checkTimingForFederatedAwsSession(&session)
		}
	}

	UpdateConfiguration(configuration, false)

	return nil
}

func checkTimingPlainAwsSession(session *model.PlainAwsSession) {
	startTime, _ := time.Parse(time.RFC3339, session.StartTime)
	secondsPassedFromStart := time.Now().Sub(startTime).Seconds()
	needRefresh := secondsPassedFromStart > time.Duration(constant.SessionTokenDurationInSeconds).Seconds()

	if needRefresh {
		refreshAwsPlainCredentials(session)
	}
}

func checkTimingForFederatedAwsSession(session *model.FederatedAwsSession) {
	startTime, _ := time.Parse(time.RFC3339, session.StartTime)
	secondsPassedFromStart := time.Now().Sub(startTime).Seconds()
	needRefresh := secondsPassedFromStart > time.Duration(constant.SessionTokenDurationInSeconds).Seconds()
	if needRefresh {
		refreshAwsFederatedCredentials(session)
	}
}

func refreshAwsPlainCredentials(session *model.PlainAwsSession) {
	session.StartTime = time.Now().Format(time.RFC3339)
}

func refreshAwsFederatedCredentials(session *model.FederatedAwsSession) {
	session.StartTime = time.Now().Format(time.RFC3339)
}
