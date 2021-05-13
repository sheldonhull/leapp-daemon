package main

import (
  "github.com/google/uuid"
  "leapp_daemon/domain/session"
  "leapp_daemon/infrastructure/encryption"
  "leapp_daemon/infrastructure/file_system"
  "leapp_daemon/infrastructure/logging"
  "leapp_daemon/interface/repository"
  "leapp_daemon/use_case"
)

func main() {
	// Test MFA
	//testMFA()

	// ======= Deferred functions ========
	//defer logging.CloseLogFile()
	//defer timer.Close()

	/*
	// Check and create config file
	_, err := configuration.ReadConfiguration()
	// TODO: check the nature of the error: if is no such file is ok, otherwise it must be panicked
	if err != nil {
		fmt.Printf("%+v", err)
		err = configuration.CreateConfiguration()
		if err != nil {
			logging.Entry().Error(err)
			panic(err)
		}
	}
	 */

	// ======== Sessions Timer ========
	//timer.Initialize(1, use_case.RotateAllSessionsCredentials)

	// ======== WebSocket Hub ========
	//go websocket.Hub.Run()

	// ======== REST API Server ========
	//eng := engine.Engine()
	//eng.Serve(8080)

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

  plainAwsSessions := config.PlainAwsSessions

  logging.Info(config)

  plainAwsSessionFacade := session.GetPlainAwsSessionsFacade()
  plainAwsSessionFacade.SetPlainAwsSessions(plainAwsSessions)

  plainAwsSessionFacade.Subscribe(&use_case.SessionsWriter{
    ConfigurationRepository: &fileConfigurationRepository,
  })

  fakePlainAwsSession := session.PlainAwsSession{
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
  }
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
