package use_case

import (
	"leapp_daemon/domain/session"
	"leapp_daemon/test/mock"
	"reflect"
	"testing"
)

var (
	oldSessions       []session.GcpPlainSession
	newSessions       []session.GcpPlainSession
	fileRepoMock      mock.FileConfigurationRepositoryMock
	gcpSessionsWriter *GcpSessionsWriter
)

func gcpSessionsWriterSetup() {
	oldSessions = []session.GcpPlainSession{}
	newSessions = []session.GcpPlainSession{{Id: "ID"}}

	fileRepoMock = mock.NewFileConfigurationRepositoryMock()
	gcpSessionsWriter = &GcpSessionsWriter{
		ConfigurationRepository: &fileRepoMock,
	}
}

func gcpSessionsWriterVerifyExpectedCalls(t *testing.T, fileRepoMockCalls []string) {
	if !reflect.DeepEqual(fileRepoMock.GetCalls(), fileRepoMockCalls) {
		t.Fatalf("fileRepoMock expectation violation.\nMock calls: %v", fileRepoMock.GetCalls())
	}
}

func TestUpdateGcpPlainSessions(t *testing.T) {
	gcpSessionsWriterSetup()

	gcpSessionsWriter.UpdateGcpPlainSessions(oldSessions, newSessions)
	gcpSessionsWriterVerifyExpectedCalls(t, []string{"GetConfiguration()", "UpdateConfiguration()"})
}

func TestUpdateGcpPlainSessions_ErrorGettingConfiguration(t *testing.T) {
	gcpSessionsWriterSetup()
	fileRepoMock.ExpErrorOnGetConfiguration = true

	gcpSessionsWriter.UpdateGcpPlainSessions(oldSessions, newSessions)
	gcpSessionsWriterVerifyExpectedCalls(t, []string{"GetConfiguration()"})
}

func TestUpdateGcpPlainSessions_ErrorUpdatingConfiguration(t *testing.T) {
	gcpSessionsWriterSetup()
	fileRepoMock.ExpErrorOnUpdateConfiguration = true

	gcpSessionsWriter.UpdateGcpPlainSessions(oldSessions, newSessions)
	gcpSessionsWriterVerifyExpectedCalls(t, []string{"GetConfiguration()", "UpdateConfiguration()"})
}
