package model

type Configuration struct {
	SsoUrl               string
	ProxyConfiguration   ProxyConfiguration
	FederatedAwsSessions []FederatedAwsSession
	PlainAwsSessions     []PlainAwsSession
}

type ProxyConfiguration struct {
	ProxyProtocol string
	ProxyUrl string
	ProxyPort uint64
	Username string
	Password string
}

type FederatedAwsSession struct {
	Id           string
	Active       bool
	Loading      bool
	StartTime    string
	Account      FederatedAwsAccount
}

type PlainAwsSession struct {
	Id           string
	Active       bool
	Loading      bool
	StartTime    string
	Account      PlainAwsAccount
}

type FederatedAwsAccount struct {
	AccountNumber string
	Name          string
	Role          FederatedAwsRole
	IdpArn        string
	Region        string
	SsoUrl        string
	// Type            string
	// ParentSessionId string
	// ParentRole      string
}

type FederatedAwsRole struct {
	Name string
	Arn  string
	// Parent string
	// ParentRole string
}

type PlainAwsAccount struct {
	AccountNumber       string
	Name                string
	Region              string
	User                string
	AwsAccessKeyId      string
	AwsSecretAccessKey  string
	MfaDevice           string
}
