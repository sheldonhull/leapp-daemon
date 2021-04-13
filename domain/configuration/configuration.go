package domain

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  session2 "leapp_daemon/domain/session"
  encryption2 "leapp_daemon/infrastructure/encryption"
  file_system2 "leapp_daemon/infrastructure/file_system"
  http_error2 "leapp_daemon/infrastructure/http/http_error"
  "os"
  "sync"
)

const configurationFilePath = `.Leapp/Leapp-lock.json`

type Configuration struct {
	ProxyConfiguration   *ProxyConfiguration
	PlainAwsSessions     []*session2.PlainAwsSession
	FederatedAwsSessions []*session2.FederatedAwsSession
	TrustedAwsSessions   []*session2.TrustedAwsSession
	NamedProfiles        []*session2.NamedProfile
}

type ProxyConfiguration struct {
	ProxyProtocol string
	ProxyUrl string
	ProxyPort uint64
	Username string
	Password string
}

// The zero value is an unlocked mutex
var configurationFileMutex sync.Mutex

func CreateConfiguration() error {
	configuration := getInitialConfiguration()
	err := UpdateConfiguration(configuration, true)
	if err != nil { return err }
	return nil
}

func ReadConfiguration() (*Configuration, error) {
	homeDir, err := file_system2.GetHomeDir()
	if err != nil { return nil, err }

	configurationFilePath := fmt.Sprintf("%s/%s", homeDir, configurationFilePath)

	encryptedText, err := ioutil.ReadFile(configurationFilePath)
	if err != nil {
		return nil, http_error2.NewInternalServerError(err)
	}

	plainText, err := encryption2.Decrypt(string(encryptedText))
	if err != nil { return nil, err }

	return unmarshalConfiguration(plainText), nil
}

func UpdateConfiguration(configuration *Configuration, deleteExistingFile bool) error {
	configurationFileMutex.Lock()
	defer configurationFileMutex.Unlock()

	homeDir, err := file_system2.GetHomeDir()
	if err != nil { return err }

	configurationFilePath := fmt.Sprintf("%s/%s", homeDir, configurationFilePath)

	if deleteExistingFile == true {
		if file_system2.DoesFileExist(configurationFilePath) {
			err = os.Remove(configurationFilePath)
			if err != nil {
				return err
			}
		}
	}

	configurationJson, err := json.Marshal(configuration)
	if err != nil { return err }

	encryptedConfigurationJson, err := encryption2.Encrypt(string(configurationJson))
	if err != nil { return err }

	err = ioutil.WriteFile(configurationFilePath, []byte(encryptedConfigurationJson), 0644)
	if err != nil { return err }

	return nil
}

func (config *Configuration) Update() error {
  err := UpdateConfiguration(config, false)
  if err != nil {
    return err
  }
  return nil
}

func (config *Configuration) GetPlainAwsSessions() ([]*session2.PlainAwsSession, error) {
	return config.PlainAwsSessions, nil
}

func (config *Configuration) SetPlainAwsSessions(plainAwsSessions []*session2.PlainAwsSession) error {
	config.PlainAwsSessions = plainAwsSessions
	return nil
}

func (config *Configuration) GetFederatedAwsSessions() ([]*session2.FederatedAwsSession, error) {
	return config.FederatedAwsSessions, nil
}

func (config *Configuration) SetFederatedAwsSessions(federatedAwsSessions []*session2.FederatedAwsSession) error {
	config.FederatedAwsSessions = federatedAwsSessions
	return nil
}

func (config *Configuration) GetTrustedAwsSessions() ([]*session2.TrustedAwsSession, error) {
  return config.TrustedAwsSessions, nil
}

func (config *Configuration) SetTrustedAwsSessions(trustedAwsSessions []*session2.TrustedAwsSession) error {
  config.TrustedAwsSessions = trustedAwsSessions
  return nil
}

func (config *Configuration) GetNamedProfiles() ([]*session2.NamedProfile, error) {
	return config.NamedProfiles, nil
}

func (config *Configuration) SetNamedProfiles(namedProfiles []*session2.NamedProfile) error {
	config.NamedProfiles = namedProfiles
	return nil
}

func (config *Configuration) GetAllSessions() []session2.Rotatable {
  sessions := make([]session2.Rotatable, 0)

  for i := range config.PlainAwsSessions {
    sess := config.PlainAwsSessions[i]
    sessions = append(sessions, sess)
  }

  for i := range config.FederatedAwsSessions {
    sess := config.FederatedAwsSessions[i]
    sessions = append(sessions, sess)
  }

  return sessions
}

func getInitialConfiguration() *Configuration {
	return &Configuration{
		ProxyConfiguration: &ProxyConfiguration{
			ProxyProtocol: "https",
			ProxyUrl: "",
			ProxyPort: 8080,
			Username: "",
			Password: "",
		},
		FederatedAwsSessions: make([]*session2.FederatedAwsSession, 0),
		PlainAwsSessions:     make([]*session2.PlainAwsSession, 0),
	}
}

func unmarshalConfiguration(configurationJson string) *Configuration {
	var tmp Configuration
	_ = json.Unmarshal([]byte(configurationJson), &tmp)
	return &tmp
}
