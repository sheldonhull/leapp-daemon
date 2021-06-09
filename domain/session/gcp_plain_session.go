package session

type GcpPlainSessionContainer interface {
	AddSession(GcpPlainSession) error
	GetAllSessions() ([]GcpPlainSession, error)
	RemoveSession(session GcpPlainSession) error
}

type GcpPlainSession struct {
	Id             string
	Name           string
	AccountId      string
	ProjectName    string
	NamedProfileId string
	Status         Status
	StartTime      string
	LastStopTime   string
}
