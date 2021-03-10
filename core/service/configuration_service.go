package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"leapp_daemon/core/model"
	"leapp_daemon/shared/const"
	"os"
	"sync"
)

// The zero value is an unlocked mutex
var configurationFileMutex sync.Mutex

func getInitialConfiguration() model.Configuration {
	return model.Configuration{
		SsoUrl: "",
		ProxyConfiguration: model.ProxyConfiguration{
			ProxyProtocol: "https",
			ProxyUrl: "",
			ProxyPort: 8080,
			Username: "",
			Password: "",
		},
		FederatedAwsSessions: make([]model.FederatedAwsSession, 0),
		PlainAwsSessions:     make([]model.PlainAwsSession, 0),
	}
}

func CreateConfiguration() error {
	configuration := getInitialConfiguration()
	err := UpdateConfiguration(&configuration, true)
	if err != nil { return err }
	return nil
}

func ReadConfiguration() (*model.Configuration, error) {
	homeDir, err := GetHomeDir()
	if err != nil { return nil, err }

	configurationFilePath := fmt.Sprintf("%s/%s", homeDir, _const.ConfigurationFilePath)

	encryptedText, err := ioutil.ReadFile(configurationFilePath)
	if err != nil { return nil, err }

	plainText, err := Decrypt(string(encryptedText))
	if err != nil { return nil, err }

	return unmarshalConfiguration(plainText), nil
}

func UpdateConfiguration(configuration *model.Configuration, deleteExistingFile bool) error {
	configurationFileMutex.Lock()
	defer configurationFileMutex.Unlock()

	homeDir, err := GetHomeDir()
	if err != nil { return err }

	configurationFilePath := fmt.Sprintf("%s/%s", homeDir, _const.ConfigurationFilePath)

	if deleteExistingFile == true {
		if DoesFileExist(configurationFilePath) {
			err = os.Remove(configurationFilePath)
			if err != nil {
				return err
			}
		}
	}

	configurationJson, err := json.Marshal(configuration)
	if err != nil { return err }

	encryptedConfigurationJson, err := Encrypt(string(configurationJson))
	if err != nil { return err }

	err = ioutil.WriteFile(configurationFilePath, []byte(encryptedConfigurationJson), 0644)
	if err != nil { return err }

	return nil
}

func unmarshalConfiguration(configurationJson string) *model.Configuration {
	var tmp model.Configuration
	_ = json.Unmarshal([]byte(configurationJson), &tmp)
	return &tmp
}