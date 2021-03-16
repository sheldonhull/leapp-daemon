package service

import (
	"leapp_daemon/core/model"
	"leapp_daemon/shared/constant"
	"time"
)

func RotateAllSessionsCredentials() error {
	configuration, err := ReadConfiguration()
	if err != nil { return err }

	for i := range configuration.PlainAwsSessions {
		sess := configuration.PlainAwsSessions[i]

		if sess.Active {
			isRotationIntervalExpired, err := isPlainAwsSessionRotationIntervalExpired(sess)
			if err != nil {
				return err
			}

			if isRotationIntervalExpired {
				isMfaTokenRequired, err := IsMfaRequiredForPlainAwsSession(sess.Id)
				if err != nil { return nil }

				if isMfaTokenRequired {
					// TODO: need to implement a way to ask for token.
					//  We will probably need WebSocket for this. We will exploit the Hub to wait for client response.
				} else {
					println("Rotating session with id", sess.Id)
					err = RotatePlainAwsSessionCredentials(sess, configuration, nil)
					if err != nil { return nil }
				}
			}
		}
	}

	/*for i := range configuration.FederatedAwsSessions {
		sess := configuration.FederatedAwsSessions[i]
		if sess.Active {
		}
	}*/

	return nil
}

func isPlainAwsSessionRotationIntervalExpired(session *model.PlainAwsSession) (bool, error) {
	startTime, _ := time.Parse(time.RFC3339, session.StartTime)
	secondsPassedFromStart := time.Now().Sub(startTime).Seconds()
	return int64(secondsPassedFromStart) > constant.RotationIntervalInSeconds, nil
}
