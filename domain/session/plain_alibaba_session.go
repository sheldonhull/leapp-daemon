package session

type PlainAlibabaSessionContainer interface {
	AddPlainAlibabaSession(PlainAlibabaSession) error
	GetAllPlainAlibabaSessions() ([]PlainAlibabaSession, error)
	RemovePlainAlibabaSession(session PlainAlibabaSession) error
}

type PlainAlibabaSession struct {
	Id      string
	Alias   string
	Status  Status
	Account *PlainAlibabaAccount
}

type PlainAlibabaAccount struct {
	Region         string
	NamedProfileId string
}
