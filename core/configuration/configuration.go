package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"leapp_daemon/constant"
	"leapp_daemon/core/encryption"
	"leapp_daemon/core/file_system"
	"os"
	"sync"
)

type Configuration struct {
	SsoUrl               string
	ProxyConfiguration   *ProxyConfiguration
	FederatedAwsSessions []*FederatedAwsSession
	PlainAwsSessions     []*PlainAwsSession
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

	configurationFilePath := fmt.Sprintf("%s/%s", homeDir, constant.ConfigurationFilePath)

	encryptedText, err := ioutil.ReadFile(configurationFilePath)
	if err != nil { return nil, err }

	plainText, err := encryption.Decrypt(string(encryptedText))
	if err != nil { return nil, err }

	return unmarshalConfiguration(plainText), nil
}

func UpdateConfiguration(configuration *Configuration, deleteExistingFile bool) error {
	configurationFileMutex.Lock()
	defer configurationFileMutex.Unlock()

	homeDir, err := file_system.GetHomeDir()
	if err != nil { return err }

	configurationFilePath := fmt.Sprintf("%s/%s", homeDir, constant.ConfigurationFilePath)

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
		FederatedAwsSessions: make([]*FederatedAwsSession, 0),
		PlainAwsSessions:     make([]*PlainAwsSession, 0),
	}
}

func unmarshalConfiguration(configurationJson string) *Configuration {
	var tmp Configuration
	_ = json.Unmarshal([]byte(configurationJson), &tmp)
	return &tmp
}
