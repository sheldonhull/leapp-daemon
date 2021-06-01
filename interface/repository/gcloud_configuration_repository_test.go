package repository

import (
	"fmt"
	"leapp_daemon/infrastructure/http/http_error"
	"leapp_daemon/mock"
	"path/filepath"
	"reflect"
	"testing"
)

var expectedConfigDirPath string
var expectedConfigFilePath string
var expectedActiveConfigFilePath string
var envMock mock.EnvironmentMock
var fsMock mock.FileSystemMock
var repo *GcloudConfigurationRepository

func setup() {
	expectedConfigDirPath = filepath.Join("c:/", "appdata", "gcloud")
	expectedConfigFilePath = filepath.Join(expectedConfigDirPath, "configurations", "config_leapp_configurationName")
	expectedActiveConfigFilePath = filepath.Join(expectedConfigDirPath, "active_config")

	envMock = mock.NewEnvironmentMock()
	fsMock = mock.NewFileSystemMock()
	repo = &GcloudConfigurationRepository{
		Environment: &envMock,
		FileSystem:  &fsMock,
	}
}

func verifyExpectedCalls(t *testing.T, envMockCalls []string, fsMockCalls []string) {
	if !reflect.DeepEqual(envMock.GetCalls(), envMockCalls) {
		t.Fatalf("envMock expectation violation.\nMock calls: %v", envMock.GetCalls())
	}
	if !reflect.DeepEqual(fsMock.GetCalls(), fsMockCalls) {
		t.Fatalf("fsMock expectation violation.\nMock calls: %v", fsMock.GetCalls())
	}
}

func expectHttpError(t *testing.T, err error, expectedStatusCode int, expectedError string) {
	customError, isCustomError := err.(http_error.CustomError)
	if !isCustomError {
		t.Fatalf("expected CustomError")
	}
	if customError.StatusCode != expectedStatusCode {
		t.Fatalf("unexpected error status code: %v", customError.StatusCode)
	}
	if customError.Error() != expectedError {
		t.Fatalf("unexpected error: %v", customError.Error())
	}
}

func TestIsGcloudCliAvailable_available(t *testing.T) {
	setup()
	envMock.ExpIsCommandAvailable = true

	if !repo.IsGcloudCliAvailable() {
		t.Fatalf("should be available")
	}
}

func TestIsGcloudCliAvailable_unavailable(t *testing.T) {
	setup()
	if repo.IsGcloudCliAvailable() {
		t.Fatalf("should be unavailable")
	}
}

func TestDoesGcloudConfigFolderExist_exists(t *testing.T) {
	setup()
	fsMock.ExpDoesFileExist = true
	envMock.ExpIsWindows = true

	exists, err := repo.DoesGcloudConfigFolderExist()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !exists {
		t.Fatalf("directory should exist")
	}

	verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
		[]string{fmt.Sprintf("DoesFileExist(%v)", expectedConfigDirPath)})
}

func TestDoesGcloudConfigFolderExist_not_exists(t *testing.T) {
	setup()
	envMock.ExpIsWindows = true

	exists, err := repo.DoesGcloudConfigFolderExist()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if exists {
		t.Fatalf("directory should not exist")
	}

	verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
		[]string{fmt.Sprintf("DoesFileExist(%v)", expectedConfigDirPath)})
}

func TestDoesGcloudConfigFolderExist_errorGettingHomeDir(t *testing.T) {
	setup()
	fsMock.ExpErrorOnGetHomeDir = true

	_, err := repo.DoesGcloudConfigFolderExist()
	expectHttpError(t, err, 500, "error getting home dir")

	verifyExpectedCalls(t, []string{"IsWindows()"}, []string{})
}

func TestGcloudConfigurationRepository_onWindows(t *testing.T) {
	setup()
	envMock.ExpIsWindows = true

	configDir, err := repo.gcloudConfigDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedConfigDir := filepath.Join("c:/", "appdata", "gcloud")
	if configDir != expectedConfigDir {
		t.Fatalf("expected gcloudConfigDir: %v", expectedConfigDir)
	}

	verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"}, []string{})
}

func TestGcloudConfigurationRepository_onUnix(t *testing.T) {
	setup()
	configDir, err := repo.gcloudConfigDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedConfigDir := filepath.Join("/", "user", "home", ".config", "gcloud")
	if configDir != expectedConfigDir {
		t.Fatalf("expected gcloudConfigDir: %v", expectedConfigDir)
	}

	verifyExpectedCalls(t, []string{"IsWindows()"}, []string{"GetHomeDir()"})
}

func TestGcloudConfigurationRepository_error(t *testing.T) {
	setup()
	fsMock.ExpErrorOnGetHomeDir = true

	_, err := repo.gcloudConfigDir()
	expectHttpError(t, err, 500, "error getting home dir")

	verifyExpectedCalls(t, []string{"IsWindows()"}, []string{})
}

func TestCreateConfiguration(t *testing.T) {
	setup()
	envMock.ExpIsWindows = true

	err := repo.CreateConfiguration("configurationName", "accountId", "projectId")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	expectedFileContent := "[core]\naccount = accountId\nproject = projectId\n"
	expectedFile := fmt.Sprintf("WriteToFile(%v, %v)", expectedConfigFilePath, []byte(expectedFileContent))

	verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"}, []string{expectedFile})
}

func TestCreateConfiguration_errorGettingHomeDir(t *testing.T) {
	setup()
	fsMock.ExpErrorOnGetHomeDir = true

	err := repo.CreateConfiguration("configurationName", "accountId", "projectId")
	expectHttpError(t, err, 500, "error getting home dir")

	verifyExpectedCalls(t, []string{"IsWindows()"}, []string{})
}

func TestCreateConfiguration_errorWritingFile(t *testing.T) {
	setup()
	fsMock.ExpErrorOnWriteToFile = true
	envMock.ExpIsWindows = true

	err := repo.CreateConfiguration("configurationName", "accountId", "projectId")
	expectHttpError(t, err, 500, "error writing file")

	verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"}, []string{})
}

func TestRemoveConfiguration(t *testing.T) {
	setup()
	envMock.ExpIsWindows = true

	err := repo.RemoveConfiguration("configurationName")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	expectedFile := fmt.Sprintf("RemoveFile(%v)", expectedConfigFilePath)
	verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"}, []string{expectedFile})
}

func TestRemoveConfiguration_errorGettingHomeDir(t *testing.T) {
	setup()
	fsMock.ExpErrorOnGetHomeDir = true

	err := repo.RemoveConfiguration("configurationName")
	expectHttpError(t, err, 500, "error getting home dir")

	verifyExpectedCalls(t, []string{"IsWindows()"}, []string{})
}

func TestRemoveConfiguration_errorRemovingFile(t *testing.T) {
	setup()
	fsMock.ExpErrorOnRemoveFile = true

	err := repo.RemoveConfiguration("configurationName")
	expectHttpError(t, err, 500, "error removing file")

	verifyExpectedCalls(t, []string{"IsWindows()"}, []string{"GetHomeDir()"})
}

func TestActivateConfiguration(t *testing.T) {
	setup()
	envMock.ExpIsWindows = true

	activeConfigurationName := "configurationName"
	err := repo.ActivateConfiguration(activeConfigurationName)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	expectedFile := fmt.Sprintf("WriteToFile(%v, %v)", expectedActiveConfigFilePath, []byte(activeConfigurationName))
	verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"}, []string{expectedFile})
}

func TestActivateConfiguration_errorGettingHomeDir(t *testing.T) {
	setup()
	fsMock.ExpErrorOnGetHomeDir = true

	err := repo.ActivateConfiguration("configurationName")
	expectHttpError(t, err, 500, "error getting home dir")

	verifyExpectedCalls(t, []string{"IsWindows()"}, []string{})
}

func TestActivateConfiguration_errorWritingFile(t *testing.T) {
	setup()
	fsMock.ExpErrorOnWriteToFile = true
	envMock.ExpIsWindows = true

	err := repo.ActivateConfiguration("configurationName")
	expectHttpError(t, err, 500, "error writing file")

	verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"}, []string{})
}

func TestDeactivateConfiguration(t *testing.T) {
	setup()
	envMock.ExpIsWindows = true

	err := repo.DeactivateConfiguration()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	expectedFile := fmt.Sprintf("RemoveFile(%v)", expectedActiveConfigFilePath)
	verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"}, []string{expectedFile})
}

func TestDeactivateConfiguration_errorGettingHomeDir(t *testing.T) {
	setup()
	fsMock.ExpErrorOnGetHomeDir = true

	err := repo.DeactivateConfiguration()
	expectHttpError(t, err, 500, "error getting home dir")

	verifyExpectedCalls(t, []string{"IsWindows()"}, []string{})
}

func TestDeactivateConfiguration_errorRemovingFile(t *testing.T) {
	setup()
	fsMock.ExpErrorOnRemoveFile = true

	err := repo.DeactivateConfiguration()
	expectHttpError(t, err, 500, "error removing file")

	verifyExpectedCalls(t, []string{"IsWindows()"}, []string{"GetHomeDir()"})
}
