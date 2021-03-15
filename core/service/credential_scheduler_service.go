package service

import (
	"leapp_daemon/core/model"
	"leapp_daemon/shared/constant"
	"time"
)

func CheckAllSessions() error {
	configuration, err := ReadConfiguration()
	if err != nil { return err }

	for i := range configuration.PlainAwsSessions {
		sess := configuration.PlainAwsSessions[i]
		if sess.Active {
			err = checkTimingPlainAwsSession(&sess)
			if err != nil { return err }
		}
		configuration.PlainAwsSessions[i] = sess
	}

	for i := range configuration.FederatedAwsSessions {
		sess := configuration.FederatedAwsSessions[i]
		if sess.Active {
			err = checkTimingForFederatedAwsSession(&sess)
			if err != nil { return err }
		}
	}

	err = UpdateConfiguration(configuration, false)
	if err != nil { return err }

	return nil
}

func checkTimingPlainAwsSession(session *model.PlainAwsSession) error {
	startTime, _ := time.Parse(time.RFC3339, session.StartTime)
	secondsPassedFromStart := time.Now().Sub(startTime).Seconds()

	println("seconds passed: ", int(secondsPassedFromStart))
	println("Duration in seconds", constant.SessionTokenDurationInSeconds)

	needRefresh := int64(secondsPassedFromStart) > constant.SessionTokenDurationInSeconds

	if needRefresh {
		err := refreshAwsPlainCredentials(session)
		if err != nil { return nil }
	}

	return nil
}

func checkTimingForFederatedAwsSession(session *model.FederatedAwsSession) error {
	startTime, _ := time.Parse(time.RFC3339, session.StartTime)
	secondsPassedFromStart := time.Now().Sub(startTime).Seconds()
	needRefresh := secondsPassedFromStart > time.Duration(constant.SessionTokenDurationInSeconds).Seconds()
	if needRefresh {
		err := refreshAwsFederatedCredentials(session)
		if err != nil { return nil }
	}

	return nil
}

func refreshAwsPlainCredentials(sess *model.PlainAwsSession) error {
	sess.StartTime = time.Now().Format(time.RFC3339)
	// 1) Check if we need MFA
	isMfaTokenRequired, err := IsMfaRequiredForPlainAwsSession(sess.Id)
	if err != nil { return nil }

	// 2) Choose if we can refresh directly or ask for MFA token before proceeding
	if isMfaTokenRequired {
		//	TODO: need to implement a way to ask for token.
		// 	We will probably need WebSocket for this. We will exploit the Hub to wait for client response.
	} else {
		// 3) Start the correct Plain Aws Session again
		println("refreshing session with id", sess.Id)
		err = StartPlainAwsSession(sess.Id, "")
		if err != nil { return nil }
	}

	return nil
}

func refreshAwsFederatedCredentials(session *model.FederatedAwsSession) error {
	session.StartTime = time.Now().Format(time.RFC3339)
	return nil
}
