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

type GcloudCredentialsTable interface {
	WriteCredentials(sqlFilePath string, accountId string, value string) error
	RemoveCredentials(sqlFilePath string, accountId string) error
}

type GcloudConfigurationRepository struct {
	FileSystem       GcloudConfigFileSystem
	Environment      GcloudEnvironment
	CredentialsTable GcloudCredentialsTable
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

func (repo *GcloudConfigurationRepository) getGcloudActiveConfigFilePath() (string, error) {
	configDir, err := repo.gcloudConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "active_config"), nil
}

func (repo *GcloudConfigurationRepository) getGcloudDefaultCredentialsFilePath() (string, error) {
	configDir, err := repo.gcloudConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "application_default_credentials.json"), nil
}

func (repo *GcloudConfigurationRepository) getGcloudCredentialsDbFilePath() (string, error) {
	configDir, err := repo.gcloudConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "credentials.db"), nil
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

func (repo *GcloudConfigurationRepository) ActivateConfiguration(configurationName string) error {
	activeConfigFilePath, err := repo.getGcloudActiveConfigFilePath()
	if err != nil {
		return err
	}

	err = repo.FileSystem.WriteToFile(activeConfigFilePath, []byte(configurationName))
	if err != nil {
		return http_error.NewInternalServerError(err)
	}
	return nil
}

func (repo *GcloudConfigurationRepository) DeactivateConfiguration() error {
	activeConfigFilePath, err := repo.getGcloudActiveConfigFilePath()
	if err != nil {
		return err
	}

	err = repo.FileSystem.RemoveFile(activeConfigFilePath)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (repo *GcloudConfigurationRepository) WriteDefaultCredentials(credentialsJson string) error {
	defaultCredentialsFilePath, err := repo.getGcloudDefaultCredentialsFilePath()
	if err != nil {
		return err
	}

	err = repo.FileSystem.WriteToFile(defaultCredentialsFilePath, []byte(credentialsJson))
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (repo *GcloudConfigurationRepository) RemoveDefaultCredentials() error {
	defaultCredentialFilePath, err := repo.getGcloudDefaultCredentialsFilePath()
	if err != nil {
		return err
	}

	err = repo.FileSystem.RemoveFile(defaultCredentialFilePath)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (repo *GcloudConfigurationRepository) WriteCredentialsToDb(accountId string, credentialsJson string) error {
	credentialsDbFilePath, err := repo.getGcloudCredentialsDbFilePath()
	if err != nil {
		return err
	}

	err = repo.CredentialsTable.WriteCredentials(credentialsDbFilePath, accountId, credentialsJson)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (repo *GcloudConfigurationRepository) RemoveCredentialsFromDb(accountId string) error {
	credentialsDbFilePath, err := repo.getGcloudCredentialsDbFilePath()
	if err != nil {
		return err
	}

	if !repo.FileSystem.DoesFileExist(credentialsDbFilePath) {
		return nil
	}

	err = repo.CredentialsTable.RemoveCredentials(credentialsDbFilePath, accountId)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}
