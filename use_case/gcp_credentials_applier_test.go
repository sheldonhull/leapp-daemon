package use_case

import (
	"leapp_daemon/domain/gcp"
	"leapp_daemon/domain/gcp/gcp_iam_user_account_oauth"
	"leapp_daemon/test/mock"
	"reflect"
	"testing"
)

var (
	keychainMock              mock.KeychainMock
	gcpRepoMock               mock.GcpConfigurationRepositoryMock
	gcpCredentialsApplier     *GcpCredentialsApplier
	expectedDeactivationCalls []string
	expectedActivationCalls   []string
)

func gcpCredentialsApplierSetup() {
	expectedDeactivationCalls = []string{"RemoveDefaultCredentials()", "DeactivateConfiguration()",
		"RemoveCredentialsFromDb()", "RemoveAccessTokensFromDb()", "RemoveConfiguration()"}
	expectedActivationCalls = []string{"WriteDefaultCredentials(accountId, credentials)",
		"CreateConfiguration(accountId, projectName)", "ActivateConfiguration()", "WriteDefaultCredentials(credentials)"}

	keychainMock = mock.NewKeychainMock()
	gcpRepoMock = mock.NewGcpConfigurationRepositoryMock()
	gcpCredentialsApplier = &GcpCredentialsApplier{
		Keychain:   &keychainMock,
		Repository: &gcpRepoMock,
	}
}

func gcpCredentialsApplierVerifyExpectedCalls(t *testing.T, keychainMockCalls []string, repoMockCalls []string) {
	if !reflect.DeepEqual(keychainMock.GetCalls(), keychainMockCalls) {
		t.Fatalf("keychainMock expectation violation.\nMock calls: %v", keychainMock.GetCalls())
	}
	if !reflect.DeepEqual(gcpRepoMock.GetCalls(), repoMockCalls) {
		t.Fatalf("gcpRepoMock expectation violation.\nMock calls: %v", gcpRepoMock.GetCalls())
	}
}

func TestUpdateGcpIamUserAccountOauthSessions_OldActiveSessionAndNoNewActiveSessions(t *testing.T) {
	gcpCredentialsApplierSetup()
	oldSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Status: gcp.Active}}
	newSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{}

	gcpCredentialsApplier.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	gcpCredentialsApplierVerifyExpectedCalls(t, []string{}, expectedDeactivationCalls)
}

func TestUpdateGcpIamUserAccountOauthSessions_OldAndNewActiveSessionWithDifferentIds(t *testing.T) {
	gcpCredentialsApplierSetup()
	keychainMock.ExpGetSecret = "credentials"
	oldSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID1", Status: gcp.Active}}
	newSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID2", CredentialsLabel: "credentialsLabel", AccountId: "accountId",
		Name: "sessionName", ProjectName: "projectName", Status: gcp.Active}}

	gcpCredentialsApplier.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	expectedRepositoryCalls := append(expectedDeactivationCalls, expectedActivationCalls...)
	gcpCredentialsApplierVerifyExpectedCalls(t, []string{"GetSecret(credentialsLabel)"}, expectedRepositoryCalls)
}

func TestUpdateGcpIamUserAccountOauthSessions_OldAndNewActiveSessionAreEqual(t *testing.T) {
	gcpCredentialsApplierSetup()
	oldSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID1", Status: gcp.Active}}
	newSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID1", Status: gcp.Active}}

	gcpCredentialsApplier.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	gcpCredentialsApplierVerifyExpectedCalls(t, []string{}, []string{})
}

func TestUpdateGcpIamUserAccountOauthSessions_OldAndNewActiveSessionWithSameIdsButDifferentParams(t *testing.T) {
	gcpCredentialsApplierSetup()
	keychainMock.ExpGetSecret = "credentials"
	oldSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID1", CredentialsLabel: "credentialsLabel", AccountId: "accountId",
		Name: "sessionName", ProjectName: "oldProjectName", Status: gcp.Active}}
	newSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID1", CredentialsLabel: "credentialsLabel", AccountId: "accountId",
		Name: "sessionName", ProjectName: "projectName", Status: gcp.Active}}

	gcpCredentialsApplier.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	gcpCredentialsApplierVerifyExpectedCalls(t, []string{"GetSecret(credentialsLabel)"}, expectedActivationCalls)
}

func TestUpdateGcpIamUserAccountOauthSessions_NoOldActiveSessionButNewActiveSessionPresent(t *testing.T) {
	gcpCredentialsApplierSetup()
	keychainMock.ExpGetSecret = "credentials"
	oldSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID1", Status: gcp.NotActive}}
	newSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID1", CredentialsLabel: "credentialsLabel", AccountId: "accountId",
		Name: "sessionName", ProjectName: "projectName", Status: gcp.Active}}

	gcpCredentialsApplier.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	gcpCredentialsApplierVerifyExpectedCalls(t, []string{"GetSecret(credentialsLabel)"}, expectedActivationCalls)
}

func TestUpdateGcpIamUserAccountOauthSessions_NoActiveSessions(t *testing.T) {
	gcpCredentialsApplierSetup()
	keychainMock.ExpGetSecret = "credentials"
	oldSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID1", Status: gcp.NotActive}}
	newSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID1", Status: gcp.NotActive}}

	gcpCredentialsApplier.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	gcpCredentialsApplierVerifyExpectedCalls(t, []string{}, []string{})
}
