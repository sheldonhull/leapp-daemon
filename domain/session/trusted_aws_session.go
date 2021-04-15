package session

type TrustedAwsSession struct {
	Id        string
	Status    Status
	StartTime string
  ParentId  string
	Account   *TrustedAwsAccount
}

type TrustedAwsAccount struct {
	AccountNumber     string
	Name              string
	Role              *TrustedAwsRole
	Region            string
	// Type            string
	// ParentSessionId string
	// ParentRole      string
}

type TrustedAwsRole struct {
	Name string
	Arn  string
	// Parent string
	// ParentRole string
}

/*
func CreateTrusterAwsSession(AccountName string, AccountNumber string, RoleName string, Region string) error {

  sessions, err := sessionContainer.GetPlainAwsSessions()
  if err != nil {
    return err
  }

  for _, sess := range sessions {
    account := sess.Account
    if account.AccountNumber == accountNumber && account.User == user {
      err := http_error.NewUnprocessableEntityError(fmt.Errorf("a session with the same account number and user is already present"))
      return err
    }
  }

  plainAwsAccount := PlainAwsAccount{
    AccountNumber: accountNumber,
    Name:          name,
    Region:        region,
    User:          user,
    AwsAccessKeyId: awsAccessKeyId,
    AwsSecretAccessKey: awsSecretAccessKey,
    MfaDevice:     mfaDevice,

  }

  uuidString := uuid.New().String()
  uuidString = strings.Replace(uuidString, "-", "", -1)

  namedProfileId, err := CreateNamedProfile(sessionContainer, profile)
  if err != nil {
    return err
  }


  sess := PlainAwsSession{
    Id:        uuidString,
    Status:    NotActive,
    StartTime: "",
    Account:   &plainAwsAccount,
    Profile: namedProfileId,
  }

  err = sessionContainer.SetPlainAwsSessions(append(sessions, &sess))
  if err != nil {
    return err
  }

  return nil
}
*/
