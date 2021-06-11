package use_case

import (
  "leapp_daemon/domain/configuration"
  "leapp_daemon/domain/session"
)

type AwsSessionsWriter struct {
  ConfigurationRepository configuration.Repository
}

func(sessionWriter *AwsSessionsWriter) UpdatePlainAwsSessions(oldSessions []session.PlainAwsSession, newSessions []session.PlainAwsSession) error {
  config, err := sessionWriter.ConfigurationRepository.GetConfiguration()
  if err != nil {
    return err
  }

  config.PlainAwsSessions = newSessions
  err = sessionWriter.ConfigurationRepository.UpdateConfiguration(config)

  return err
}
