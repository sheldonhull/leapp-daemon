package session

import (
	"fmt"
	"leapp_daemon/infrastructure/http/http_error"
	"sync"
)

var federatedAlibabaSessionsFacadeSingleton *federatedAlibabaSessionsFacade
var federatedAlibabaSessionsFacadeLock sync.Mutex
var federatedAlibabaSessionsLock sync.Mutex

type FederatedAlibabaSessionsObserver interface {
	UpdateFederatedAlibabaSessions(oldFederatedAlibabaSessions []FederatedAlibabaSession, newFederatedAlibabaSessions []FederatedAlibabaSession) error
}

type federatedAlibabaSessionsFacade struct {
	federatedAlibabaSessions []FederatedAlibabaSession
	observers                []FederatedAlibabaSessionsObserver
}

func GetFederatedAlibabaSessionsFacade() *federatedAlibabaSessionsFacade {
	federatedAlibabaSessionsFacadeLock.Lock()
	defer federatedAlibabaSessionsFacadeLock.Unlock()

	if federatedAlibabaSessionsFacadeSingleton == nil {
		federatedAlibabaSessionsFacadeSingleton = &federatedAlibabaSessionsFacade{
			federatedAlibabaSessions: make([]FederatedAlibabaSession, 0),
		}
	}

	return federatedAlibabaSessionsFacadeSingleton
}

func (fac *federatedAlibabaSessionsFacade) Subscribe(observer FederatedAlibabaSessionsObserver) {
	fac.observers = append(fac.observers, observer)
}

func (fac *federatedAlibabaSessionsFacade) GetFederatedAlibabaSessions() []FederatedAlibabaSession {
	return fac.federatedAlibabaSessions
}

func (fac *federatedAlibabaSessionsFacade) SetFederatedAlibabaSessions(federatedAlibabaSessions []FederatedAlibabaSession) {
	fac.federatedAlibabaSessions = federatedAlibabaSessions
}

func (fac *federatedAlibabaSessionsFacade) AddFederatedAlibabaSession(federatedAlibabaSession FederatedAlibabaSession) error {
	federatedAlibabaSessionsLock.Lock()
	defer federatedAlibabaSessionsLock.Unlock()

	oldFederatedAlibabaSessions := fac.GetFederatedAlibabaSessions()
	newFederatedAlibabaSessions := make([]FederatedAlibabaSession, 0)

	for i := range oldFederatedAlibabaSessions {
		newFederatedAlibabaSession := oldFederatedAlibabaSessions[i]
		newFederatedAlibabaSessionAccount := *oldFederatedAlibabaSessions[i].Account
		newFederatedAlibabaSession.Account = &newFederatedAlibabaSessionAccount
		newFederatedAlibabaSessions = append(newFederatedAlibabaSessions, newFederatedAlibabaSession)
	}

	for _, sess := range newFederatedAlibabaSessions {
		if federatedAlibabaSession.Id == sess.Id {
			return http_error.NewConflictError(fmt.Errorf("a FederatedAlibabaSession with id " + federatedAlibabaSession.Id +
				" is already present"))
		}

		/*if federatedAlibabaSession.Alias == sess.Alias {
			return http_error.NewUnprocessableEntityError(fmt.Errorf("a session with the same alias " +
				"is already present"))
		}*/
	}

	newFederatedAlibabaSessions = append(newFederatedAlibabaSessions, federatedAlibabaSession)

	err := fac.updateState(newFederatedAlibabaSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *federatedAlibabaSessionsFacade) RemoveFederatedAlibabaSession(id string) error {
	federatedAlibabaSessionsLock.Lock()
	defer federatedAlibabaSessionsLock.Unlock()

	oldFederatedAlibabaSessions := fac.GetFederatedAlibabaSessions()
	newFederatedAlibabaSessions := make([]FederatedAlibabaSession, 0)

	for i := range oldFederatedAlibabaSessions {
		newFederatedAlibabaSession := oldFederatedAlibabaSessions[i]
		newFederatedAlibabaSessionAccount := *oldFederatedAlibabaSessions[i].Account
		newFederatedAlibabaSession.Account = &newFederatedAlibabaSessionAccount
		newFederatedAlibabaSessions = append(newFederatedAlibabaSessions, newFederatedAlibabaSession)
	}

	for i, sess := range newFederatedAlibabaSessions {
		if sess.Id == id {
			newFederatedAlibabaSessions = append(newFederatedAlibabaSessions[:i], newFederatedAlibabaSessions[i+1:]...)
			break
		}
	}

	if len(fac.GetFederatedAlibabaSessions()) == len(newFederatedAlibabaSessions) {
		return http_error.NewNotFoundError(fmt.Errorf("plain Alibaba session with id %s not found", id))
	}

	err := fac.updateState(newFederatedAlibabaSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *federatedAlibabaSessionsFacade) GetFederatedAlibabaSessionById(id string) (*FederatedAlibabaSession, error) {
	for _, federatedAlibabaSession := range fac.GetFederatedAlibabaSessions() {
		if federatedAlibabaSession.Id == id {
			return &federatedAlibabaSession, nil
		}
	}
	return nil, http_error.NewNotFoundError(fmt.Errorf("plain Alibaba session with id %s not found", id))
}

func (fac *federatedAlibabaSessionsFacade) SetFederatedAlibabaSessionStatusToPending(id string) error {
	federatedAlibabaSessionsLock.Lock()
	defer federatedAlibabaSessionsLock.Unlock()

	federatedAlibabaSession, err := fac.GetFederatedAlibabaSessionById(id)
	if err != nil {
		return err
	}

	if !(federatedAlibabaSession.Status == NotActive) {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("plain Alibaba session with id " + id + "cannot be started because it's in pending or active state"))
	}

	oldFederatedAlibabaSessions := fac.GetFederatedAlibabaSessions()
	newFederatedAlibabaSessions := make([]FederatedAlibabaSession, 0)

	for i := range oldFederatedAlibabaSessions {
		newFederatedAlibabaSession := oldFederatedAlibabaSessions[i]
		newFederatedAlibabaSessionAccount := *oldFederatedAlibabaSessions[i].Account
		newFederatedAlibabaSession.Account = &newFederatedAlibabaSessionAccount
		newFederatedAlibabaSessions = append(newFederatedAlibabaSessions, newFederatedAlibabaSession)
	}

	for i, session := range newFederatedAlibabaSessions {
		if session.Id == id {
			newFederatedAlibabaSessions[i].Status = Pending
		}
	}

	err = fac.updateState(newFederatedAlibabaSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *federatedAlibabaSessionsFacade) SetFederatedAlibabaSessionStatusToActive(id string) error {
	federatedAlibabaSessionsLock.Lock()
	defer federatedAlibabaSessionsLock.Unlock()

	federatedAlibabaSession, err := fac.GetFederatedAlibabaSessionById(id)
	if err != nil {
		return err
	}

	if !(federatedAlibabaSession.Status == Pending) {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("plain Alibaba session with id " + id + "cannot be started because it's not in pending state"))
	}

	oldFederatedAlibabaSessions := fac.GetFederatedAlibabaSessions()
	newFederatedAlibabaSessions := make([]FederatedAlibabaSession, 0)

	for i := range oldFederatedAlibabaSessions {
		newFederatedAlibabaSession := oldFederatedAlibabaSessions[i]
		newFederatedAlibabaSessionAccount := *oldFederatedAlibabaSessions[i].Account
		newFederatedAlibabaSession.Account = &newFederatedAlibabaSessionAccount
		newFederatedAlibabaSessions = append(newFederatedAlibabaSessions, newFederatedAlibabaSession)
	}

	err = fac.updateState(newFederatedAlibabaSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *federatedAlibabaSessionsFacade) updateState(newState []FederatedAlibabaSession) error {
	oldFederatedAlibabaSessions := fac.GetFederatedAlibabaSessions()
	fac.SetFederatedAlibabaSessions(newState)

	for _, observer := range fac.observers {
		err := observer.UpdateFederatedAlibabaSessions(oldFederatedAlibabaSessions, newState)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
func CreateFederatedAlibabaSession(sessionContainer Container, name string, accountNumber string, roleName string, roleArn string, idpArn string,
	region string, ssoUrl string, profile string) error {

	sessions, err := sessionContainer.GetFederatedAlibabaSessions()
	if err != nil {
		return err
	}

	for _, session := range sessions {
		account := session.Account
		if account.AccountNumber == accountNumber && account.Role.Name == roleName {
			err = http_error2.NewUnprocessableEntityError(fmt.Errorf("an account with the same account number and " +
				"role name is already present"))
			return err
		}
	}

	role := FederatedAlibabaRole{
		Name: roleName,
		Arn:  roleArn,
	}

	federatedAlibabaAccount := FederatedAlibabaAccount{
		AccountNumber: accountNumber,
		Name:          name,
		Role:          &role,
		IdpArn:        idpArn,
		Region:        region,
		SsoUrl:        ssoUrl,
	}

	uuidString := uuid.New().String()
	uuidString = strings.Replace(uuidString, "-", "", -1)

	namedProfileId, err := named_profile.CreateNamedProfile(sessionContainer, profile)
	if err != nil {
		return err
	}

	session := FederatedAlibabaSession{
		Id:        uuidString,
		Status:    NotActive,
		StartTime: "",
		Account:   &federatedAlibabaAccount,
		Profile:   namedProfileId,
	}

	err = sessionContainer.SetFederatedAlibabaSessions(append(sessions, &session))
	if err != nil { return err }

	return nil
}

func GetFederatedAlibabaSession(sessionContainer Container, id string) (*FederatedAlibabaSession, error) {
	sessions, err := sessionContainer.GetFederatedAlibabaSessions()
	if err != nil {
		return nil, err
	}

	for index, _ := range sessions {
		if sessions[index].Id == id {
			return sessions[index], nil
		}
	}

	return nil, http_error2.NewNotFoundError(fmt.Errorf("No session found with id:" + id))
}

func ListFederatedAlibabaSession(sessionContainer Container, query string) ([]*FederatedAlibabaSession, error) {
	sessions, err := sessionContainer.GetFederatedAlibabaSessions()
	if err != nil {
		return nil, err
	}

	filteredList := make([]*FederatedAlibabaSession, 0)

	if query == "" {
		return append(filteredList, sessions...), nil
	} else {
		for _, session := range sessions {
			if  strings.Contains(session.Id, query) ||
				strings.Contains(session.Profile, query) ||
				strings.Contains(session.Account.Name, query) ||
				strings.Contains(session.Account.IdpArn, query) ||
				strings.Contains(session.Account.SsoUrl, query) ||
				strings.Contains(session.Account.Region, query) ||
				strings.Contains(session.Account.AccountNumber, query) ||
				strings.Contains(session.Account.Role.Name, query) ||
				strings.Contains(session.Account.Role.Arn, query) {

				filteredList = append(filteredList, session)
			}
		}

		return filteredList, nil
	}
}

func UpdateFederatedAlibabaSession(sessionContainer Container, id string, name string, accountNumber string, roleName string, roleArn string, idpArn string,
	region string, ssoUrl string, profile string) error {

	sessions, err := sessionContainer.GetFederatedAlibabaSessions()
	if err != nil { return err }

	found := false
	for index := range sessions {
		if sessions[index].Id == id {
			namedProfileId, err := named_profile.EditNamedProfile(sessionContainer, sessions[index].Profile, profile)
			if err != nil { return err }

			sessions[index].Profile = namedProfileId
			sessions[index].Account = &FederatedAlibabaAccount{
				AccountNumber: accountNumber,
				Name:          name,
				Region:        region,
				IdpArn: 	   idpArn,
				SsoUrl:        ssoUrl,
			}

			sessions[index].Account.Role = &FederatedAlibabaRole{
				Name: roleName,
				Arn:  roleArn,
			}

			found = true
		}
	}

	if found == false {
		err = http_error2.NewNotFoundError(fmt.Errorf("federated AWS session with id " + id + " not found"))
		return err
	}

	err = sessionContainer.SetFederatedAlibabaSessions(sessions)
	if err != nil { return err }

	return nil
}

func DeleteFederatedAlibabaSession(sessionContainer Container, id string) error {
	sessions, err := sessionContainer.GetFederatedAlibabaSessions()
	if err != nil {
		return err
	}

	found := false
	for index := range sessions {
		if sessions[index].Id == id {
			sessions = append(sessions[:index], sessions[index+1:]...)
			found = true
			break
		}
	}

	if found == false {
		err = http_error2.NewNotFoundError(fmt.Errorf("federated AWS session with id " + id + " not found"))
		return err
	}

	err = sessionContainer.SetFederatedAlibabaSessions(sessions)
	if err != nil {
		return err
	}

	return nil
}

func StartFederatedAlibabaSession(sessionContainer Container, id string) error {
	sess, err := GetFederatedAlibabaSession(sessionContainer, id)
	if err != nil {
		return err
	}

	println("Rotating session with id", sess.Id)
	err = sess.Rotate(nil)
	if err != nil { return err }

	return nil
}

func StopFederatedAlibabaSession(sessionContainer Container, id string) error {
	sess, err := GetFederatedAlibabaSession(sessionContainer, id)
	if err != nil {
		return err
	}

	sess.Status = NotActive
	return nil
}
*/
