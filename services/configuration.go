package services

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"leapp_daemon/services/domain"
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
		Sessions: make([]domain.Session, 0),
	}
}

func CreateConfiguration() error {
	configuration := getInitialConfiguration()
	err := UpdateConfiguration(&configuration, false)
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
	var configuration domain.Configuration
	configuration.SsoUrl = gjson.Get(configurationJson, "SsoUrl").String()

	var proxyConfiguration domain.ProxyConfiguration
	proxyConfiguration.ProxyProtocol = gjson.Get(configurationJson, "ProxyConfiguration.ProxyProtocol").String()
	proxyConfiguration.ProxyUrl = gjson.Get(configurationJson, "ProxyConfiguration.ProxyUrl").String()
	proxyConfiguration.ProxyPort = gjson.Get(configurationJson, "ProxyConfiguration.ProxyPort").Uint()
	proxyConfiguration.Username = gjson.Get(configurationJson, "ProxyConfiguration.Username").String()
	proxyConfiguration.Password = gjson.Get(configurationJson, "ProxyConfiguration.Password").String()

	configuration.ProxyConfiguration = proxyConfiguration

	var sessions []domain.Session
	sessionsJsonArray := gjson.Get(configurationJson, "Sessions").Array()

	if len(sessionsJsonArray) > 0 {
		for _, sessionJsonValue := range sessionsJsonArray {
			sessionJson := sessionJsonValue.String()
			var session domain.Session

			session.Id = gjson.Get(sessionJson, "Id").String()
			session.Active = gjson.Get(sessionJson, "Active").Bool()
			session.Loading = gjson.Get(sessionJson, "Loading").Bool()
			session.LastStopDate = gjson.Get(sessionJson, "LastStopDate").String()

			accountJsonValue := gjson.Get(sessionJson, "Account").Value()

			if accountJsonValue != nil {
				accountJson := gjson.Get(sessionJson, "Account").String()
				var account domain.Account

				account.Id = gjson.Get(accountJson, "Id").String()
				account.Name = gjson.Get(sessionJson, "Name").String()
				account.AccountNumber = gjson.Get(sessionJson, "AccountNumber").String()

				roleJsonValue := gjson.Get(sessionJson, "Role").Value()

				if roleJsonValue != nil {
					roleJson := gjson.Get(sessionJson, "Role").String()
					var role domain.Role

					role.Name = gjson.Get(roleJson, "Name").String()
					role.RoleArn = gjson.Get(roleJson, "RoleArn").String()
					role.Parent = gjson.Get(roleJson, "Parent").String()
					role.ParentRole = gjson.Get(roleJson, "ParentRole").String()

					account.Role = role
				} else {
					account.Role = domain.Role{}
				}

				account.IdpArn = gjson.Get(sessionJson, "IdpArn").String()
				account.Region = gjson.Get(sessionJson, "Region").String()
				account.SsoUrl = gjson.Get(sessionJson, "SsoUrl").String()
				account.Type = gjson.Get(sessionJson, "Type").String()
				account.ParentSessionId = gjson.Get(sessionJson, "ParentSessionId").String()
				account.ParentRole = gjson.Get(sessionJson, "ParentRole").String()

				session.Account = account
			} else {
				session.Account = domain.Account{}
			}

			sessions = append(sessions, session)
		}
	}

	configuration.Sessions = sessions
	return &configuration
}
