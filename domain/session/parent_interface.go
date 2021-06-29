package session

type ParentSession interface {
	GetId() string
	GetTypeString() string
}