package session

type GcpPlainSession struct {
	Id               string
	Name             string
	AccountId        string
	ProjectName      string
	CredentialsLabel string
	Status           Status
	StartTime        string // TODO: initialize in actions
	LastStopTime     string // TODO: initialize in actions
}
