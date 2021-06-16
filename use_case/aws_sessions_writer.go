package use_case

import (
	"leapp_daemon/domain/session"
)

type AwsSessionsWriter struct {
	ConfigurationRepository ConfigurationRepository
}

func (sessionWriter *AwsSessionsWriter) UpdateAwsPlainSessions(oldSessions []session.AwsPlainSession, newSessions []session.AwsPlainSession) error {
	config, err := sessionWriter.ConfigurationRepository.GetConfiguration()
	if err != nil {
		return err
	}

	config.PlainAwsSessions = newSessions
	err = sessionWriter.ConfigurationRepository.UpdateConfiguration(config)

	return err
}
