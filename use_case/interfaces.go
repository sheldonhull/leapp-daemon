package use_case

import (
	"golang.org/x/oauth2"
	"leapp_daemon/domain/configuration"
	"leapp_daemon/domain/named_profile"
	"leapp_daemon/domain/session"
	"time"
)

type FileSystem interface {
	DoesFileExist(path string) bool
	GetHomeDir() (string, error)
}

type Environment interface {
	GenerateUuid() string
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
	CreateConfiguration(configurationName string, account string, project string) error
	RemoveConfiguration(configurationName string) error
	ActivateConfiguration() error
	DeactivateConfiguration() error
	WriteDefaultCredentials(credentialsJson string) error
	RemoveDefaultCredentials() error
	WriteCredentialsToDb(accountId string, credentialsJson string) error
	RemoveCredentialsFromDb(accountId string) error
	RemoveAccessTokensFromDb(accountId string) error
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

type AwsPlainSessionsFacade interface {
	Subscribe(observer session.AwsPlainSessionsObserver)
	GetSessions() []session.AwsPlainSession
	SetSessions(sessions []session.AwsPlainSession)
	AddSession(session session.AwsPlainSession) error
	RemoveSession(id string) error
	GetSessionById(id string) (*session.AwsPlainSession, error)
	SetSessionStatusToPending(id string) error
	SetSessionStatusToActive(id string) error
	SetSessionTokenExpiration(sessionId string, sessionTokenExpiration time.Time) error
}

type GcpPlainSessionsFacade interface {
	GetSessions() []session.GcpPlainSession
	GetSessionById(id string) (session.GcpPlainSession, error)
	AddSession(session session.GcpPlainSession) error
	SetSessionStatus(sessionId string, status session.Status) error
	RemoveSession(id string) error
	EditSession(sessionId string, name string, projectName string, profileId string) error
}
