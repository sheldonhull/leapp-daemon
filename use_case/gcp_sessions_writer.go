package use_case

import (
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/logging"
)

type GcpSessionsWriter struct {
	ConfigurationRepository ConfigurationRepository
}

func (sessionWriter *GcpSessionsWriter) UpdateGcpPlainSessions(oldSessions []session.GcpPlainSession, newSessions []session.GcpPlainSession) {
	config, err := sessionWriter.ConfigurationRepository.GetConfiguration()
	if err != nil {
		logging.Entry().Error(err)
		return
	}

	config.GcpPlainSessions = newSessions
	err = sessionWriter.ConfigurationRepository.UpdateConfiguration(config)
	if err != nil {
		logging.Entry().Error(err)
		return
	}
}
