package use_case

import (
	"leapp_daemon/domain/session"
	"leapp_daemon/test/mock"
	"reflect"
	"testing"
)

var (
	keychainMock              mock.KeychainMock
	repoMock                  mock.GcpConfigurationRepositoryMock
	applier                   *GcpCredentialsApplier
	expectedDeactivationCalls []string
	expectedActivationCalls   []string
)

func gcpCredentialsApplierSetup() {
	expectedDeactivationCalls = []string{"RemoveDefaultCredentials()", "DeactivateConfiguration()",
		"RemoveCredentialsFromDb()", "RemoveAccessTokensFromDb()", "RemoveConfiguration()"}
	expectedActivationCalls = []string{"WriteDefaultCredentials(accountId, credentials)",
		"CreateConfiguration(sessionName, accountId, projectName)", "ActivateConfiguration()", "WriteDefaultCredentials(credentials)"}

	keychainMock = mock.NewKeychainMock()
	repoMock = mock.NewGcpConfigurationRepositoryMock()
	applier = &GcpCredentialsApplier{
		Keychain:   &keychainMock,
		Repository: &repoMock,
	}
}

func verifyExpectedCalls(t *testing.T, keychainMockCalls []string, repoMockCalls []string) {
	if !reflect.DeepEqual(keychainMock.GetCalls(), keychainMockCalls) {
		t.Fatalf("keychainMock expectation violation.\nMock calls: %v", keychainMock.GetCalls())
	}
	if !reflect.DeepEqual(repoMock.GetCalls(), repoMockCalls) {
		t.Fatalf("repoMock expectation violation.\nMock calls: %v", repoMock.GetCalls())
	}
}

func TestUpdateGcpPlainSessions_OldActiveSessionAndNoNewActiveSessions(t *testing.T) {
	gcpCredentialsApplierSetup()
	oldSessions := []session.GcpPlainSession{{Status: session.Active}}
	newSessions := []session.GcpPlainSession{}

	applier.UpdateGcpPlainSessions(oldSessions, newSessions)
	verifyExpectedCalls(t, []string{}, expectedDeactivationCalls)
}

func TestUpdateGcpPlainSessions_OldAndNewActiveSessionWithDifferentIds(t *testing.T) {
	gcpCredentialsApplierSetup()
	keychainMock.ExpGetSecret = "credentials"
	oldSessions := []session.GcpPlainSession{{Id: "ID1", Status: session.Active}}
	newSessions := []session.GcpPlainSession{{Id: "ID2", CredentialsLabel: "credentialsLabel", AccountId: "accountId",
		Name: "sessionName", ProjectName: "projectName", Status: session.Active}}

	applier.UpdateGcpPlainSessions(oldSessions, newSessions)
	expectedRepositoryCalls := append(expectedDeactivationCalls, expectedActivationCalls...)
	verifyExpectedCalls(t, []string{"GetSecret(credentialsLabel)"}, expectedRepositoryCalls)
}

func TestUpdateGcpPlainSessions_OldAndNewActiveSessionAreEqual(t *testing.T) {
	gcpCredentialsApplierSetup()
	oldSessions := []session.GcpPlainSession{{Id: "ID1", Status: session.Active}}
	newSessions := []session.GcpPlainSession{{Id: "ID1", Status: session.Active}}

	applier.UpdateGcpPlainSessions(oldSessions, newSessions)
	verifyExpectedCalls(t, []string{}, []string{})
}

func TestUpdateGcpPlainSessions_OldAndNewActiveSessionWithSameIdsButDifferentParams(t *testing.T) {
	gcpCredentialsApplierSetup()
	keychainMock.ExpGetSecret = "credentials"
	oldSessions := []session.GcpPlainSession{{Id: "ID1", CredentialsLabel: "credentialsLabel", AccountId: "accountId",
		Name: "sessionName", ProjectName: "oldProjectName", Status: session.Active}}
	newSessions := []session.GcpPlainSession{{Id: "ID1", CredentialsLabel: "credentialsLabel", AccountId: "accountId",
		Name: "sessionName", ProjectName: "projectName", Status: session.Active}}

	applier.UpdateGcpPlainSessions(oldSessions, newSessions)
	verifyExpectedCalls(t, []string{"GetSecret(credentialsLabel)"}, expectedActivationCalls)
}

func TestUpdateGcpPlainSessions_NoOldActiveSessionButNewActiveSessionPresent(t *testing.T) {
	gcpCredentialsApplierSetup()
	keychainMock.ExpGetSecret = "credentials"
	oldSessions := []session.GcpPlainSession{{Id: "ID1", Status: session.NotActive}}
	newSessions := []session.GcpPlainSession{{Id: "ID1", CredentialsLabel: "credentialsLabel", AccountId: "accountId",
		Name: "sessionName", ProjectName: "projectName", Status: session.Active}}

	applier.UpdateGcpPlainSessions(oldSessions, newSessions)
	verifyExpectedCalls(t, []string{"GetSecret(credentialsLabel)"}, expectedActivationCalls)
}

func TestUpdateGcpPlainSessions_NoActiveSessions(t *testing.T) {
	gcpCredentialsApplierSetup()
	keychainMock.ExpGetSecret = "credentials"
	oldSessions := []session.GcpPlainSession{{Id: "ID1", Status: session.NotActive}}
	newSessions := []session.GcpPlainSession{{Id: "ID1", Status: session.Pending}}

	applier.UpdateGcpPlainSessions(oldSessions, newSessions)
	verifyExpectedCalls(t, []string{}, []string{})
}

/*
credentials, err := applier.Keychain.GetSecret(session.CredentialsLabel)
	if err != nil {
		logging.Entry().Error(err)
		return
	}

	err = applier.Repository.WriteCredentialsToDb(session.AccountId, credentials)
	if err != nil {
		logging.Entry().Error(err)
		return
	}

	err = applier.Repository.CreateConfiguration(session.Name, session.AccountId, session.ProjectName)
	if err != nil {
		logging.Entry().Error(err)
		return
	}

	err = applier.Repository.ActivateConfiguration()
	if err != nil {
		logging.Entry().Error(err)
		return
	}

	err = applier.Repository.WriteDefaultCredentials(credentials)
	if err != nil {
		logging.Entry().Error(err)
		return
	}
*/
