package configuration

import (
  "leapp_daemon/domain/session"
  "reflect"
  "testing"
)

var gcpSession = session.PlainGcpSession{Id: "123"}

func TestConfiguration_AddPlainGcpSession(t *testing.T) {
  config := Configuration{PlainGcpSessions: []session.PlainGcpSession{gcpSession}}
  newGcpSession := session.PlainGcpSession{Id: "234"}
  if err := config.AddPlainGcpSession(newGcpSession); err != nil {
    t.Errorf("unexpected error: %v", err)
  }
  if !reflect.DeepEqual(config.PlainGcpSessions, []session.PlainGcpSession{gcpSession, newGcpSession}) {
    t.Errorf("session not added")
  }
}

func TestConfiguration_AddPlainGcpSession_ConfigAlreadyPresent(t *testing.T) {
  config := Configuration{PlainGcpSessions: []session.PlainGcpSession{gcpSession}}
  expectedError := "a PlainGcpSession with id 123 is already present"
  if err := config.AddPlainGcpSession(gcpSession); expectedError != err.Error() {
    t.Fatalf("expected error: %v", expectedError)
  }
}

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
