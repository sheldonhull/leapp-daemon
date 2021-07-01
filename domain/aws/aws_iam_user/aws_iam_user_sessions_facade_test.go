package aws_iam_user

import (
	"leapp_daemon/domain/aws"
	"leapp_daemon/test"
	"net/http"
	"reflect"
	"testing"
)

var (
	facade               *AwsIamUserSessionsFacade
	sessionsBeforeUpdate []AwsIamUserSession
	sessionsAfterUpdate  []AwsIamUserSession
)

func facadeSetup() {
	facade = NewAwsIamUserSessionsFacade()
	sessionsBeforeUpdate = []AwsIamUserSession{}
	sessionsAfterUpdate = []AwsIamUserSession{}
}

func TestAwsIamUserSessionsFacade_GetSessions(t *testing.T) {
	facadeSetup()

	newSessions := []AwsIamUserSession{{Id: "id"}}
	facade.awsIamUserSessions = newSessions

	if !reflect.DeepEqual(facade.GetSessions(), newSessions) {
		t.Errorf("unexpected sessions")
	}
}

func TestAwsIamUserSessionsFacade_SetSessions(t *testing.T) {
	facadeSetup()

	newSessions := []AwsIamUserSession{{Id: "id"}}
	facade.SetSessions(newSessions)

	if !reflect.DeepEqual(facade.awsIamUserSessions, newSessions) {
		t.Errorf("unexpected sessions")
	}
}

func TestAwsIamUserSessionsFacade_AddSession(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	newSession := AwsIamUserSession{Id: "id"}
	facade.AddSession(newSession)

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AwsIamUserSession{}) {
		t.Errorf("sessions were not empty")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AwsIamUserSession{newSession}) {
		t.Errorf("unexpected session")
	}
}

func TestAwsIamUserSessionsFacade_AddSession_alreadyExistentId(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	newSession := AwsIamUserSession{Id: "ID"}
	facade.awsIamUserSessions = []AwsIamUserSession{newSession}

	err := facade.AddSession(newSession)
	test.ExpectHttpError(t, err, http.StatusConflict, "a session with id ID is already present")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAwsIamUserSessionsFacade_AddSession_alreadyExistentName(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	facade.awsIamUserSessions = []AwsIamUserSession{{Id: "1", Name: "NAME"}}

	err := facade.AddSession(AwsIamUserSession{Id: "2", Name: "NAME"})
	test.ExpectHttpError(t, err, http.StatusConflict, "a session named NAME is already present")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAwsIamUserSessionsFacade_RemoveSession(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	session1 := AwsIamUserSession{Id: "ID1"}
	session2 := AwsIamUserSession{Id: "ID2"}
	facade.awsIamUserSessions = []AwsIamUserSession{session1, session2}

	facade.RemoveSession("ID1")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AwsIamUserSession{session1, session2}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AwsIamUserSession{session2}) {
		t.Errorf("sessions were not empty")
	}
}

func TestAwsIamUserSessionsFacade_RemoveSession_notFound(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	err := facade.RemoveSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAwsIamUserSessionsFacade_SetSessionTokenExpiration(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	newSession := AwsIamUserSession{Id: "ID", SessionTokenExpiration: "sessionTokenExpiration"}
	facade.awsIamUserSessions = []AwsIamUserSession{newSession}

	facade.SetSessionTokenExpiration("ID", "newSessionTokenExpiration")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AwsIamUserSession{newSession}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AwsIamUserSession{{Id: "ID", SessionTokenExpiration: "newSessionTokenExpiration"}}) {
		t.Errorf("sessions were not updated")
	}
}

func TestAwsIamUserSessionsFacade_SetSessionTokenExpiration_notFound(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	err := facade.SetSessionTokenExpiration("ID", "newSessionTokenExpiration")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAwsIamUserSessionsFacade_StartingSession(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	newSession := AwsIamUserSession{Id: "ID", Status: aws.NotActive}
	facade.awsIamUserSessions = []AwsIamUserSession{newSession}

	facade.StartingSession("ID")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AwsIamUserSession{newSession}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AwsIamUserSession{{Id: "ID", Status: aws.Pending, StartTime: "", LastStopTime: ""}}) {
		t.Errorf("sessions were not updated")
	}
}

func TestAwsIamUserSessionsFacade_StartingSession_notFound(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	err := facade.StartingSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAwsIamUserSessionsFacade_StartSession(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	newSession := AwsIamUserSession{Id: "ID", Status: aws.NotActive}
	facade.awsIamUserSessions = []AwsIamUserSession{newSession}

	facade.StartSession("ID", "start-time")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AwsIamUserSession{newSession}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AwsIamUserSession{{Id: "ID", Status: aws.Active, StartTime: "start-time"}}) {
		t.Errorf("sessions were not updated")
	}
}

func TestAwsIamUserSessionsFacade_StartSession_notFound(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	err := facade.StartSession("ID", "start-time")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAwsIamUserSessionsFacade_StopSession(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	newSession := AwsIamUserSession{Id: "ID", Status: aws.Active}
	facade.awsIamUserSessions = []AwsIamUserSession{newSession}

	facade.StopSession("ID", "stop-time")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AwsIamUserSession{newSession}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AwsIamUserSession{{Id: "ID", Status: aws.NotActive, LastStopTime: "stop-time"}}) {
		t.Errorf("sessions were not updated")
	}
}

func TestAwsIamUserSessionsFacade_StopSession_notFound(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	err := facade.StopSession("ID", "stop-time")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAwsIamUserSessionsFacade_EditSession(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	session1 := AwsIamUserSession{Id: "ID1", Name: "Name1", Region: "region", AccessKeyIdLabel: "accessKeyIdLabel",
		SecretKeyLabel: "secretKeyLabel", SessionTokenLabel: "sessionTokenLabel", MfaDevice: "mfaDevice",
		SessionTokenExpiration: "sessionTokenExpiration", NamedProfileId: "ProfileId1"}
	session2 := AwsIamUserSession{Id: "ID2", Name: "Name2", Region: "region2", AccessKeyIdLabel: "accessKeyIdLabel2",
		SecretKeyLabel: "secretKeyLabel2", SessionTokenLabel: "sessionTokenLabel2", MfaDevice: "mfaDevice2",
		SessionTokenExpiration: "sessionTokenExpiration2", NamedProfileId: "ProfileId2"}
	facade.awsIamUserSessions = []AwsIamUserSession{session1, session2}

	facade.EditSession("ID1", "newName", "newRegion", "newAccountNumber",
		"newUserName", "newMfaDevice", "newProfileId")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AwsIamUserSession{session1, session2}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AwsIamUserSession{
		{Id: "ID1", Name: "newName", Region: "newRegion", AccountNumber: "newAccountNumber", UserName: "newUserName",
			AccessKeyIdLabel: "accessKeyIdLabel", SecretKeyLabel: "secretKeyLabel", SessionTokenLabel: "sessionTokenLabel",
			MfaDevice: "newMfaDevice", NamedProfileId: "newProfileId", Status: aws.NotActive, StartTime: "",
			LastStopTime: "", SessionTokenExpiration: ""}, session2}) {
		t.Errorf("sessions were not updated")
	}
}

func TestAwsIamUserSessionsFacade_EditSession_DuplicateSessionNameAttempt(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	session1 := AwsIamUserSession{Id: "ID1", Name: "Name1", Region: "region", AccessKeyIdLabel: "accessKeyIdLabel",
		SecretKeyLabel: "secretKeyLabel", SessionTokenLabel: "sessionTokenLabel", MfaDevice: "mfaDevice",
		SessionTokenExpiration: "sessionTokenExpiration", NamedProfileId: "ProfileId1"}
	session2 := AwsIamUserSession{Id: "ID2", Name: "Name2", Region: "region2", AccessKeyIdLabel: "accessKeyIdLabel2",
		SecretKeyLabel: "secretKeyLabel2", SessionTokenLabel: "sessionTokenLabel2", MfaDevice: "mfaDevice2",
		SessionTokenExpiration: "sessionTokenExpiration2", NamedProfileId: "ProfileId2"}
	facade.awsIamUserSessions = []AwsIamUserSession{session1, session2}

	err := facade.EditSession("ID1", "Name2", "newRegion", "newAccountNumber",
		"newUserName", "newMfaDevice", "newProfileId")
	test.ExpectHttpError(t, err, http.StatusConflict, "a session named Name2 is already present")

	err = facade.EditSession("ID2", "Name1", "newRegion", "newAccountNumber",
		"newUserName", "newMfaDevice", "newProfileId")
	test.ExpectHttpError(t, err, http.StatusConflict, "a session named Name1 is already present")
}

func TestAwsIamUserSessionsFacade_EditSession_notFound(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	err := facade.EditSession("ID", "Name2", "NewRegion", "newAccountNumber",
		"newUserName", "newMfaDevice", "newProfileId")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

type fakeSessionsObserver struct {
}

func (observer fakeSessionsObserver) UpdateAwsIamUserSessions(oldSessions []AwsIamUserSession, newSessions []AwsIamUserSession) {
	sessionsBeforeUpdate = oldSessions
	sessionsAfterUpdate = newSessions
}
