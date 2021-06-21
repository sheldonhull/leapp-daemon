package use_case

import (
	"golang.org/x/oauth2"
	"leapp_daemon/domain/session"
	"leapp_daemon/test"
	"leapp_daemon/test/mock"
	"net/http"
	"reflect"
	"testing"
)

var (
	gcpApiMock                         mock.GcpApiMock
	envMock                            mock.EnvironmentMock
	gcpPlainSessionActionsKeychainMock mock.KeychainMock
	gcpPlainSessionFacadeMock          mock.GcpPlainSessionsFacadeMock
	gcpPlainSessionActions             *GcpPlainSessionActions
)

func gcpPlainSessionActionsSetup() {
	gcpApiMock = mock.NewGcpApiMock()
	envMock = mock.NewEnvironmentMock()
	gcpPlainSessionActionsKeychainMock = mock.NewKeychainMock()
	gcpPlainSessionFacadeMock = mock.NewGcpPlainSessionsFacadeMock()
	gcpPlainSessionActions = &GcpPlainSessionActions{
		GcpApi:                &gcpApiMock,
		Environment:           &envMock,
		Keychain:              &gcpPlainSessionActionsKeychainMock,
		GcpPlainSessionFacade: &gcpPlainSessionFacadeMock,
	}
}

func gcpPlainSessionActionsVerifyExpectedCalls(t *testing.T, gcpApiMockCalls, envMockCalls, keychainMockCalls, facadeMockCalls []string) {
	if !reflect.DeepEqual(gcpApiMock.GetCalls(), gcpApiMockCalls) {
		t.Fatalf("gcpApiMock expectation violation.\nMock calls: %v", gcpApiMock.GetCalls())
	}
	if !reflect.DeepEqual(envMock.GetCalls(), envMockCalls) {
		t.Fatalf("envMock expectation violation.\nMock calls: %v", envMock.GetCalls())
	}
	if !reflect.DeepEqual(gcpPlainSessionActionsKeychainMock.GetCalls(), keychainMockCalls) {
		t.Fatalf("keychainMock expectation violation.\nMock calls: %v", gcpPlainSessionActionsKeychainMock.GetCalls())
	}
	if !reflect.DeepEqual(gcpPlainSessionFacadeMock.GetCalls(), facadeMockCalls) {
		t.Fatalf("facadeMock expectation violation.\nMock calls: %v", gcpPlainSessionFacadeMock.GetCalls())
	}
}

func TestGetSession(t *testing.T) {
	gcpPlainSessionActionsSetup()

	session := session.GcpPlainSession{Name: "test_session"}
	gcpPlainSessionFacadeMock.ExpGetSessionById = session

	actualSession, err := gcpPlainSessionActions.GetSession("ID")
	if err != nil && !reflect.DeepEqual(session, actualSession) {
		t.Fatalf("Returned unexpected session")
	}
	gcpPlainSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{}, []string{"GetSessionById(ID)"})
}

func TestGetSession_SessionFacadeReturnsError(t *testing.T) {
	gcpPlainSessionActionsSetup()
	gcpPlainSessionFacadeMock.ExpErrorOnGetSessionById = true

	_, err := gcpPlainSessionActions.GetSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session not found")
	gcpPlainSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{}, []string{"GetSessionById(ID)"})
}

func TestGetOAuthUrl(t *testing.T) {
	gcpPlainSessionActionsSetup()
	gcpApiMock.ExpOauthUrl = "url"

	actualOauthUrl, err := gcpPlainSessionActions.GetOAuthUrl()
	if err != nil && !reflect.DeepEqual("url", actualOauthUrl) {
		t.Fatalf("Returned unexpected oauth url")
	}
	gcpPlainSessionActionsVerifyExpectedCalls(t, []string{"GetOauthUrl()"}, []string{}, []string{}, []string{})
}

func TestGetOAuthUrl_GcpApiReturnsError(t *testing.T) {
	gcpPlainSessionActionsSetup()
	gcpApiMock.ExpErrorOnGetOauthUrl = true

	_, err := gcpPlainSessionActions.GetOAuthUrl()
	test.ExpectHttpError(t, err, http.StatusNotFound, "error getting oauth url")
	gcpPlainSessionActionsVerifyExpectedCalls(t, []string{"GetOauthUrl()"}, []string{}, []string{}, []string{})
}

func TestCreateSession(t *testing.T) {
	gcpPlainSessionActionsSetup()

	sessionName := "sessionName"
	accountId := "accountId"
	projectName := "projectName"
	oauthCode := "oauthCode"
	uuid := "uuid"
	credentials := "credentialsJson"
	envMock.ExpUuid = uuid
	gcpApiMock.ExpOauthToken = oauth2.Token{}
	gcpApiMock.ExpCredentials = credentials

	err := gcpPlainSessionActions.CreateSession(sessionName, accountId, projectName, oauthCode)
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpPlainSessionActionsVerifyExpectedCalls(t, []string{"GetOauthUrl(oauthCode)", "GetCredentials()"}, []string{"GenerateUuid()"}, []string{"SetSecret(credentialsJson, uuid-gcp-plain-session-credentials)"}, []string{"AddSession(sessionName)"})
}

func TestCreateSession_GcpApiGetOauthTokenReturnsError(t *testing.T) {
	gcpPlainSessionActionsSetup()

	sessionName := "sessionName"
	accountId := "accountId"
	projectName := "projectName"
	oauthCode := "oauthCode"
	uuid := "uuid"
	envMock.ExpUuid = uuid
	gcpApiMock.ExpErrorOnGetOauth = true

	err := gcpPlainSessionActions.CreateSession(sessionName, accountId, projectName, oauthCode)
	test.ExpectHttpError(t, err, http.StatusBadRequest, "error getting oauth token")
	gcpPlainSessionActionsVerifyExpectedCalls(t, []string{"GetOauthUrl(oauthCode)"}, []string{"GenerateUuid()"}, []string{}, []string{})
}

func TestCreateSession_KeychainSetSecretReturnsError(t *testing.T) {
	gcpPlainSessionActionsSetup()

	sessionName := "sessionName"
	accountId := "accountId"
	projectName := "projectName"
	oauthCode := "oauthCode"
	uuid := "uuid"
	credentials := "credentialsJson"
	envMock.ExpUuid = uuid
	gcpApiMock.ExpOauthToken = oauth2.Token{}
	gcpApiMock.ExpCredentials = credentials
	gcpPlainSessionActionsKeychainMock.ExpErrorOnSetSecret = true

	err := gcpPlainSessionActions.CreateSession(sessionName, accountId, projectName, oauthCode)
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "unable to set secret")
	gcpPlainSessionActionsVerifyExpectedCalls(t, []string{"GetOauthUrl(oauthCode)", "GetCredentials()"}, []string{"GenerateUuid()"}, []string{"SetSecret(credentialsJson, uuid-gcp-plain-session-credentials)"}, []string{})
}

func TestCreateSession_FacadeAddSessionReturnsError(t *testing.T) {
	gcpPlainSessionActionsSetup()

	sessionName := "sessionName"
	accountId := "accountId"
	projectName := "projectName"
	oauthCode := "oauthCode"
	uuid := "uuid"
	credentials := "credentialsJson"
	envMock.ExpUuid = uuid
	gcpApiMock.ExpOauthToken = oauth2.Token{}
	gcpApiMock.ExpCredentials = credentials
	gcpPlainSessionFacadeMock.ExpErrorOnAddSession = true

	err := gcpPlainSessionActions.CreateSession(sessionName, accountId, projectName, oauthCode)
	test.ExpectHttpError(t, err, http.StatusConflict, "session already exist")
	gcpPlainSessionActionsVerifyExpectedCalls(t, []string{"GetOauthUrl(oauthCode)", "GetCredentials()"}, []string{"GenerateUuid()"}, []string{"SetSecret(credentialsJson, uuid-gcp-plain-session-credentials)"}, []string{"AddSession(sessionName)"})
}

func TestStartSession_NoPreviousActiveSession(t *testing.T) {
	gcpPlainSessionActionsSetup()
	gcpPlainSessionFacadeMock.ExpGetSessions = []session.GcpPlainSession{{Id: "ID2", Status: session.NotActive}}
	sessionId := "ID1"

	err := gcpPlainSessionActions.StartSession(sessionId)
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpPlainSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{}, []string{"GetSessions()", "SetSessionStatus(ID1, 2)"})
}

func TestStartSession_PreviousActiveSessionDiffersFromNewActiveSession(t *testing.T) {
	gcpPlainSessionActionsSetup()
	gcpPlainSessionFacadeMock.ExpGetSessions = []session.GcpPlainSession{{Id: "ID2", Status: session.Active}}
	sessionId := "ID1"

	err := gcpPlainSessionActions.StartSession(sessionId)
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpPlainSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{}, []string{"GetSessions()", "SetSessionStatus(ID2, 0)", "SetSessionStatus(ID1, 2)"})
}

func TestStartSession_SessionWasAlreadyActive(t *testing.T) {
	gcpPlainSessionActionsSetup()
	gcpPlainSessionFacadeMock.ExpGetSessions = []session.GcpPlainSession{{Id: "ID1", Status: session.Active}}
	sessionId := "ID1"

	err := gcpPlainSessionActions.StartSession(sessionId)
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpPlainSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{}, []string{"GetSessions()", "SetSessionStatus(ID1, 2)"})
}

func TestStartSession_PreviousActiveSessionDifferentAndFacadeSetSessionStatusReturnsError(t *testing.T) {
	gcpPlainSessionActionsSetup()
	gcpPlainSessionFacadeMock.ExpGetSessions = []session.GcpPlainSession{{Id: "ID2", Status: session.Active}}
	gcpPlainSessionFacadeMock.ExpErrorOnSetSessionStatus = true
	sessionId := "ID1"

	err := gcpPlainSessionActions.StartSession(sessionId)
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "unable to set the session status")
	gcpPlainSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{}, []string{"GetSessions()", "SetSessionStatus(ID2, 0)"})
}

func TestStartSession_FacadeSetSessionStatusReturnsError(t *testing.T) {
	gcpPlainSessionActionsSetup()
	gcpPlainSessionFacadeMock.ExpGetSessions = []session.GcpPlainSession{{Id: "ID2", Status: session.NotActive}}
	gcpPlainSessionFacadeMock.ExpErrorOnSetSessionStatus = true
	sessionId := "ID1"

	err := gcpPlainSessionActions.StartSession(sessionId)
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "unable to set the session status")
	gcpPlainSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{}, []string{"GetSessions()", "SetSessionStatus(ID1, 2)"})
}

func TestStopSession(t *testing.T) {
	gcpPlainSessionActionsSetup()
	sessionId := "ID"
	err := gcpPlainSessionActions.StopSession(sessionId)
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpPlainSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{}, []string{"SetSessionStatus(ID, 0)"})
}

func TestStopSession_FacadeReturnsError(t *testing.T) {
	gcpPlainSessionActionsSetup()
	gcpPlainSessionFacadeMock.ExpErrorOnSetSessionStatus = true
	sessionId := "ID"
	err := gcpPlainSessionActions.StopSession(sessionId)
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "unable to set the session status")
	gcpPlainSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{}, []string{"SetSessionStatus(ID, 0)"})
}

func TestDeleteSession(t *testing.T) {
	gcpPlainSessionActionsSetup()
	sessionId := "ID"
	credentialsLabel := "credentialLabel"
	gcpPlainSessionFacadeMock.ExpGetSessionById = session.GcpPlainSession{Id: "ID", CredentialsLabel: credentialsLabel}

	err := gcpPlainSessionActions.DeleteSession(sessionId)
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpPlainSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{"DeleteSecret(credentialLabel)"}, []string{"GetSessionById(ID)", "RemoveSession(ID)"})
}

func TestDeleteSession_FacadeGetSessionByIdReturnsError(t *testing.T) {
	gcpPlainSessionActionsSetup()
	sessionId := "ID"
	gcpPlainSessionFacadeMock.ExpErrorOnGetSessionById = true

	err := gcpPlainSessionActions.DeleteSession(sessionId)
	test.ExpectHttpError(t, err, http.StatusNotFound, "session not found")
	gcpPlainSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{}, []string{"GetSessionById(ID)"})
}

func TestDeleteSession_KeychainDeleteSecretReturnsError(t *testing.T) {
	gcpPlainSessionActionsSetup()
	sessionId := "ID"
	credentialsLabel := "credentialLabel"
	gcpPlainSessionFacadeMock.ExpGetSessionById = session.GcpPlainSession{Id: "ID", CredentialsLabel: credentialsLabel}
	gcpPlainSessionActionsKeychainMock.ExpErrorOnDeleteSecret = true

	err := gcpPlainSessionActions.DeleteSession(sessionId)
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpPlainSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{"DeleteSecret(credentialLabel)"}, []string{"GetSessionById(ID)", "RemoveSession(ID)"})
}

func TestDeleteSession_FacadeRemoveSessionReturnsError(t *testing.T) {
	gcpPlainSessionActionsSetup()
	sessionId := "ID"
	credentialsLabel := "credentialLabel"
	gcpPlainSessionFacadeMock.ExpGetSessionById = session.GcpPlainSession{Id: "ID", CredentialsLabel: credentialsLabel}
	gcpPlainSessionFacadeMock.ExpErrorOnRemoveSession = true

	err := gcpPlainSessionActions.DeleteSession(sessionId)
	test.ExpectHttpError(t, err, http.StatusNotFound, "session not found")
	gcpPlainSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{"DeleteSecret(credentialLabel)"}, []string{"GetSessionById(ID)", "RemoveSession(ID)"})
}

func TestEditSession(t *testing.T) {
	gcpPlainSessionActionsSetup()

	sessionId := "ID"
	sessionName := "sessionName"
	projectName := "projectName"

	err := gcpPlainSessionActions.EditSession(sessionId, sessionName, projectName)
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpPlainSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{}, []string{"EditSession(ID, sessionName, projectName)"})
}

func TestEditSession_FacadeEditSessionReturnsError(t *testing.T) {
	gcpPlainSessionActionsSetup()

	sessionId := "ID"
	sessionName := "sessionName"
	projectName := "projectName"
	gcpPlainSessionFacadeMock.ExpErrorOnEditSession = true

	err := gcpPlainSessionActions.EditSession(sessionId, sessionName, projectName)
	test.ExpectHttpError(t, err, http.StatusConflict, "unable to edit session, collision detected")

	gcpPlainSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{}, []string{"EditSession(ID, sessionName, projectName)"})
}
