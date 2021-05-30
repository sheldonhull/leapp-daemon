package use_case

import (
  "leapp_daemon/domain/configuration"
  "leapp_daemon/domain/session"
)

type SessionsWriter struct {
  ConfigurationRepository configuration.Repository
}

func(sessionWriter *SessionsWriter) UpdatePlainAwsSessions(oldPlainAwsSessions []session.PlainAwsSession, newPlainAwsSessions []session.PlainAwsSession) error {
  config, err := sessionWriter.ConfigurationRepository.GetConfiguration()
  if err != nil {
    return err
  }

  config.PlainAwsSessions = newPlainAwsSessions
  err = sessionWriter.ConfigurationRepository.UpdateConfiguration(config)

  return err
}
