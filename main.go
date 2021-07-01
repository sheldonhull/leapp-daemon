package main

import (
	"fmt"
	"leapp_daemon/domain/constant"
	"leapp_daemon/domain/named_profile"
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/encryption"
	"leapp_daemon/infrastructure/file_system"
	"leapp_daemon/infrastructure/http/engine"
	"leapp_daemon/infrastructure/keychain"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/interface/repository"
	"leapp_daemon/use_case"
)

func main() {
	// TODO: add logging state observers

	defer logging.CloseLogFile()
	//defer timer.Close()

	fileSystem := &file_system.FileSystem{}

	fileConfigurationRepository := repository.FileConfigurationRepository{
		FileSystem: fileSystem,
		Encryption: &encryption.Encryption{},
	}

	config, err := fileConfigurationRepository.GetConfiguration()
	if err != nil {
		logging.Entry().Error(err)
		panic(err)
	}

	homeDir, err := fileSystem.GetHomeDir()
	if err != nil {
		logging.Entry().Error(err)
		panic(err)
	}

	configurationFileBackupPath := fmt.Sprintf("%s/%s", homeDir, ".Leapp/Leapp-lock.json")
	doesConfigurationFileExist := fileSystem.DoesFileExist(configurationFileBackupPath)

	credentialsFilePath := fmt.Sprintf("%s/%s", homeDir, constant.AwsCredentialsFilePath)
	doesCredentialsFileExist := fileSystem.DoesFileExist(credentialsFilePath)

	credentialsFileBackupPath := fmt.Sprintf("%s/%s", homeDir, ".aws/credentials.leapp.bkp")
	doesCredentialsFileBackupExist := fileSystem.DoesFileExist(credentialsFileBackupPath)

	if !doesConfigurationFileExist && doesCredentialsFileExist && !doesCredentialsFileBackupExist {
		err = fileSystem.RenameFile(credentialsFilePath, credentialsFileBackupPath)
		if err != nil {
			logging.Entry().Error(err)
			panic(err)
		}
	}

	plainAwsSessions := config.PlainAwsSessions

	logging.Info(fmt.Sprintf("%+v", plainAwsSessions))

	plainAwsSessionFacade := session.GetPlainAwsSessionsFacade()
	plainAwsSessionFacade.SetPlainAwsSessions(plainAwsSessions)

	plainAwsSessionFacade.Subscribe(&use_case.SessionsWriter{
		ConfigurationRepository: &fileConfigurationRepository,
	})

	plainAwsSessionFacade.Subscribe(&use_case.AwsCredentialsApplier{
		FileSystem: fileSystem,
		Keychain:   &keychain.Keychain{},
	})

	namedProfiles := config.NamedProfiles

	logging.Info(fmt.Sprintf("%+v", namedProfiles))

	namedProfilesFacade := named_profile.GetNamedProfilesFacade()
	namedProfilesFacade.SetNamedProfiles(namedProfiles)

	namedProfilesFacade.Subscribe(&use_case.NamedProfilesWriter{
		ConfigurationRepository: &fileConfigurationRepository,
	})

	// TODO: subscribe observer that reads session token from keychain and writes it down into credentials file

	plainAlibabaSessions := config.PlainAlibabaSessions
	logging.Info(fmt.Sprintf("%+v", plainAlibabaSessions))
	plainAlibabaSessionFacade := session.GetPlainAlibabaSessionsFacade()
	plainAlibabaSessionFacade.SetSessions(plainAlibabaSessions)
	plainAlibabaSessionFacade.Subscribe(&use_case.SessionsWriter{
		ConfigurationRepository: &fileConfigurationRepository,
	})
	plainAlibabaSessionFacade.Subscribe(&use_case.AlibabaCredentialsApplier{
		FileSystem: fileSystem,
		Keychain:   &keychain.Keychain{},
	})

	federatedAlibabaSessions := config.FederatedAlibabaSessions
	logging.Info(fmt.Sprintf("%+v", federatedAlibabaSessions))
	federatedAlibabaSessionFacade := session.GetFederatedAlibabaSessionsFacade()
	federatedAlibabaSessionFacade.SetSessions(federatedAlibabaSessions)
	federatedAlibabaSessionFacade.Subscribe(&use_case.SessionsWriter{
		ConfigurationRepository: &fileConfigurationRepository,
	})
	federatedAlibabaSessionFacade.Subscribe(&use_case.AlibabaCredentialsApplier{
		FileSystem: fileSystem,
		Keychain:   &keychain.Keychain{},
	})

	trustedAlibabaSessions := config.TrustedAlibabaSessions
	logging.Info(fmt.Sprintf("%+v", trustedAlibabaSessions))
	trustedAlibabaSessionFacade := session.GetTrustedAlibabaSessionsFacade()
	trustedAlibabaSessionFacade.SetSessions(trustedAlibabaSessions)
	trustedAlibabaSessionFacade.Subscribe(&use_case.SessionsWriter{
		ConfigurationRepository: &fileConfigurationRepository,
	})
	
	trustedAlibabaSessionFacade.Subscribe(&use_case.AlibabaCredentialsApplier{
		FileSystem: fileSystem,
		Keychain:   &keychain.Keychain{},
	})

	//timer.Initialize(1, use_case.RotateAllSessionsCredentials)
	//go websocket.Hub.Run()
	eng := engine.Engine()
	eng.Serve(8080)
}
