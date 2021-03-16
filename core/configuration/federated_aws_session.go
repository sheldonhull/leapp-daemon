package configuration

type FederatedAwsSession struct {
	Id           string
	Active       bool
	Loading      bool
	StartTime    string
	Account      *FederatedAwsAccount
}

type FederatedAwsAccount struct {
	AccountNumber string
	Name          string
	Role          *FederatedAwsRole
	IdpArn        string
	Region        string
	SsoUrl        string
	// Type            string
	// ParentSessionId string
	// ParentRole      string
}

type FederatedAwsRole struct {
	Name string
	Arn  string
	// Parent string
	// ParentRole string
}
