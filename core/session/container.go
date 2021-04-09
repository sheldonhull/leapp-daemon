package session

type Container interface {
	GetPlainAwsSessions() ([]*PlainAwsSession, error)
	SetPlainAwsSessions([]*PlainAwsSession) error

	GetFederatedAwsSessions() ([]*FederatedAwsSession, error)
	SetFederatedAwsSessions([]*FederatedAwsSession) error

	GetNamedProfiles() ([]*NamedProfile, error)
	SetNamedProfiles([]*NamedProfile) error
}
