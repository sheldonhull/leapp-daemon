package domain

type Configuration struct {
	SsoUrl string
	ProxyConfiguration ProxyConfiguration
	Sessions []Session
}

type ProxyConfiguration struct {
	ProxyProtocol string
	ProxyUrl string
	proxyPort int
	Username string
	Password string
}

type Session struct {
	Id string
	Active bool
	Loading bool
	LastStopDate string
	Account Account
}

type Account struct {
	Id string
	Name string
	AccountNumber string
	Role Role
	IdpArn string
	Region string
	SsoUrl string
	Type string
	ParentSessionId string
	ParentRole string
}

type Role struct {
	Name string
	RoleArn string
	Parent string
	ParentRole string
}
