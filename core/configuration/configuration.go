package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"leapp_daemon/core/encryption"
	"leapp_daemon/core/file_system"
	"leapp_daemon/core/session"
	"leapp_daemon/custom_error"
	"os"
	"sync"
)

const configurationFilePath = `.Leapp/Leapp-lock.json`

type Configuration struct {
	SsoUrl               string
	ProxyConfiguration   *ProxyConfiguration
	PlainAwsSessions     []*session.PlainAwsSession
	FederatedAwsSessions []*session.FederatedAwsSession
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
	homeDir, err := file_system.GetHomeDir()
	if err != nil { return nil, err }

	configurationFilePath := fmt.Sprintf("%s/%s", homeDir, configurationFilePath)

	encryptedText, err := ioutil.ReadFile(configurationFilePath)
	if err != nil { return nil, custom_error.NewNotFoundError(err) }

	plainText, err := encryption.Decrypt(string(encryptedText))
	if err != nil { return nil, custom_error.NewBadRequestError(err) }

	return unmarshalConfiguration(plainText), nil
}

func UpdateConfiguration(configuration *Configuration, deleteExistingFile bool) error {
	configurationFileMutex.Lock()
	defer configurationFileMutex.Unlock()

	homeDir, err := file_system.GetHomeDir()
	if err != nil { return custom_error.NewNotFoundError(err) }

	configurationFilePath := fmt.Sprintf("%s/%s", homeDir, configurationFilePath)

	if deleteExistingFile == true {
		if file_system.DoesFileExist(configurationFilePath) {
			err = os.Remove(configurationFilePath)
			if err != nil {
				return custom_error.NewBadRequestError(err)
			}
		}
	}

	configurationJson, err := json.Marshal(configuration)
	if err != nil { return custom_error.NewBadRequestError(err) }

	encryptedConfigurationJson, err := encryption.Encrypt(string(configurationJson))
	if err != nil { return custom_error.NewBadRequestError(err) }

	err = ioutil.WriteFile(configurationFilePath, []byte(encryptedConfigurationJson), 0644)
	if err != nil { return custom_error.NewBadRequestError(err) }

	return nil
}

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

func (config *Configuration) Update() error {
	err := UpdateConfiguration(config, false)
	if err != nil {
		return err
	}
	return nil
}

func getInitialConfiguration() *Configuration {
	return &Configuration{
		SsoUrl: "",
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

func unmarshalConfiguration(configurationJson string) *Configuration {
	var tmp Configuration
	_ = json.Unmarshal([]byte(configurationJson), &tmp)
	return &tmp
}
