package session

type PlainGcpSessionContainer interface {
	AddPlainGcpSession(PlainGcpSession) error
	GetAllPlainGcpSessions() ([]PlainGcpSession, error)
	RemovePlainGcpSession(session PlainGcpSession) error
}

type PlainGcpSession struct {
	Id           string
	Alias        string
	Status       Status
	StartTime    string
	LastStopTime string
	AccountId    string
}
