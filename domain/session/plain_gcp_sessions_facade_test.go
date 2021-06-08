package session

import (
	"leapp_daemon/test"
	"net/http"
	"reflect"
	"testing"
)

var oldSessions []PlainGcpSession
var newSessions []PlainGcpSession

func TestSingleton(t *testing.T) {
	facadeInstance1 := GetPlainGcpSessionsFacade()
	facadeInstance2 := GetPlainGcpSessionsFacade()

	if facadeInstance1 != facadeInstance2 {
		t.Fatalf("singleton is not returning the same instance")
	}
}

type FakeObserver struct {
}

func (observer FakeObserver) UpdatePlainGcpSessions(oldPlainGcpSessions []PlainGcpSession, newPlainGcpSessions []PlainGcpSession) error {
	oldSessions = oldPlainGcpSessions
	newSessions = newPlainGcpSessions
	return nil
}

func TestGetSessions(t *testing.T) {
	facade := GetPlainGcpSessionsFacade()

	newSessions := []PlainGcpSession{{Id: "id"}}
	facade.plainGcpSessions = newSessions

	if !reflect.DeepEqual(facade.GetSessions(), newSessions) {
		t.Errorf("unexpected sessions")
	}
}

func TestSetSessions(t *testing.T) {
	facade := GetPlainGcpSessionsFacade()

	newSessions := []PlainGcpSession{{Id: "id"}}
	facade.SetSessions(newSessions)

	if !reflect.DeepEqual(facade.plainGcpSessions, newSessions) {
		t.Errorf("unexpected sessions")
	}
}

func TestAddSession(t *testing.T) {
	facade := GetPlainGcpSessionsFacade()
	facade.Subscribe(FakeObserver{})

	newSession := PlainGcpSession{Id: "id"}
	facade.AddSession(newSession)

	if !reflect.DeepEqual(oldSessions, []PlainGcpSession{}) {
		t.Errorf("sessions were not empty")
	}

	if !reflect.DeepEqual(newSessions, []PlainGcpSession{newSession}) {
		t.Errorf("unexpected session")
	}
}

func TestAddSession_alreadyExistent(t *testing.T) {
	facade := GetPlainGcpSessionsFacade()
	facade.Subscribe(FakeObserver{})

	newSession := PlainGcpSession{Id: "ID"}
	facade.plainGcpSessions = []PlainGcpSession{newSession}

	err := facade.AddSession(newSession)
	test.ExpectHttpError(t, err, http.StatusConflict, "a PlainGcpSession with id ID is already present")

	if oldSessions != nil || newSessions != nil {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestRemoveSession(t *testing.T) {
	facade := GetPlainGcpSessionsFacade()
	facade.Subscribe(FakeObserver{})

	newSession := PlainGcpSession{Id: "ID"}
	facade.plainGcpSessions = []PlainGcpSession{newSession}

	facade.RemoveSession("ID")

	if !reflect.DeepEqual(oldSessions, []PlainGcpSession{newSession}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(newSessions, []PlainGcpSession{}) {
		t.Errorf("sessions were not empty")
	}
}

func TestRemoveSession_notFound(t *testing.T) {
	facade := GetPlainGcpSessionsFacade()
	facade.Subscribe(FakeObserver{})

	err := facade.RemoveSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "plain gcp session with id ID not found")

	if oldSessions != nil || newSessions != nil {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestSetSessionStatus(t *testing.T) {
	facade := GetPlainGcpSessionsFacade()
	facade.Subscribe(FakeObserver{})

	newSession := PlainGcpSession{Id: "ID", Status: NotActive}
	facade.plainGcpSessions = []PlainGcpSession{newSession}

	facade.SetSessionStatus("ID", Pending)

	if !reflect.DeepEqual(oldSessions, []PlainGcpSession{newSession}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(newSessions, []PlainGcpSession{{Id: "ID", Status: Pending}}) {
		t.Errorf("sessions were not updated")
	}
}

func TestSetSessionStatus_notFound(t *testing.T) {
	facade := GetPlainGcpSessionsFacade()
	facade.Subscribe(FakeObserver{})

	err := facade.SetSessionStatus("ID", Pending)
	test.ExpectHttpError(t, err, http.StatusNotFound, "plain gcp session with id ID not found")

	if oldSessions != nil || newSessions != nil {
		t.Errorf("sessions was unexpectedly changed")
	}
}

/*

func TestConfiguration_GetAllPlainGcpSessions(t *testing.T) {
  expectedSessions := []session.PlainGcpSession{gcpSession}
  config := Configuration{PlainGcpSessions: expectedSessions}
  if sessions, err := config.GetAllPlainGcpSessions(); err != nil {
    t.Fatalf("unexpected error: %v", err)
  } else if !reflect.DeepEqual(sessions, expectedSessions) {
    t.Fatalf("expected sessions: %v", expectedSessions)
  }
}

func TestConfiguration_RemovePlainGcpSession(t *testing.T) {
  config := Configuration{PlainGcpSessions: []session.PlainGcpSession{gcpSession}}
  if err := config.RemovePlainGcpSession(gcpSession); err != nil {
    t.Errorf("unexpected error: %v", err)
  }
  if !reflect.DeepEqual(config.PlainGcpSessions, []session.PlainGcpSession{}) {
    t.Errorf("session not removed")
  }
}

func TestConfiguration_RemovePlainGcpSession_NoConfigs(t *testing.T) {
  config := Configuration{PlainGcpSessions: []session.PlainGcpSession{}}
  expectedError := "PlainGcpSession with id 123 not found"
  if err := config.RemovePlainGcpSession(gcpSession); err.Error() != expectedError {
    t.Errorf("expected error: %v", expectedError)
  }
}
*/
