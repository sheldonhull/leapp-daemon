package main

import (
  "fmt"
  "leapp_daemon/domain/named_profile"
  "leapp_daemon/domain/session"
  "leapp_daemon/infrastructure/encryption"
  "leapp_daemon/infrastructure/environment"
  "leapp_daemon/infrastructure/file_system"
  "leapp_daemon/infrastructure/http/engine"
  "leapp_daemon/infrastructure/keychain"
  "leapp_daemon/infrastructure/logging"
  "leapp_daemon/interface/gcp"
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
  plainAwsSessionFacade := session.GetPlainAwsSessionsFacade()

  plainAwsSessions := config.PlainAwsSessions
  logging.Info(fmt.Sprintf("%+v", plainAwsSessions))
  plainAwsSessionFacade.SetPlainAwsSessions(plainAwsSessions)

  plainAwsSessionFacade.Subscribe(&use_case.AwsSessionsWriter{
    ConfigurationRepository: &fileConfigurationRepository,
  })

  keyChain := &keychain.Keychain{}
  plainAwsSessionFacade.Subscribe(&use_case.AwsCredentialsApplier{
    FileSystem: fileSystem,
    Keychain:   keyChain,
  })

  // GCP PLAIN
  gcpPlainSessionFacade := session.GetGcpPlainSessionFacade()

  gcpPlainSessions := config.GcpPlainSessions
  logging.Info(fmt.Sprintf("%+v", gcpPlainSessions))
  gcpPlainSessionFacade.SetSessions(gcpPlainSessions)

  gcpPlainSessionFacade.Subscribe(&use_case.GcpSessionsWriter{
    ConfigurationRepository: &fileConfigurationRepository,
  })

  gcpPlainSessionFacade.Subscribe(&use_case.GcpCredentialsApplier{
    Repository: repository.GcloudConfigurationRepository{
      FileSystem:        fileSystem,
      Environment:       &environment.Environment{},
      CredentialsTable:  &gcp.CredentialsTable{},
      AccessTokensTable: &gcp.AccessTokensTable{},
    },
    Keychain: keyChain,
  })

  // NAMED PROFILES
  namedProfilesFacade := named_profile.GetNamedProfilesFacade()

  namedProfiles := config.NamedProfiles
  logging.Info(fmt.Sprintf("%+v", namedProfiles))
  namedProfilesFacade.SetNamedProfiles(namedProfiles)

  namedProfilesFacade.Subscribe(&use_case.NamedProfilesWriter{
    ConfigurationRepository: &fileConfigurationRepository,
  })

  // TODO: subscribe observer that reads session token from keychain and writes it down into credentials file

  //timer.Initialize(1, use_case.RotateAllSessionsCredentials)
  //go websocket.Hub.Run()
  eng := engine.Engine()
  eng.Serve(8080)
}
