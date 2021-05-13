package configuration

import (
  "leapp_daemon/domain/session"
  "sync"
)

var facadeSingleton *facade
var configurationLock sync.Mutex
var plainAwsSessionsLock sync.Mutex

type Observer interface {
  Update(configuration Configuration)
}

type facade struct {
  // ConfigurationRepository Repository
  Configuration Configuration
  Observers     []Observer
}

func GetFacade() *facade {
  if facadeSingleton != nil {
    facadeSingleton = &facade {
      Configuration: Configuration{
        ProxyConfiguration: ProxyConfiguration{
          ProxyProtocol: "https",
          ProxyUrl: "",
          ProxyPort: 8080,
          Username: "",
          Password: "",
        },
        FederatedAwsSessions: make([]session.FederatedAwsSession, 0),
        PlainAwsSessions:     make([]session.PlainAwsSession, 0),
      },
    }
  }
  return facadeSingleton
}

func(fac *facade) Subscribe(observer Observer) {
  fac.Observers = append(fac.Observers, observer)
}

func(fac *facade) GetConfiguration() Configuration {
  return fac.Configuration
}

func(fac *facade) SetConfiguration(configuration Configuration) {
  configurationLock.Lock()
  defer configurationLock.Unlock()
  fac.Configuration = configuration
}

/*func(service *facade) Create() error {
  config := GetInitialConfiguration()
  err := service.ConfigurationRepository.CreateConfiguration(config)
  if err != nil {
    return http_error.NewInternalServerError(err)
  }
  return nil
}

func(service *facade) Get() (Configuration, error) {
  var config Configuration
  config, err := service.ConfigurationRepository.GetConfiguration()
  if err != nil {
    return config, err
  }
  return config, nil
}

func(service *facade) Update(configuration Configuration) error {
  err := service.ConfigurationRepository.UpdateConfiguration(configuration)
  if err != nil {
    return err
  }
  return nil
}
*/
