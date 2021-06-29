package session

import (
	"leapp_daemon/test"
	"net/http"
	"reflect"
	"testing"
)

var (
	gcpIamUserSessionFacade        *GcpIamUserAccountOauthSessionsFacade
	gcpIamUserSessionsBeforeUpdate []GcpIamUserAccountOauthSession
	gcpIamUserSessionsAfterUpdate  []GcpIamUserAccountOauthSession
)

func gcpIamUserAccountOauthSessionFacadeSetup() {
	gcpIamUserSessionFacade = NewGcpIamUserAccountOauthSessionsFacade()
	gcpIamUserSessionsBeforeUpdate = []GcpIamUserAccountOauthSession{}
	gcpIamUserSessionsAfterUpdate = []GcpIamUserAccountOauthSession{}
}

func TestGcpIamUserSessionFacade_GetSessions(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()

	newSessions := []GcpIamUserAccountOauthSession{{Id: "id"}}
	gcpIamUserSessionFacade.gcpIamUserAccountOauthSessions = newSessions

	if !reflect.DeepEqual(gcpIamUserSessionFacade.GetSessions(), newSessions) {
		t.Errorf("unexpected sessions")
	}
}

func TestGcpIamUserSessionFacade_SetSessions(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()

	newSessions := []GcpIamUserAccountOauthSession{{Id: "id"}}
	gcpIamUserSessionFacade.SetSessions(newSessions)

	if !reflect.DeepEqual(gcpIamUserSessionFacade.gcpIamUserAccountOauthSessions, newSessions) {
		t.Errorf("unexpected sessions")
	}
}

func TestGcpIamUserSessionFacade_AddSession(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()
	gcpIamUserSessionFacade.Subscribe(FakeGcpIamUserSessionObserver{})

	newSession := GcpIamUserAccountOauthSession{Id: "id"}
	gcpIamUserSessionFacade.AddSession(newSession)

	if !reflect.DeepEqual(gcpIamUserSessionsBeforeUpdate, []GcpIamUserAccountOauthSession{}) {
		t.Errorf("sessions were not empty")
	}

	if !reflect.DeepEqual(gcpIamUserSessionsAfterUpdate, []GcpIamUserAccountOauthSession{newSession}) {
		t.Errorf("unexpected session")
	}
}

func TestGcpIamUserSessionFacade_AddSession_alreadyExistentId(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()
	gcpIamUserSessionFacade.Subscribe(FakeGcpIamUserSessionObserver{})

	newSession := GcpIamUserAccountOauthSession{Id: "ID"}
	gcpIamUserSessionFacade.gcpIamUserAccountOauthSessions = []GcpIamUserAccountOauthSession{newSession}

	err := gcpIamUserSessionFacade.AddSession(newSession)
	test.ExpectHttpError(t, err, http.StatusConflict, "a session with id ID is already present")

	if len(gcpIamUserSessionsBeforeUpdate) > 0 || len(gcpIamUserSessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestGcpIamUserSessionFacade_AddSession_alreadyExistentName(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()
	gcpIamUserSessionFacade.Subscribe(FakeGcpIamUserSessionObserver{})

	gcpIamUserSessionFacade.gcpIamUserAccountOauthSessions = []GcpIamUserAccountOauthSession{{Id: "1", Name: "NAME"}}

	err := gcpIamUserSessionFacade.AddSession(GcpIamUserAccountOauthSession{Id: "2", Name: "NAME"})
	test.ExpectHttpError(t, err, http.StatusConflict, "a session named NAME is already present")

	if len(gcpIamUserSessionsBeforeUpdate) > 0 || len(gcpIamUserSessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestGcpIamUserSessionFacade_RemoveSession(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()
	gcpIamUserSessionFacade.Subscribe(FakeGcpIamUserSessionObserver{})

	session1 := GcpIamUserAccountOauthSession{Id: "ID1"}
	session2 := GcpIamUserAccountOauthSession{Id: "ID2"}
	gcpIamUserSessionFacade.gcpIamUserAccountOauthSessions = []GcpIamUserAccountOauthSession{session1, session2}

	gcpIamUserSessionFacade.RemoveSession("ID1")

	if !reflect.DeepEqual(gcpIamUserSessionsBeforeUpdate, []GcpIamUserAccountOauthSession{session1, session2}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(gcpIamUserSessionsAfterUpdate, []GcpIamUserAccountOauthSession{session2}) {
		t.Errorf("sessions were not empty")
	}
}

func TestGcpIamUserSessionFacade_RemoveSession_notFound(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()
	gcpIamUserSessionFacade.Subscribe(FakeGcpIamUserSessionObserver{})

	err := gcpIamUserSessionFacade.RemoveSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(gcpIamUserSessionsBeforeUpdate) > 0 || len(gcpIamUserSessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestGcpIamUserSessionFacade_StartSession(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()
	gcpIamUserSessionFacade.Subscribe(FakeGcpIamUserSessionObserver{})

	newSession := GcpIamUserAccountOauthSession{Id: "ID", Status: NotActive}
	gcpIamUserSessionFacade.gcpIamUserAccountOauthSessions = []GcpIamUserAccountOauthSession{newSession}

	gcpIamUserSessionFacade.StartSession("ID", "start-time")

	if !reflect.DeepEqual(gcpIamUserSessionsBeforeUpdate, []GcpIamUserAccountOauthSession{newSession}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(gcpIamUserSessionsAfterUpdate, []GcpIamUserAccountOauthSession{{Id: "ID", Status: Active, StartTime: "start-time"}}) {
		t.Errorf("sessions were not updated")
	}
}

func TestGcpIamUserSessionFacade_StartSession_notFound(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()
	gcpIamUserSessionFacade.Subscribe(FakeGcpIamUserSessionObserver{})

	err := gcpIamUserSessionFacade.StartSession("ID", "start-time")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(gcpIamUserSessionsBeforeUpdate) > 0 || len(gcpIamUserSessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestGcpIamUserSessionFacade_StopSession(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()
	gcpIamUserSessionFacade.Subscribe(FakeGcpIamUserSessionObserver{})

	newSession := GcpIamUserAccountOauthSession{Id: "ID", Status: Active}
	gcpIamUserSessionFacade.gcpIamUserAccountOauthSessions = []GcpIamUserAccountOauthSession{newSession}

	gcpIamUserSessionFacade.StopSession("ID", "stop-time")

	if !reflect.DeepEqual(gcpIamUserSessionsBeforeUpdate, []GcpIamUserAccountOauthSession{newSession}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(gcpIamUserSessionsAfterUpdate, []GcpIamUserAccountOauthSession{{Id: "ID", Status: NotActive, LastStopTime: "stop-time"}}) {
		t.Errorf("sessions were not updated")
	}
}

func TestGcpIamUserSessionFacade_StopSession_notFound(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()
	gcpIamUserSessionFacade.Subscribe(FakeGcpIamUserSessionObserver{})

	err := gcpIamUserSessionFacade.StopSession("ID", "stop-time")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(gcpIamUserSessionsBeforeUpdate) > 0 || len(gcpIamUserSessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestGcpIamUserSessionFacade_EditSession(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()
	gcpIamUserSessionFacade.Subscribe(FakeGcpIamUserSessionObserver{})

	session1 := GcpIamUserAccountOauthSession{Id: "ID1", Name: "Name1", ProjectName: "Project1"}
	session2 := GcpIamUserAccountOauthSession{Id: "ID2", Name: "Name2", ProjectName: "Project2"}
	gcpIamUserSessionFacade.gcpIamUserAccountOauthSessions = []GcpIamUserAccountOauthSession{session1, session2}

	gcpIamUserSessionFacade.EditSession("ID1", "NewName", "NewProject")

	if !reflect.DeepEqual(gcpIamUserSessionsBeforeUpdate, []GcpIamUserAccountOauthSession{session1, session2}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(gcpIamUserSessionsAfterUpdate, []GcpIamUserAccountOauthSession{
		{Id: "ID1", Name: "NewName", ProjectName: "NewProject"}, session2}) {
		t.Errorf("sessions were not updated")
	}
}

func TestGcpIamUserSessionFacade_EditSession_DuplicateSessionNameAttempt(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()
	gcpIamUserSessionFacade.Subscribe(FakeGcpIamUserSessionObserver{})

	session1 := GcpIamUserAccountOauthSession{Id: "ID1", Name: "Name1", ProjectName: "Project1"}
	session2 := GcpIamUserAccountOauthSession{Id: "ID2", Name: "Name2", ProjectName: "Project2"}
	gcpIamUserSessionFacade.gcpIamUserAccountOauthSessions = []GcpIamUserAccountOauthSession{session1, session2}

	err := gcpIamUserSessionFacade.EditSession("ID1", "Name2", "NewProject")
	test.ExpectHttpError(t, err, http.StatusConflict, "a session named Name2 is already present")

	err = gcpIamUserSessionFacade.EditSession("ID2", "Name1", "NewProject")
	test.ExpectHttpError(t, err, http.StatusConflict, "a session named Name1 is already present")
}

func TestGcpIamUserSessionFacade_EditSession_notFound(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()
	gcpIamUserSessionFacade.Subscribe(FakeGcpIamUserSessionObserver{})

	err := gcpIamUserSessionFacade.EditSession("ID", "", "")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(gcpIamUserSessionsBeforeUpdate) > 0 || len(gcpIamUserSessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

type FakeGcpIamUserSessionObserver struct {
}

func (observer FakeGcpIamUserSessionObserver) UpdateGcpIamUserAccountOauthSessions(oldSessions []GcpIamUserAccountOauthSession, newSessions []GcpIamUserAccountOauthSession) {
	gcpIamUserSessionsBeforeUpdate = oldSessions
	gcpIamUserSessionsAfterUpdate = newSessions
}
