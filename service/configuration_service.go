package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"leapp_daemon/service/domain"
	"os"
)

func getInitialConfiguration() domain.Configuration {
	return domain.Configuration{
		SsoUrl: "",
		ProxyConfiguration: domain.ProxyConfiguration{
			ProxyProtocol: "https",
			ProxyUrl: "",
			ProxyPort: 8080,
			Username: "",
			Password: "",
		},
		FederatedAwsAccountSessions: make([]domain.FederatedAwsAccountSession, 0),
		PlainAwsAccountSessions: make([]domain.PlainAwsAccountSession, 0),
	}
}

func CreateConfiguration() error {
	configuration := getInitialConfiguration()
	err := UpdateConfiguration(&configuration, true)
	if err != nil { return err }
	return nil
}

func ReadConfiguration() (*domain.Configuration, error) {
	homeDir, err := GetHomeDir()
	if err != nil { return nil, err }

	configurationFilePath := fmt.Sprintf("%s/%s", homeDir, domain.ConfigurationFilePath)

	encryptedText, err := ioutil.ReadFile(configurationFilePath)
	if err != nil { return nil, err }

	plainText, err := Decrypt(string(encryptedText))
	if err != nil { return nil, err }

	return unmarshalConfiguration(plainText), nil
}

func UpdateConfiguration(configuration *domain.Configuration, deleteExistingFile bool) error {
	homeDir, err := GetHomeDir()
	if err != nil { return err }

	configurationFilePath := fmt.Sprintf("%s/%s", homeDir, domain.ConfigurationFilePath)

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

func unmarshalConfiguration(configurationJson string) *domain.Configuration {
	var tmp domain.Configuration
	_ = json.Unmarshal([]byte(configurationJson), &tmp)
	return &tmp
}