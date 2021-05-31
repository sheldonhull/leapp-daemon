package repository

import (
  "leapp_daemon/infrastructure/http/http_error"
  "leapp_daemon/mock"
  "reflect"
  "testing"
)

func TestGcloudConfigurationRepository_onWindows(t *testing.T) {
  fsMock := mock.NewFileSystemMock(false)
  envMock := mock.NewEnvironmentMock(true)
  repo := &GcloudConfigurationRepository{
    FileSystem:  &fsMock,
    Environment: &envMock,
  }

  configDir, err := repo.configDir()
  if err != nil {
    t.Fatalf("unexpected error: %v", err)
  }

  expectedConfigDir := "c:/appdata/gcloud"
  if configDir != expectedConfigDir {
    t.Fatalf("expected configDir: %v", expectedConfigDir)
  }

  if !reflect.DeepEqual(fsMock.GetCalls(), []string{}) ||
    !reflect.DeepEqual(envMock.GetCalls(), []string{"IsWindows()", "GetEnvironmentVariable(\"APPDATA\")"}) {
    t.Fatalf("mock expectation violation.\nfsMock calls: %v\nenvMock calls: %v", fsMock.GetCalls(), envMock.GetCalls())
  }
}

func TestGcloudConfigurationRepository_onUnix(t *testing.T) {
  fsMock := mock.NewFileSystemMock(false)
  envMock := mock.NewEnvironmentMock(false)
  repo := &GcloudConfigurationRepository{
    FileSystem:  &fsMock,
    Environment: &envMock,
  }

  configDir, err := repo.configDir()
  if err != nil {
    t.Fatalf("unexpected error: %v", err)
  }

  expectedConfigDir := "/user/home/.config/gcloud"
  if configDir != expectedConfigDir {
    t.Fatalf("expected configDir: %v", expectedConfigDir)
  }

  if !reflect.DeepEqual(fsMock.GetCalls(), []string{"GetHomeDir()"}) ||
    !reflect.DeepEqual(envMock.GetCalls(), []string{"IsWindows()"}) {
    t.Fatalf("mock expectation violation.\nfsMock calls: %v\nenvMock calls: %v", fsMock.GetCalls(), envMock.GetCalls())
  }
}

func TestGcloudConfigurationRepository_error(t *testing.T) {
  fsMock := mock.NewFileSystemMock(true)
  envMock := mock.NewEnvironmentMock(false)
  repo := &GcloudConfigurationRepository{
    FileSystem:  &fsMock,
    Environment: &envMock,
  }

  expectedError := "error"
  _, err := repo.configDir()
  customError, isCustomError := err.(http_error.CustomError)
  if !isCustomError {
    t.Fatalf("expected CustomError")
  }
  if customError.StatusCode != 500 || customError.Error() != expectedError {
    t.Fatalf("unexpected error")
  }
}
