package use_case

import (
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/logging"
)

type AwsSessionsWriter struct {
	ConfigurationRepository ConfigurationRepository
}

func (sessionWriter *AwsSessionsWriter) UpdateAwsIamUserSessions(oldSessions []session.AwsIamUserSession, newSessions []session.AwsIamUserSession) {
	config, err := sessionWriter.ConfigurationRepository.GetConfiguration()
	if err != nil {
		logging.Entry().Error(err)
		return
	}

	config.AwsIamUserSessions = newSessions
	err = sessionWriter.ConfigurationRepository.UpdateConfiguration(config)
	if err != nil {
		logging.Entry().Error(err)
		return
	}
}
