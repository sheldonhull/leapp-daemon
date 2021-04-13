package use_case

import (
  "encoding/json"
  "fmt"
  "leapp_daemon/domain/configuration"
  "leapp_daemon/infrastructure/http/http_error"
  "sync"
)

const configurationFilePath = `.Leapp/Leapp-lock.json`

// The zero value is an unlocked mutex
var configurationFileMutex sync.Mutex

type FileSystem interface {
  DoesFileExist(path string) bool
  GetHomeDir() (string, error)
  ReadFile(path string) ([]byte, error)
  RemoveFile(path string) error
  WriteFile(path string, data []byte) error
}

type Encryption interface {
  Decrypt(encryptedText string) (string, error)
  Encrypt(plainText string) (string, error)
}

type ConfigurationService struct {
  FileSystem FileSystem
  Encryption Encryption
}

func(service *ConfigurationService) Create() error {
  config := configuration.GetInitialConfiguration()
  err := service.Update(config, true)
  if err != nil {
    return http_error.NewInternalServerError(err)
  }
  return nil
}

func(service *ConfigurationService) Read() (*configuration.Configuration, error) {
  homeDir, err := service.FileSystem.GetHomeDir()
  if err != nil {
    return nil, http_error.NewInternalServerError(err)
  }

  configurationFilePath := fmt.Sprintf("%s/%s", homeDir, configurationFilePath)

  encryptedText, err := service.FileSystem.ReadFile(configurationFilePath)
  if err != nil {
    return nil, http_error.NewInternalServerError(err)
  }

  plainText, err := service.Encryption.Decrypt(string(encryptedText))
  if err != nil {
    return nil, http_error.NewInternalServerError(err)
  }

  return configuration.UnmarshalConfiguration(plainText), nil
}

func(service *ConfigurationService) Update(configuration *configuration.Configuration, deleteExistingFile bool) error {
  configurationFileMutex.Lock()
  defer configurationFileMutex.Unlock()

  homeDir, err := service.FileSystem.GetHomeDir()
  if err != nil {
    return http_error.NewInternalServerError(err)
  }

  configurationFilePath := fmt.Sprintf("%s/%s", homeDir, configurationFilePath)

  if deleteExistingFile == true {
    if service.FileSystem.DoesFileExist(configurationFilePath) {
      err = service.FileSystem.RemoveFile(configurationFilePath)
      if err != nil {
        return http_error.NewInternalServerError(err)
      }
    }
  }

  configurationJson, err := json.Marshal(configuration)
  if err != nil {
    return http_error.NewInternalServerError(err)
  }

  encryptedConfigurationJson, err := service.Encryption.Encrypt(string(configurationJson))
  if err != nil {
    return http_error.NewInternalServerError(err)
  }

  err = service.FileSystem.WriteFile(configurationFilePath, []byte(encryptedConfigurationJson))
  if err != nil {
    return http_error.NewInternalServerError(err)
  }

  return nil
}
