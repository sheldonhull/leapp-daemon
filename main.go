package main

import (
	"fmt"
	"leapp_daemon/infrastructure/http/engine"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/providers"
)

func main() {

	prov := providers.NewProviders()

	defer logging.CloseLogFile()
	//defer timer.Close()

	// TODO: create BootstrapActions calling Gcp/AWS-BootstrapActions

	fileSystem := prov.GetFileSystem()
	fileConfigurationRepository := prov.GetFileConfigurationRepository()

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

	//AWS Stuff
	credentialsFilePath := fmt.Sprintf("%s/%s", homeDir, ".aws/credentials")
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

	// AWS PLAIN
	awsPlainSessionFacade := prov.GetAwsPlainSessionFacade()
	awsPlainSessions := config.PlainAwsSessions
	awsPlainSessionFacade.SetSessions(awsPlainSessions)
	awsPlainSessionFacade.Subscribe(prov.GetAwsSessionWriter())
	awsPlainSessionFacade.Subscribe(prov.GetAwsCredentialsApplier())
	logging.Info(fmt.Sprintf("%+v", awsPlainSessions))

	// GCP PLAIN
	gcpPlainSessionFacade := prov.GetGcpPlainSessionFacade()
	gcpPlainSessions := config.GcpPlainSessions
	gcpPlainSessionFacade.SetSessions(gcpPlainSessions)
	gcpPlainSessionFacade.Subscribe(prov.GetGcpSessionWriter())
	gcpPlainSessionFacade.Subscribe(prov.GetGcpCredentialsApplier())
	logging.Info(fmt.Sprintf("%+v", gcpPlainSessions))

	// NAMED PROFILES
	namedProfilesFacade := prov.GetNamedProfilesFacade()
	namedProfiles := config.NamedProfiles
	namedProfilesFacade.SetNamedProfiles(namedProfiles)
	namedProfilesFacade.Subscribe(prov.GetNamedProfilesWriter())
	logging.Info(fmt.Sprintf("%+v", namedProfiles))

	//timer.Initialize(1, use_case.RotateAllSessionsCredentials)
	//go websocket.Hub.Run()
	eng := engine.Engine(prov)
	eng.Serve(8080)
}
