package repository

import (
	"fmt"
	"leapp_daemon/infrastructure/http/http_error"
	"leapp_daemon/mock"
	"path/filepath"
	"reflect"
	"testing"
)

func TestIsGcloudCliAvailable_available(t *testing.T) {
	envMock := mock.NewEnvironmentMock()
	envMock.ExpIsCommandAvailable = true
	repo := &GcloudConfigurationRepository{
		Environment: &envMock,
	}

	if !repo.IsGcloudCliAvailable() {
		t.Fatalf("should be available")
	}
}

func TestIsGcloudCliAvailable_unavailable(t *testing.T) {
	envMock := mock.NewEnvironmentMock()
	repo := &GcloudConfigurationRepository{
		Environment: &envMock,
	}

	if repo.IsGcloudCliAvailable() {
		t.Fatalf("should be unavailable")
	}
}

func TestDoesGcloudConfigFolderExist_exists(t *testing.T) {
	envMock := mock.NewEnvironmentMock()
	fsMock := mock.NewFileSystemMock()
	fsMock.ExpDoesFileExist = true
	envMock.ExpIsWindows = true

	repo := &GcloudConfigurationRepository{
		Environment: &envMock,
		FileSystem:  &fsMock,
	}

	exists, err := repo.DoesGcloudConfigFolderExist()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !exists {
		t.Fatalf("directory should exist")
	}

	expectedConfigDirPath := filepath.Join("c:/", "appdata", "gcloud")
	if !reflect.DeepEqual(fsMock.GetCalls(), []string{fmt.Sprintf("DoesFileExist(%v)", expectedConfigDirPath)}) {
		t.Fatalf("mock expectation violation.\nfsMock calls: %v", fsMock.GetCalls())
	}
}

func TestDoesGcloudConfigFolderExist_not_exists(t *testing.T) {
	envMock := mock.NewEnvironmentMock()
	fsMock := mock.NewFileSystemMock()
	envMock.ExpIsWindows = true

	repo := &GcloudConfigurationRepository{
		Environment: &envMock,
		FileSystem:  &fsMock,
	}

	exists, err := repo.DoesGcloudConfigFolderExist()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if exists {
		t.Fatalf("directory should not exist")
	}

	expectedConfigDirPath := filepath.Join("c:/", "appdata", "gcloud")
	if !reflect.DeepEqual(fsMock.GetCalls(), []string{fmt.Sprintf("DoesFileExist(%v)", expectedConfigDirPath)}) {
		t.Fatalf("mock expectation violation.\nfsMock calls: %v", fsMock.GetCalls())
	}
}

func TestGcloudConfigurationRepository_onWindows(t *testing.T) {
	envMock := mock.NewEnvironmentMock()
	envMock.ExpIsWindows = true
	repo := &GcloudConfigurationRepository{
		Environment: &envMock,
	}

	configDir, err := repo.configDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedConfigDir := filepath.Join("c:/", "appdata", "gcloud")
	if configDir != expectedConfigDir {
		t.Fatalf("expected configDir: %v", expectedConfigDir)
	}

	if !reflect.DeepEqual(envMock.GetCalls(), []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"}) {
		t.Fatalf("mock expectation violation.\nenvMock calls: %v", envMock.GetCalls())
	}
}

func TestGcloudConfigurationRepository_onUnix(t *testing.T) {
	fsMock := mock.NewFileSystemMock()
	envMock := mock.NewEnvironmentMock()
	repo := &GcloudConfigurationRepository{
		FileSystem:  &fsMock,
		Environment: &envMock,
	}

	configDir, err := repo.configDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedConfigDir := filepath.Join("/", "user", "home", ".config", "gcloud")
	if configDir != expectedConfigDir {
		t.Fatalf("expected configDir: %v", expectedConfigDir)
	}

	if !reflect.DeepEqual(fsMock.GetCalls(), []string{"GetHomeDir()"}) ||
		!reflect.DeepEqual(envMock.GetCalls(), []string{"IsWindows()"}) {
		t.Fatalf("mock expectation violation.\nfsMock calls: %v\nenvMock calls: %v", fsMock.GetCalls(), envMock.GetCalls())
	}
}

func TestGcloudConfigurationRepository_error(t *testing.T) {
	fsMock := mock.NewFileSystemMock()
	fsMock.ExpErrorOnGetHomeDir = true
	envMock := mock.NewEnvironmentMock()
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
