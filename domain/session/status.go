package session

type Status int

const (
	NotActive Status = iota
	Pending
	Active
)
