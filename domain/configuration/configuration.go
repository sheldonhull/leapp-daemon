package configuration

import (
  "encoding/json"
  "leapp_daemon/domain/session"
)

//const configurationFilePath = `.Leapp/Leapp-lock.json`

type Configuration struct {
	ProxyConfiguration   *ProxyConfiguration
	PlainAwsSessions     []*session.PlainAwsSession
	FederatedAwsSessions []*session.FederatedAwsSession
	TrustedAwsSessions   []*session.TrustedAwsSession
	NamedProfiles        []*session.NamedProfile
}

type ProxyConfiguration struct {
	ProxyProtocol string
	ProxyUrl string
	ProxyPort uint64
	Username string
	Password string
}

func GetInitialConfiguration() *Configuration {
  return &Configuration{
    ProxyConfiguration: &ProxyConfiguration{
      ProxyProtocol: "https",
      ProxyUrl: "",
      ProxyPort: 8080,
      Username: "",
      Password: "",
    },
    FederatedAwsSessions: make([]*session.FederatedAwsSession, 0),
    PlainAwsSessions:     make([]*session.PlainAwsSession, 0),
  }
}

func UnmarshalConfiguration(configurationJson string) *Configuration {
  var tmp Configuration
  _ = json.Unmarshal([]byte(configurationJson), &tmp)
  return &tmp
}

/*
// The zero value is an unlocked mutex
var configurationFileMutex sync.Mutex

func CreateConfiguration() error {
	configuration := GetInitialConfiguration()
	err := UpdateConfiguration(configuration, true)
	if err != nil { return err }
	return nil
}

func ReadConfiguration() (*Configuration, error) {
	homeDir, err := file_system.GetHomeDir()
	if err != nil { return nil, err }

	configurationFilePath := fmt.Sprintf("%s/%s", homeDir, configurationFilePath)

	encryptedText, err := ioutil.ReadFile(configurationFilePath)
	if err != nil {
		return nil, http_error.NewInternalServerError(err)
	}

	plainText, err := encryption.Decrypt(string(encryptedText))
	if err != nil { return nil, err }

	return UnmarshalConfiguration(plainText), nil
}

func UpdateConfiguration(configuration *Configuration, deleteExistingFile bool) error {
	configurationFileMutex.Lock()
	defer configurationFileMutex.Unlock()

	homeDir, err := file_system.GetHomeDir()
	if err != nil { return err }

	configurationFilePath := fmt.Sprintf("%s/%s", homeDir, configurationFilePath)

	if deleteExistingFile == true {
		if file_system.DoesFileExist(configurationFilePath) {
			err = os.Remove(configurationFilePath)
			if err != nil {
				return err
			}
		}
	}

	configurationJson, err := json.Marshal(configuration)
	if err != nil { return err }

	encryptedConfigurationJson, err := encryption.Encrypt(string(configurationJson))
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
*/

func (config *Configuration) GetPlainAwsSessions() ([]*session.PlainAwsSession, error) {
	return config.PlainAwsSessions, nil
}

func (config *Configuration) SetPlainAwsSessions(plainAwsSessions []*session.PlainAwsSession) error {
	config.PlainAwsSessions = plainAwsSessions
	return nil
}

func (config *Configuration) GetFederatedAwsSessions() ([]*session.FederatedAwsSession, error) {
	return config.FederatedAwsSessions, nil
}

func (config *Configuration) SetFederatedAwsSessions(federatedAwsSessions []*session.FederatedAwsSession) error {
	config.FederatedAwsSessions = federatedAwsSessions
	return nil
}

func (config *Configuration) GetTrustedAwsSessions() ([]*session.TrustedAwsSession, error) {
  return config.TrustedAwsSessions, nil
}

func (config *Configuration) SetTrustedAwsSessions(trustedAwsSessions []*session.TrustedAwsSession) error {
  config.TrustedAwsSessions = trustedAwsSessions
  return nil
}

func (config *Configuration) GetNamedProfiles() ([]*session.NamedProfile, error) {
	return config.NamedProfiles, nil
}

func (config *Configuration) SetNamedProfiles(namedProfiles []*session.NamedProfile) error {
	config.NamedProfiles = namedProfiles
	return nil
}

func (config *Configuration) GetAllSessions() []session.Rotatable {
  sessions := make([]session.Rotatable, 0)

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
