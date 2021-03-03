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
	var configuration domain.Configuration
	configuration.SsoUrl = gjson.Get(configurationJson, "SsoUrl").String()

	var proxyConfiguration domain.ProxyConfiguration
	proxyConfiguration.ProxyProtocol = gjson.Get(configurationJson, "ProxyConfiguration.ProxyProtocol").String()
	proxyConfiguration.ProxyUrl = gjson.Get(configurationJson, "ProxyConfiguration.ProxyUrl").String()
	proxyConfiguration.ProxyPort = gjson.Get(configurationJson, "ProxyConfiguration.ProxyPort").Uint()
	proxyConfiguration.Username = gjson.Get(configurationJson, "ProxyConfiguration.Username").String()
	proxyConfiguration.Password = gjson.Get(configurationJson, "ProxyConfiguration.Password").String()

	configuration.ProxyConfiguration = proxyConfiguration

	configuration.FederatedAwsAccountSessions = unrollFederatedAwsAccountSessions(configurationJson)

	configuration.PlainAwsAccountSessions = unrollPlainAwsAccountSessions(configurationJson)

	return &configuration
}

func unrollFederatedAwsAccountSessions(configurationJson string) []domain.FederatedAwsAccountSession {
	var federatedAwsAccountSessions []domain.FederatedAwsAccountSession
	federatedAwsAccountSessionsJsonArray := gjson.Get(configurationJson, "FederatedAwsAccountSessions").Array()

	if len(federatedAwsAccountSessionsJsonArray) > 0 {
		for _, sessionJsonValue := range federatedAwsAccountSessionsJsonArray {
			federatedAwsAccountSessionJson := sessionJsonValue.String()
			var federatedAwsAccountSession domain.FederatedAwsAccountSession

			federatedAwsAccountSession.Id = gjson.Get(federatedAwsAccountSessionJson, "Id").String()
			federatedAwsAccountSession.Active = gjson.Get(federatedAwsAccountSessionJson, "Active").Bool()
			federatedAwsAccountSession.Loading = gjson.Get(federatedAwsAccountSessionJson, "Loading").Bool()
			federatedAwsAccountSession.LastStopDate = gjson.Get(federatedAwsAccountSessionJson, "LastStopDate").String()

			federatedAwsAccountJsonValue := gjson.Get(federatedAwsAccountSessionJson, "Account").Value()

			if federatedAwsAccountJsonValue != nil {
				federatedAwsAccountJson := gjson.Get(federatedAwsAccountSessionJson, "Account").String()
				var federatedAwsAccount domain.FederatedAwsAccount

				federatedAwsAccount.AccountNumber = gjson.Get(federatedAwsAccountJson, "AccountNumber").String()
				federatedAwsAccount.Name = gjson.Get(federatedAwsAccountJson, "Name").String()

				awsRoleJsonValue := gjson.Get(federatedAwsAccountJson, "Role").Value()

				if awsRoleJsonValue != nil {
					awsRoleJson := gjson.Get(federatedAwsAccountJson, "Role").String()
					var awsRole domain.FederatedAwsRole

					awsRole.Name = gjson.Get(awsRoleJson, "Name").String()
					awsRole.Arn = gjson.Get(awsRoleJson, "Arn").String()
					//awsRole.Parent = gjson.Get(awsRoleJson, "Parent").String()
					//awsRole.ParentRole = gjson.Get(awsRoleJson, "ParentRole").String()

					federatedAwsAccount.Role = awsRole
				} else {
					federatedAwsAccount.Role = domain.FederatedAwsRole{}
				}

				federatedAwsAccount.IdpArn = gjson.Get(federatedAwsAccountJson, "IdpArn").String()
				federatedAwsAccount.Region = gjson.Get(federatedAwsAccountJson, "Region").String()
				federatedAwsAccount.SsoUrl = gjson.Get(federatedAwsAccountJson, "SsoUrl").String()
				//federatedAwsAccount.Type = gjson.Get(federatedAwsAccountJson, "Type").String()
				//federatedAwsAccount.ParentSessionId = gjson.Get(federatedAwsAccountJson, "ParentSessionId").String()
				//federatedAwsAccount.ParentRole = gjson.Get(federatedAwsAccountJson, "ParentRole").String()

				federatedAwsAccountSession.Account = federatedAwsAccount
			} else {
				federatedAwsAccountSession.Account = domain.FederatedAwsAccount{}
			}

			federatedAwsAccountSessions = append(federatedAwsAccountSessions, federatedAwsAccountSession)
		}
	}
	return federatedAwsAccountSessions
}

func unrollPlainAwsAccountSessions(configurationJson string) []domain.PlainAwsAccountSession {
	var plainAwsAccountSessions []domain.PlainAwsAccountSession
	plainAwsAccountSessionsJsonArray := gjson.Get(configurationJson, "PlainAwsAccountSessions").Array()

	if len(plainAwsAccountSessionsJsonArray) > 0 {
		for _, sessionJsonValue := range plainAwsAccountSessionsJsonArray {
			plainAwsAccountSessionJson := sessionJsonValue.String()
			var plainAwsAccountSession domain.PlainAwsAccountSession

			plainAwsAccountSession.Id = gjson.Get(plainAwsAccountSessionJson, "Id").String()
			plainAwsAccountSession.Active = gjson.Get(plainAwsAccountSessionJson, "Active").Bool()
			plainAwsAccountSession.Loading = gjson.Get(plainAwsAccountSessionJson, "Loading").Bool()
			plainAwsAccountSession.LastStopDate = gjson.Get(plainAwsAccountSessionJson, "LastStopDate").String()

			plainAwsAccountJsonValue := gjson.Get(plainAwsAccountSessionJson, "Account").Value()

			if plainAwsAccountJsonValue != nil {
				plainAwsAccountJson := gjson.Get(plainAwsAccountSessionJson, "Account").String()
				var plainAwsAccount domain.PlainAwsAccount

				plainAwsAccount.AccountNumber = gjson.Get(plainAwsAccountJson, "AccountNumber").String()
				plainAwsAccount.Name = gjson.Get(plainAwsAccountJson, "Name").String()
				plainAwsAccount.Region = gjson.Get(plainAwsAccountJson, "Region").String()
				plainAwsAccount.User = gjson.Get(plainAwsAccountJson, "User").String()
				plainAwsAccount.MfaDevice = gjson.Get(plainAwsAccountJson, "MfaDevice").String()

				plainAwsAccountSession.Account = plainAwsAccount
			} else {
				plainAwsAccountSession.Account = domain.PlainAwsAccount{}
			}

			plainAwsAccountSessions = append(plainAwsAccountSessions, plainAwsAccountSession)
		}
	}
	return plainAwsAccountSessions
}
