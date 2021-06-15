package repository

import (
  "fmt"
  "leapp_daemon/test"
  "leapp_daemon/test/mock"
  "path/filepath"
  "reflect"
  "testing"
)

var (
  expectedConfigDirPath              string
  expectedConfigFilePath             string
  expectedActiveConfigFilePath       string
  expectedDefaultCredentialsFilePath string
  expectedCredentialsDbFilePath      string
  envMock                            mock.EnvironmentMock
  fsMock                             mock.FileSystemMock
  dbTableMock                        mock.GcpCredentialsDbTableMock
  repo                               *GcloudConfigurationRepository
)

func setup() {
  expectedConfigDirPath = filepath.Join("c:/", "appdata", "gcloud")
  expectedConfigFilePath = filepath.Join(expectedConfigDirPath, "configurations", "config_leapp_configurationName")
  expectedActiveConfigFilePath = filepath.Join(expectedConfigDirPath, "active_config")
  expectedDefaultCredentialsFilePath = filepath.Join(expectedConfigDirPath, "application_default_credentials.json")
  expectedCredentialsDbFilePath = filepath.Join(expectedConfigDirPath, "credentials.db")

  envMock = mock.NewEnvironmentMock()
  fsMock = mock.NewFileSystemMock()
  dbTableMock = mock.NewGcpCredentialsDbTableMock()

  repo = &GcloudConfigurationRepository{
    Environment:      &envMock,
    FileSystem:       &fsMock,
    CredentialsTable: &dbTableMock,
  }
}

func verifyExpectedCalls(t *testing.T, envMockCalls []string, fsMockCalls []string, dbTableMockCalls []string) {
  if !reflect.DeepEqual(envMock.GetCalls(), envMockCalls) {
    t.Fatalf("envMock expectation violation.\nMock calls: %v", envMock.GetCalls())
  }
  if !reflect.DeepEqual(fsMock.GetCalls(), fsMockCalls) {
    t.Fatalf("fsMock expectation violation.\nMock calls: %v", fsMock.GetCalls())
  }
  if !reflect.DeepEqual(dbTableMock.GetCalls(), dbTableMockCalls) {
    t.Fatalf("dbTableMockCalls expectation violation.\nMock calls: %v", dbTableMock.GetCalls())
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
    []string{fmt.Sprintf("DoesFileExist(%v)", expectedConfigDirPath)}, []string{})
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
    []string{fmt.Sprintf("DoesFileExist(%v)", expectedConfigDirPath)}, []string{})
}

func TestDoesGcloudConfigFolderExist_errorGettingHomeDir(t *testing.T) {
  setup()
  fsMock.ExpErrorOnGetHomeDir = true

  _, err := repo.DoesGcloudConfigFolderExist()
  test.ExpectHttpError(t, err, 500, "error getting home dir")

  verifyExpectedCalls(t, []string{"IsWindows()"}, []string{}, []string{})
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

  verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"}, []string{}, []string{})
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

  verifyExpectedCalls(t, []string{"IsWindows()"}, []string{"GetHomeDir()"}, []string{})
}

func TestGcloudConfigurationRepository_error(t *testing.T) {
  setup()
  fsMock.ExpErrorOnGetHomeDir = true

  _, err := repo.gcloudConfigDir()
  test.ExpectHttpError(t, err, 500, "error getting home dir")

  verifyExpectedCalls(t, []string{"IsWindows()"}, []string{}, []string{})
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

  verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"}, []string{expectedFile}, []string{})
}

func TestCreateConfiguration_errorGettingHomeDir(t *testing.T) {
  setup()
  fsMock.ExpErrorOnGetHomeDir = true

  err := repo.CreateConfiguration("configurationName", "accountId", "projectId")
  test.ExpectHttpError(t, err, 500, "error getting home dir")

  verifyExpectedCalls(t, []string{"IsWindows()"}, []string{}, []string{})
}

func TestCreateConfiguration_errorWritingFile(t *testing.T) {
  setup()
  fsMock.ExpErrorOnWriteToFile = true
  envMock.ExpIsWindows = true

  err := repo.CreateConfiguration("configurationName", "accountId", "projectId")
  test.ExpectHttpError(t, err, 500, "error writing file")

  verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"}, []string{}, []string{})
}

func TestRemoveConfiguration(t *testing.T) {
  setup()
  envMock.ExpIsWindows = true

  err := repo.RemoveConfiguration("configurationName")
  if err != nil {
    t.Fatalf("unexpected error %v", err)
  }

  expectedFile := fmt.Sprintf("RemoveFile(%v)", expectedConfigFilePath)
  verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"}, []string{expectedFile}, []string{})
}

func TestRemoveConfiguration_errorGettingHomeDir(t *testing.T) {
  setup()
  fsMock.ExpErrorOnGetHomeDir = true

  err := repo.RemoveConfiguration("configurationName")
  test.ExpectHttpError(t, err, 500, "error getting home dir")

  verifyExpectedCalls(t, []string{"IsWindows()"}, []string{}, []string{})
}

func TestRemoveConfiguration_errorRemovingFile(t *testing.T) {
  setup()
  fsMock.ExpErrorOnRemoveFile = true

  err := repo.RemoveConfiguration("configurationName")
  test.ExpectHttpError(t, err, 500, "error removing file")

  verifyExpectedCalls(t, []string{"IsWindows()"}, []string{"GetHomeDir()"}, []string{})
}

func TestActivateConfiguration(t *testing.T) {
  setup()
  envMock.ExpIsWindows = true

  activeConfigurationName := "configurationName"
  err := repo.ActivateConfiguration(activeConfigurationName)
  if err != nil {
    t.Fatalf("unexpected error %v", err)
  }

  fullActiveConfigurationName := fmt.Sprintf("config_leapp_%v", activeConfigurationName)
  expectedFile := fmt.Sprintf("WriteToFile(%v, %v)", expectedActiveConfigFilePath, []byte(fullActiveConfigurationName))
  verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"}, []string{expectedFile}, []string{})
}

func TestActivateConfiguration_errorGettingHomeDir(t *testing.T) {
  setup()
  fsMock.ExpErrorOnGetHomeDir = true

  err := repo.ActivateConfiguration("configurationName")
  test.ExpectHttpError(t, err, 500, "error getting home dir")

  verifyExpectedCalls(t, []string{"IsWindows()"}, []string{}, []string{})
}

func TestActivateConfiguration_errorWritingFile(t *testing.T) {
  setup()
  fsMock.ExpErrorOnWriteToFile = true
  envMock.ExpIsWindows = true

  err := repo.ActivateConfiguration("configurationName")
  test.ExpectHttpError(t, err, 500, "error writing file")

  verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"}, []string{}, []string{})
}

func TestDeactivateConfiguration(t *testing.T) {
  setup()
  envMock.ExpIsWindows = true

  err := repo.DeactivateConfiguration()
  if err != nil {
    t.Fatalf("unexpected error %v", err)
  }

  expectedFile := fmt.Sprintf("RemoveFile(%v)", expectedActiveConfigFilePath)
  verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"}, []string{expectedFile}, []string{})
}

func TestDeactivateConfiguration_errorGettingHomeDir(t *testing.T) {
  setup()
  fsMock.ExpErrorOnGetHomeDir = true

  err := repo.DeactivateConfiguration()
  test.ExpectHttpError(t, err, 500, "error getting home dir")

  verifyExpectedCalls(t, []string{"IsWindows()"}, []string{}, []string{})
}

func TestDeactivateConfiguration_errorRemovingFile(t *testing.T) {
  setup()
  fsMock.ExpErrorOnRemoveFile = true

  err := repo.DeactivateConfiguration()
  test.ExpectHttpError(t, err, 500, "error removing file")

  verifyExpectedCalls(t, []string{"IsWindows()"}, []string{"GetHomeDir()"}, []string{})
}

func TestWriteDefaultCredentials(t *testing.T) {
  setup()
  envMock.ExpIsWindows = true

  defaultCredentialJson := "{\"credentials\":\"json\"}"
  err := repo.WriteDefaultCredentials(defaultCredentialJson)
  if err != nil {
    t.Fatalf("unexpected error %v", err)
  }

  expectedFile := fmt.Sprintf("WriteToFile(%v, %v)", expectedDefaultCredentialsFilePath, []byte(defaultCredentialJson))
  verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"}, []string{expectedFile}, []string{})
}

func TestWriteDefaultCredentials_errorGettingHomeDir(t *testing.T) {
  setup()
  fsMock.ExpErrorOnGetHomeDir = true

  err := repo.WriteDefaultCredentials("{\"credentials\":\"json\"}")
  test.ExpectHttpError(t, err, 500, "error getting home dir")

  verifyExpectedCalls(t, []string{"IsWindows()"}, []string{}, []string{})
}

func TestWriteDefaultCredentials_errorWritingFile(t *testing.T) {
  setup()
  fsMock.ExpErrorOnWriteToFile = true
  envMock.ExpIsWindows = true

  err := repo.WriteDefaultCredentials("{\"credentials\":\"json\"}")
  test.ExpectHttpError(t, err, 500, "error writing file")

  verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"}, []string{}, []string{})
}

func TestRemoveDefaultCredentials(t *testing.T) {
  setup()
  envMock.ExpIsWindows = true

  err := repo.RemoveDefaultCredentials()
  if err != nil {
    t.Fatalf("unexpected error %v", err)
  }

  expectedFile := fmt.Sprintf("RemoveFile(%v)", expectedDefaultCredentialsFilePath)
  verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"}, []string{expectedFile}, []string{})
}

func TestRemoveDefaultCredentials_errorGettingHomeDir(t *testing.T) {
  setup()
  fsMock.ExpErrorOnGetHomeDir = true

  err := repo.RemoveDefaultCredentials()
  test.ExpectHttpError(t, err, 500, "error getting home dir")

  verifyExpectedCalls(t, []string{"IsWindows()"}, []string{}, []string{})
}

func TestRemoveDefaultCredentials_errorRemovingFile(t *testing.T) {
  setup()
  fsMock.ExpErrorOnRemoveFile = true

  err := repo.RemoveDefaultCredentials()
  test.ExpectHttpError(t, err, 500, "error removing file")

  verifyExpectedCalls(t, []string{"IsWindows()"}, []string{"GetHomeDir()"}, []string{})
}

func TestWriteCredentialsToDb(t *testing.T) {
  setup()
  envMock.ExpIsWindows = true

  accountId := "account_id@domain.com"
  defaultCredentialJson := "{\"credentials\":\"json\"}"
  err := repo.WriteCredentialsToDb(accountId, defaultCredentialJson)
  if err != nil {
    t.Fatalf("unexpected error %v", err)
  }

  verifyExpectedCalls(t,
    []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
    []string{},
    []string{fmt.Sprintf("WriteCredentials(%v, %v, %v)", expectedCredentialsDbFilePath, accountId, defaultCredentialJson)})
}

func TestWriteCredentialsToDb_errorGettingHomeDir(t *testing.T) {
  setup()
  fsMock.ExpErrorOnGetHomeDir = true

  accountId := "account_id@domain.com"
  defaultCredentialJson := "{\"credentials\":\"json\"}"
  err := repo.WriteCredentialsToDb(accountId, defaultCredentialJson)
  test.ExpectHttpError(t, err, 500, "error getting home dir")

  verifyExpectedCalls(t, []string{"IsWindows()"}, []string{}, []string{})
}

func TestWriteCredentialsToDb_errorWritingCredentials(t *testing.T) {
  setup()
  envMock.ExpIsWindows = true
  dbTableMock.ExpErrorOnExecInsertQuery = true

  accountId := "account_id@domain.com"
  defaultCredentialJson := "{\"credentials\":\"json\"}"
  err := repo.WriteCredentialsToDb(accountId, defaultCredentialJson)
  test.ExpectHttpError(t, err, 500, "error executing insert query")

  verifyExpectedCalls(t,
    []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
    []string{},
    []string{})
}

func TestRemoveCredentialsFromDb(t *testing.T) {
  setup()
  envMock.ExpIsWindows = true
  fsMock.ExpDoesFileExist = true

  accountId := "account_id@domain.com"
  err := repo.RemoveCredentialsFromDb(accountId)
  if err != nil {
    t.Fatalf("unexpected error %v", err)
  }

  verifyExpectedCalls(t,
    []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
    []string{fmt.Sprintf("DoesFileExist(%v)", expectedCredentialsDbFilePath)},
    []string{fmt.Sprintf("RemoveCredentials(%v, %v)", expectedCredentialsDbFilePath, accountId)})
}

func TestRemoveCredentialsFromDb_DbFileDoesNotExist(t *testing.T) {
  setup()
  envMock.ExpIsWindows = true

  accountId := "account_id@domain.com"
  err := repo.RemoveCredentialsFromDb(accountId)
  if err != nil {
    t.Fatalf("unexpected error %v", err)
  }

  verifyExpectedCalls(t,
    []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
    []string{fmt.Sprintf("DoesFileExist(%v)", expectedCredentialsDbFilePath)},
    []string{})
}

func TestRemoveCredentialsFromDb_errorGettingHomeDir(t *testing.T) {
  setup()
  fsMock.ExpErrorOnGetHomeDir = true

  accountId := "account_id@domain.com"
  err := repo.RemoveCredentialsFromDb(accountId)
  test.ExpectHttpError(t, err, 500, "error getting home dir")

  verifyExpectedCalls(t, []string{"IsWindows()"}, []string{}, []string{})
}

func TestRemoveCredentialsFromDb_errorRemovingCredentialsFromDb(t *testing.T) {
  setup()
  envMock.ExpIsWindows = true
  fsMock.ExpDoesFileExist = true
  dbTableMock.ExpErrorOnExecDeleteQuery = true

  accountId := "account_id@domain.com"
  err := repo.RemoveCredentialsFromDb(accountId)
  test.ExpectHttpError(t, err, 500, "error executing delete query")

  verifyExpectedCalls(t,
    []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
    []string{fmt.Sprintf("DoesFileExist(%v)", expectedCredentialsDbFilePath)},
    []string{})
}
