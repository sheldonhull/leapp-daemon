package repository

import (
	"fmt"
	"leapp_daemon/infrastructure/http/http_error"
	"path/filepath"
)

type GcloudConfigFileSystem interface {
	DoesFileExist(path string) bool
	GetHomeDir() (string, error)
	WriteToFile(path string, data []byte) error
	RemoveFile(path string) error
}

type GcloudEnvironment interface {
	GetEnvironmentVariable(variable string) string
	IsCommandAvailable(command string) bool
	IsWindows() bool
}

type GcloudConfigurationRepository struct {
	FileSystem  GcloudConfigFileSystem
	Environment GcloudEnvironment
}

func (repo *GcloudConfigurationRepository) gcloudConfigDir() (string, error) {
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

func (repo *GcloudConfigurationRepository) getGcloudConfigFilePath(configurationName string) (string, error) {
	configDir, err := repo.gcloudConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "configurations", fmt.Sprintf("config_leapp_%v", configurationName)), nil
}

func (repo *GcloudConfigurationRepository) IsGcloudCliAvailable() bool {
	return repo.Environment.IsCommandAvailable("gcloud")
}

func (repo *GcloudConfigurationRepository) DoesGcloudConfigFolderExist() (bool, error) {
	configDir, err := repo.gcloudConfigDir()
	if err != nil {
		return false, err
	}
	return repo.FileSystem.DoesFileExist(configDir), nil
}

func (repo *GcloudConfigurationRepository) CreateConfiguration(configurationName string, account string, project string) error {
	configFilePath, err := repo.getGcloudConfigFilePath(configurationName)
	if err != nil {
		return err
	}
	configFileContent := fmt.Sprintf("[core]\naccount = %v\nproject = %v\n", account, project)
	err = repo.FileSystem.WriteToFile(configFilePath, []byte(configFileContent))
	if err != nil {
		return http_error.NewInternalServerError(err)
	}
	return nil
}

func (repo *GcloudConfigurationRepository) RemoveConfiguration(configurationName string) error {
	configFilePath, err := repo.getGcloudConfigFilePath(configurationName)
	if err != nil {
		return err
	}
	err = repo.FileSystem.RemoveFile(configFilePath)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}
	return nil
}
