package use_case

import "golang.org/x/oauth2"

type FileSystem interface {
	DoesFileExist(path string) bool
	GetHomeDir() (string, error)
	ReadFile(path string) ([]byte, error)
	RemoveFile(path string) error
	WriteToFile(path string, data []byte) error
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

type Environment interface {
	GenerateUuid() string
}
