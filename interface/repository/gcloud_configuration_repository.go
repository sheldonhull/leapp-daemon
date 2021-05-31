package repository

import (
  "leapp_daemon/infrastructure/http/http_error"
  "path/filepath"
  "sync"
)

var gcloudFileMutex sync.Mutex

type GcloudConfigFileSystem interface {
  DoesFileExist(path string) bool
  GetHomeDir() (string, error)
}

type GcloudEnvironment interface {
  GetEnvironmentVariable(variable string) string
  IsWindows() bool
}

type GcloudConfigurationRepository struct {
  FileSystem  GcloudConfigFileSystem
  Environment GcloudEnvironment
}

func (repo *GcloudConfigurationRepository) configDir() (string, error) {
  gcloudDirectory := "gcloud"
  if repo.Environment.IsWindows() {
    return filepath.Join(repo.Environment.GetEnvironmentVariable("APPDATA"), gcloudDirectory), nil
  }
  dir, err := repo.FileSystem.GetHomeDir()
  if err != nil {
    return "", http_error.NewInternalServerError(err)
  }
  return filepath.Join(dir, ".config", gcloudDirectory), nil
}
