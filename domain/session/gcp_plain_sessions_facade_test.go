package session

import (
	"leapp_daemon/test"
	"net/http"
	"reflect"
	"testing"
)

var (
	facade               *gcpPlainSessionsFacade
	sessionsBeforeUpdate []GcpPlainSession
	sessionsAfterUpdate  []GcpPlainSession
)

func setup() {
	facade = &gcpPlainSessionsFacade{
		gcpPlainSessions: make([]GcpPlainSession, 0),
	}
	sessionsBeforeUpdate = []GcpPlainSession{}
	sessionsAfterUpdate = []GcpPlainSession{}
}

func TestSingleton(t *testing.T) {
	facadeInstance1 := GetGcpPlainSessionFacade()
	facadeInstance2 := GetGcpPlainSessionFacade()

	if facadeInstance1 != facadeInstance2 {
		t.Fatalf("singleton is not returning the same instance")
	}
}

func TestGetSessions(t *testing.T) {
	setup()

	newSessions := []GcpPlainSession{{Id: "id"}}
	facade.gcpPlainSessions = newSessions

	if !reflect.DeepEqual(facade.GetSessions(), newSessions) {
		t.Errorf("unexpected sessions")
	}
}

func TestSetSessions(t *testing.T) {
	setup()

	newSessions := []GcpPlainSession{{Id: "id"}}
	facade.SetSessions(newSessions)

	if !reflect.DeepEqual(facade.gcpPlainSessions, newSessions) {
		t.Errorf("unexpected sessions")
	}
}

func TestAddSession(t *testing.T) {
	setup()
	facade.Subscribe(FakeObserver{})

	newSession := GcpPlainSession{Id: "id"}
	facade.AddSession(newSession)

	if !reflect.DeepEqual(sessionsBeforeUpdate, []GcpPlainSession{}) {
		t.Errorf("sessions were not empty")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []GcpPlainSession{newSession}) {
		t.Errorf("unexpected session")
	}
}

func TestAddSession_alreadyExistentId(t *testing.T) {
	setup()
	facade.Subscribe(FakeObserver{})

	newSession := GcpPlainSession{Id: "ID"}
	facade.gcpPlainSessions = []GcpPlainSession{newSession}

	err := facade.AddSession(newSession)
	test.ExpectHttpError(t, err, http.StatusConflict, "a session with id ID is already present")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAddSession_alreadyExistentName(t *testing.T) {
	setup()
	facade.Subscribe(FakeObserver{})

	facade.gcpPlainSessions = []GcpPlainSession{{Id: "1", Name: "NAME"}}

	err := facade.AddSession(GcpPlainSession{Id: "2", Name: "NAME"})
	test.ExpectHttpError(t, err, http.StatusConflict, "a session named NAME is already present")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestRemoveSession(t *testing.T) {
	setup()
	facade.Subscribe(FakeObserver{})

	newSession := GcpPlainSession{Id: "ID"}
	facade.gcpPlainSessions = []GcpPlainSession{newSession}

	facade.RemoveSession("ID")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []GcpPlainSession{newSession}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []GcpPlainSession{}) {
		t.Errorf("sessions were not empty")
	}
}

func TestRemoveSession_notFound(t *testing.T) {
	setup()
	facade.Subscribe(FakeObserver{})

	err := facade.RemoveSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "plain gcp session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestSetSessionStatus(t *testing.T) {
	setup()
	facade.Subscribe(FakeObserver{})

	newSession := GcpPlainSession{Id: "ID", Status: NotActive}
	facade.gcpPlainSessions = []GcpPlainSession{newSession}

	facade.SetSessionStatus("ID", Pending)

	if !reflect.DeepEqual(sessionsBeforeUpdate, []GcpPlainSession{newSession}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []GcpPlainSession{{Id: "ID", Status: Pending}}) {
		t.Errorf("sessions were not updated")
	}
}

func TestSetSessionStatus_notFound(t *testing.T) {
	setup()
	facade.Subscribe(FakeObserver{})

	err := facade.SetSessionStatus("ID", Pending)
	test.ExpectHttpError(t, err, http.StatusNotFound, "plain gcp session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

type FakeObserver struct {
}

func (observer FakeObserver) UpdateGcpPlainSessions(oldSessions []GcpPlainSession, newSessions []GcpPlainSession) error {
	sessionsBeforeUpdate = oldSessions
	sessionsAfterUpdate = newSessions
	return nil
}
