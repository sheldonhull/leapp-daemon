package aws

type AwsSessionStatus int

const (
	NotActive AwsSessionStatus = iota
	Pending
	Active
)
