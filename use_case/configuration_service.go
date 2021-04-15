package use_case

import (
  "leapp_daemon/domain/configuration"
  "leapp_daemon/infrastructure/http/http_error"
)

type ConfigurationService struct {
  ConfigurationRepository configuration.Repository
}

func(service *ConfigurationService) Create() error {
  config := configuration.GetInitialConfiguration()
  err := service.ConfigurationRepository.CreateConfiguration(config)
  if err != nil {
    return http_error.NewInternalServerError(err)
  }
  return nil
}

func(service *ConfigurationService) Get() (configuration.Configuration, error) {
  var config configuration.Configuration
  config, err := service.ConfigurationRepository.GetConfiguration()
  if err != nil {
    return config, err
  }
  return config, nil
}

func(service *ConfigurationService) Update(configuration configuration.Configuration) error {
  err := service.ConfigurationRepository.UpdateConfiguration(configuration)
  if err != nil {
    return err
  }
  return nil
}
