package use_case

import (
  "leapp_daemon/domain/configuration"
  "leapp_daemon/domain/session"
)

type SessionsWriter struct {
  ConfigurationRepository configuration.Repository
}

func(sessionWriter *SessionsWriter) UpdatePlainAwsSessions(plainAwsSessions []session.PlainAwsSession) error {
  config, err := sessionWriter.ConfigurationRepository.GetConfiguration()
  if err != nil {
    return err
  }

  config.PlainAwsSessions = plainAwsSessions
  err = sessionWriter.ConfigurationRepository.UpdateConfiguration(config)

  return err
}
