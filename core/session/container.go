package session

type Container interface {
	GetPlainAwsSessions() ([]*PlainAwsSession, error)
	SetPlainAwsSessions([]*PlainAwsSession) error
}
