package use_case

import (
	"leapp_daemon/domain/session"
	"leapp_daemon/interface/repository"
)

type GcpCredentialsApplier struct {
	Repository repository.GcloudConfigurationRepository
	Keychain   Keychain
}

func (awsCredentialsApplier *AwsCredentialsApplier) UpdateGcpPlainSessions(oldSessions []session.GcpPlainSession, newSessions []session.GcpPlainSession) error {
	return nil
}
