package main

import (
  "fmt"
  "leapp_daemon/domain/session"
  "leapp_daemon/infrastructure/encryption"
  "leapp_daemon/infrastructure/file_system"
  "leapp_daemon/infrastructure/http/engine"
  "leapp_daemon/infrastructure/logging"
  "leapp_daemon/interface/repository"
  "leapp_daemon/use_case"
)

func main() {

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

  plainAwsSessions := config.PlainAwsSessions

  plainAwsSessionFacade := session.GetPlainAwsSessionsFacade()
  plainAwsSessionFacade.SetPlainAwsSessions(plainAwsSessions)

  plainAwsSessionFacade.Subscribe(&use_case.SessionsWriter{
    ConfigurationRepository: &fileConfigurationRepository,
  })

  //timer.Initialize(1, use_case.RotateAllSessionsCredentials)
  //go websocket.Hub.Run()
  eng := engine.Engine()
  eng.Serve(8080)

  /*fakePlainAwsSession := session.PlainAwsSession{
      Id:        uuid.New().String(),
      Status:    0,
      StartTime: "",
      Account:   nil,
      Profile:   "",
    }

    err = plainAwsSessionFacade.AddPlainAwsSession(fakePlainAwsSession)
    if err != nil {
      logging.Entry().Error(err)
      panic(err)
    }*/
}

/*
func testMFA() {
	config, err := configuration.ReadConfiguration()
	if err != nil {
		logging.Info(err)
	}

	isMfaTokenRequired, err := session2.IsMfaRequiredForPlainAwsSession(config, "dc6b8f6015084ab885c00b5bc0fceb7b")

	if isMfaTokenRequired {
		var token = "014729"
		err = session2.StartPlainAwsSession(config, "dc6b8f6015084ab885c00b5bc0fceb7b", &token)
		if err != nil {
			logging.Info(err)
		}
	} else {
		err = session2.StartPlainAwsSession(config,"dc6b8f6015084ab885c00b5bc0fceb7b", nil)
		if err != nil {
			logging.Info(err)
		}
	}
}
 */
