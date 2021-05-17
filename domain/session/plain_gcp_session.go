package session

type PlainGcpSessionContainer interface {
  AddPlainGcpSession(PlainGcpSession) error
  GetAllPlainGcpSessions() ([]PlainGcpSession, error)
  RemovePlainGcpSession(session PlainGcpSession) error
}

type PlainGcpSession struct {
	Id        string
	Status    Status
	StartTime string
	Account   *PlainGcpAccount
	Profile   string
}

type PlainGcpAccount struct {
	ServiceAccountName  string
	ProjectId           string
	APIKey              string
}
