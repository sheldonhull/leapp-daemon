package service

import (
	"fmt"
	"leapp_daemon/core/model"
	"leapp_daemon/shared/constant"
	"time"
)

func CheckAllSessions() error {
	configuration, err := ReadConfiguration()
	if err != nil {
		return err
	}

	federatedAwsSessions := configuration.FederatedAwsSessions
	plainAwsSessions := configuration.PlainAwsSessions

	for _, session := range federatedAwsSessions {
		if session.Active {
			checkTimingForFederatedAwsSession(&session)
		}
	}

	for _, session := range plainAwsSessions {
		// if session.Active {
			checkTimingPlainAwsSession(&session)
		// }
	}

	UpdateConfiguration(configuration, false)

	return nil
}

func checkTimingForFederatedAwsSession(session *model.FederatedAwsSession) {
	startTime, _ := time.Parse(time.RFC3339, session.StartTime)
	secondsPassedFromStart := time.Now().Sub(startTime).Seconds()

	fmt.Println(secondsPassedFromStart)

	needRefresh := secondsPassedFromStart > time.Duration(constant.SessionTokenDurationInSeconds).Seconds()
	if needRefresh {
		go refreshAwsFederatedCredentials(session)
	}
}

func checkTimingPlainAwsSession(session *model.PlainAwsSession) {
	fmt.Println(session.StartTime)

	startTime, _ := time.Parse(time.RFC3339, session.StartTime)

	fmt.Println(startTime)

	secondsPassedFromStart := time.Now().Sub(startTime).Seconds()

	fmt.Println(secondsPassedFromStart)

	needRefresh := secondsPassedFromStart > time.Duration(constant.SessionTokenDurationInSeconds).Seconds()
	if needRefresh {
		go refreshAwsPlainCredentials(session)
	}
}

func refreshAwsFederatedCredentials(session *model.FederatedAwsSession) {
	session.StartTime = time.Now().Format(time.RFC3339)
	println("Federated session " + session.Account.Name + " started at " + session.StartTime)

}

func refreshAwsPlainCredentials(session *model.PlainAwsSession) {
	session.StartTime = time.Now().Format(time.RFC3339)
	println("Plain session " + session.Account.Name + " started at " + session.StartTime)
}
