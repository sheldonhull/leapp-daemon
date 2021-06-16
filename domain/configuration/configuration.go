package configuration

import (
	"encoding/json"
	"fmt"
	"leapp_daemon/domain/named_profile"
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/http/http_error"
)

type Configuration struct {
	ProxyConfiguration   ProxyConfiguration
	PlainAwsSessions     []session.AwsPlainSession
	FederatedAwsSessions []session.FederatedAwsSession
	TrustedAwsSessions   []session.TrustedAwsSession
	GcpPlainSessions     []session.GcpPlainSession
	NamedProfiles        []named_profile.NamedProfile
}

type ProxyConfiguration struct {
	ProxyProtocol string
	ProxyUrl      string
	ProxyPort     uint64
	Username      string
	Password      string
}

func GetDefaultConfiguration() Configuration {
	return Configuration{
		ProxyConfiguration: ProxyConfiguration{
			ProxyProtocol: "https",
			ProxyUrl:      "",
			ProxyPort:     8080,
			Username:      "",
			Password:      "",
		},
		FederatedAwsSessions: make([]session.FederatedAwsSession, 0),
		PlainAwsSessions:     make([]session.AwsPlainSession, 0),
		GcpPlainSessions:     make([]session.GcpPlainSession, 0),
	}
}

func FromJson(configurationJson string) Configuration {
	var config Configuration
	_ = json.Unmarshal([]byte(configurationJson), &config)
	return config
}

func (config *Configuration) AddAwsPlainSession(plainAwsSession session.AwsPlainSession) error {
	sessions, err := config.GetAllAwsPlainSessions()
	if err != nil {
		return err
	}

	for _, sess := range sessions {
		if plainAwsSession.Id == sess.Id {
			return http_error.NewConflictError(fmt.Errorf("a AwsPlainSession with id " + plainAwsSession.Id +
				" is already present"))
		}
	}

	sessions = append(sessions, plainAwsSession)
	config.PlainAwsSessions = sessions

	return nil
}

func (config *Configuration) GetAllAwsPlainSessions() ([]session.AwsPlainSession, error) {
	return config.PlainAwsSessions, nil
}

func (config *Configuration) RemoveAwsPlainSession(plainAwsSession session.AwsPlainSession) error {
	sessions, err := config.GetAllAwsPlainSessions()
	if err != nil {
		return err
	}

	for i, sess := range sessions {
		if plainAwsSession.Id == sess.Id {
			config.PlainAwsSessions = append(config.PlainAwsSessions[:i], config.PlainAwsSessions[i+1:]...)
			return nil
		}
	}

	return http_error.NewNotFoundError(fmt.Errorf("AwsPlainSession with id " + plainAwsSession.Id +
		" not found"))
}

func (config *Configuration) AddNamedProfile(namedProfile named_profile.NamedProfile) error {
	for _, tmpNamedProfile := range config.NamedProfiles {
		if namedProfile.Name == tmpNamedProfile.Name {
			return http_error.NewConflictError(fmt.Errorf("a NamedProfile with name " + namedProfile.Name +
				" is already present"))
		}
	}
	config.NamedProfiles = append(config.NamedProfiles, namedProfile)
	return nil
}

func (config *Configuration) FindNamedProfileByName(name string) (named_profile.NamedProfile, error) {
	var namedProfile named_profile.NamedProfile

	for _, tmpNamedProfile := range config.NamedProfiles {
		if name == tmpNamedProfile.Name {
			return tmpNamedProfile, nil
		}
	}

	return namedProfile, http_error.NewNotFoundError(fmt.Errorf("NamedProfile with name " + name + " not found"))
}

func (config *Configuration) DoesNamedProfileExist(name string) bool {
	for _, tmpNamedProfile := range config.NamedProfiles {
		if name == tmpNamedProfile.Name {
			return true
		}
	}
	return false
}

/*
// The zero value is an unlocked mutex
var configurationFileMutex sync.Mutex

func CreateConfiguration() error {
	configuration := GetDefaultConfiguration()
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

	return FromJson(plainText), nil
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

	err = ioutil.WriteToFile(configurationFilePath, []byte(encryptedConfigurationJson), 0644)
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

func (config *Configuration) GetNamedProfiles() ([]*named_profile.NamedProfile, error) {
	return config.NamedProfiles, nil
}

func (config *Configuration) SetNamedProfiles(namedProfiles []*named_profile.NamedProfile) error {
	config.NamedProfiles = namedProfiles
	return nil
}

func (config *Configuration) GetAllSessions() []session.Rotatable {
  sessions := make([]session.Rotatable, 0)

  for i := range config.plainAwsSessions {
    sess := config.plainAwsSessions[i]
    sessions = append(sessions, sess)
  }

  for i := range config.FederatedAwsSessions {
    sess := config.FederatedAwsSessions[i]
    sessions = append(sessions, sess)
  }

  return sessions
}
*/
