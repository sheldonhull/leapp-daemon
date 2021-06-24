package use_case

import (
	"golang.org/x/oauth2"
	"leapp_daemon/domain/configuration"
	"leapp_daemon/domain/named_profile"
	"leapp_daemon/domain/session"
	"leapp_daemon/interface/repository"
)

type FileSystem interface {
	DoesFileExist(path string) bool
	GetHomeDir() (string, error)
}

type Environment interface {
	GenerateUuid() string
	GetTime() string
}

type Keychain interface {
	DoesSecretExist(label string) (bool, error)
	GetSecret(label string) (string, error)
	DeleteSecret(label string) error
	SetSecret(secret string, label string) error
}

type GcpApi interface {
	GetOauthUrl() (string, error)
	GetOauthToken(authCode string) (*oauth2.Token, error)
	GetCredentials(oauthToken *oauth2.Token) string
}

type ConfigurationRepository interface {
	CreateConfiguration(configuration.Configuration) error
	GetConfiguration() (configuration.Configuration, error)
	UpdateConfiguration(configuration.Configuration) error
}

type GcpConfigurationRepository interface {
	DoesGcloudConfigFolderExist() (bool, error)
	CreateConfiguration(account string, project string) error
	RemoveConfiguration() error
	ActivateConfiguration() error
	DeactivateConfiguration() error
	WriteDefaultCredentials(credentialsJson string) error
	RemoveDefaultCredentials() error
	WriteCredentialsToDb(accountId string, credentialsJson string) error
	RemoveCredentialsFromDb(accountId string) error
	RemoveAccessTokensFromDb(accountId string) error
}

type AwsConfigurationRepository interface {
	WriteCredentials(credentials []repository.AwsTempCredentials) error
}

type NamedProfilesFacade interface {
	GetNamedProfiles() []named_profile.NamedProfile
	GetNamedProfileById(id string) (named_profile.NamedProfile, error)
	GetNamedProfileByName(name string) (named_profile.NamedProfile, error)
	AddNamedProfile(namedProfile named_profile.NamedProfile) error
}

type NamedProfilesActionsInterface interface {
	GetOrCreateNamedProfile(profileName string) (named_profile.NamedProfile, error)
}

type AwsIamUserSessionsFacade interface {
	Subscribe(observer session.AwsIamUserSessionsObserver)
	GetSessions() []session.AwsIamUserSession
	GetSessionById(sessionId string) (session.AwsIamUserSession, error)
	SetSessions(sessions []session.AwsIamUserSession)
	AddSession(session session.AwsIamUserSession) error
	RemoveSession(sessionId string) error
	EditSession(sessionId string, sessionName string, region string, accessKeyIdLabel string, secretKeyLabel string,
		sessionTokenLabel string, mfaDevice string, sessionTokenExpiration string, namedProfileId string) error
	SetSessionTokenExpiration(sessionId string, sessionTokenExpiration string) error
	StartingSession(sessionId string) error
	StartSession(sessionId string, startTime string) error
	StopSession(sessionId string, stopTime string) error
}

type GcpIamUserAccountOauthSessionsFacade interface {
	GetSessions() []session.GcpIamUserAccountOauthSession
	GetSessionById(sessionId string) (session.GcpIamUserAccountOauthSession, error)
	AddSession(session session.GcpIamUserAccountOauthSession) error
	StartSession(sessionId string, startTime string) error
	StopSession(sessionId string, stopTime string) error
	RemoveSession(sessionId string) error
	EditSession(sessionId string, name string, projectName string) error
}
