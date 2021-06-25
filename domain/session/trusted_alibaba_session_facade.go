package session

import (
	"fmt"
	"leapp_daemon/infrastructure/http/http_error"
	"sync"
)

var trustedAlibabaSessionsFacadeSingleton *trustedAlibabaSessionsFacade
var trustedAlibabaSessionsFacadeLock sync.Mutex
var trustedAlibabaSessionsLock sync.Mutex

type TrustedAlibabaSessionsObserver interface {
	UpdateTrustedAlibabaSessions(oldTrustedAlibabaSessions []TrustedAlibabaSession, newTrustedAlibabaSessions []TrustedAlibabaSession) error
}

type trustedAlibabaSessionsFacade struct {
	trustedAlibabaSessions []TrustedAlibabaSession
	observers              []TrustedAlibabaSessionsObserver
}

func GetTrustedAlibabaSessionsFacade() *trustedAlibabaSessionsFacade {
	trustedAlibabaSessionsFacadeLock.Lock()
	defer trustedAlibabaSessionsFacadeLock.Unlock()

	if trustedAlibabaSessionsFacadeSingleton == nil {
		trustedAlibabaSessionsFacadeSingleton = &trustedAlibabaSessionsFacade{
			trustedAlibabaSessions: make([]TrustedAlibabaSession, 0),
		}
	}

	return trustedAlibabaSessionsFacadeSingleton
}

func (fac *trustedAlibabaSessionsFacade) Subscribe(observer TrustedAlibabaSessionsObserver) {
	fac.observers = append(fac.observers, observer)
}

func (fac *trustedAlibabaSessionsFacade) GetTrustedAlibabaSessions() []TrustedAlibabaSession {
	return fac.trustedAlibabaSessions
}

func (fac *trustedAlibabaSessionsFacade) SetTrustedAlibabaSessions(trustedAlibabaSessions []TrustedAlibabaSession) error {
	oldTrustedAlibabaSessions := fac.GetTrustedAlibabaSessions()
	fac.trustedAlibabaSessions = trustedAlibabaSessions

	for _, observer := range fac.observers {
		err := observer.UpdateTrustedAlibabaSessions(oldTrustedAlibabaSessions, trustedAlibabaSessions)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fac *trustedAlibabaSessionsFacade) AddTrustedAlibabaSession(trustedAlibabaSession TrustedAlibabaSession) error {
	trustedAlibabaSessionsLock.Lock()
	defer trustedAlibabaSessionsLock.Unlock()

	oldTrustedAlibabaSessions := fac.GetTrustedAlibabaSessions()
	newTrustedAlibabaSessions := make([]TrustedAlibabaSession, 0)

	for i := range oldTrustedAlibabaSessions {
		newTrustedAlibabaSession := oldTrustedAlibabaSessions[i]
		newTrustedAlibabaSessionAccount := *oldTrustedAlibabaSessions[i].Account
		newTrustedAlibabaSession.Account = &newTrustedAlibabaSessionAccount
		newTrustedAlibabaSessions = append(newTrustedAlibabaSessions, newTrustedAlibabaSession)
	}

	for _, sess := range newTrustedAlibabaSessions {
		if trustedAlibabaSession.Id == sess.Id {
			return http_error.NewConflictError(fmt.Errorf("a TrustedAlibabaSession with id " + trustedAlibabaSession.Id +
				" is already present"))
		}

		/*if trustedAlibabaSession.Alias == sess.Alias {
			return http_error.NewUnprocessableEntityError(fmt.Errorf("a session with the same alias " +
				"is already present"))
		}*/
	}

	newTrustedAlibabaSessions = append(newTrustedAlibabaSessions, trustedAlibabaSession)

	err := fac.updateState(newTrustedAlibabaSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *trustedAlibabaSessionsFacade) RemoveTrustedAlibabaSession(id string) error {
	trustedAlibabaSessionsLock.Lock()
	defer trustedAlibabaSessionsLock.Unlock()

	oldTrustedAlibabaSessions := fac.GetTrustedAlibabaSessions()
	newTrustedAlibabaSessions := make([]TrustedAlibabaSession, 0)

	for i := range oldTrustedAlibabaSessions {
		newTrustedAlibabaSession := oldTrustedAlibabaSessions[i]
		newTrustedAlibabaSessionAccount := *oldTrustedAlibabaSessions[i].Account
		newTrustedAlibabaSession.Account = &newTrustedAlibabaSessionAccount
		newTrustedAlibabaSessions = append(newTrustedAlibabaSessions, newTrustedAlibabaSession)
	}

	for i, sess := range newTrustedAlibabaSessions {
		if sess.Id == id {
			newTrustedAlibabaSessions = append(newTrustedAlibabaSessions[:i], newTrustedAlibabaSessions[i+1:]...)
			break
		}
	}

	if len(fac.GetTrustedAlibabaSessions()) == len(newTrustedAlibabaSessions) {
		return http_error.NewNotFoundError(fmt.Errorf("plain Alibaba session with id %s not found", id))
	}

	err := fac.updateState(newTrustedAlibabaSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *trustedAlibabaSessionsFacade) GetTrustedAlibabaSessionById(id string) (*TrustedAlibabaSession, error) {
	for _, trustedAlibabaSession := range fac.GetTrustedAlibabaSessions() {
		if trustedAlibabaSession.Id == id {
			return &trustedAlibabaSession, nil
		}
	}
	return nil, http_error.NewNotFoundError(fmt.Errorf("plain Alibaba session with id %s not found", id))
}

func (fac *trustedAlibabaSessionsFacade) SetTrustedAlibabaSessionById(newSession TrustedAlibabaSession) {
	allSessions := fac.GetTrustedAlibabaSessions()
	for i, trustedAlibabaSession := range allSessions {
		if trustedAlibabaSession.Id == newSession.Id {
			allSessions[i] = newSession
		}
	}
	fac.SetTrustedAlibabaSessions(allSessions)
}

func (fac *trustedAlibabaSessionsFacade) SetTrustedAlibabaSessionStatusToPending(id string) error {

	trustedAlibabaSessionsLock.Lock()
	defer trustedAlibabaSessionsLock.Unlock()

	trustedAlibabaSession, err := fac.GetTrustedAlibabaSessionById(id)
	if err != nil {
		return err
	}

	if !(trustedAlibabaSession.Status == NotActive) {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("plain Alibaba session with id " + id + "cannot be started because it's in pending or active state"))
	}

	oldTrustedAlibabaSessions := fac.GetTrustedAlibabaSessions()
	newTrustedAlibabaSessions := make([]TrustedAlibabaSession, 0)

	for i := range oldTrustedAlibabaSessions {
		newTrustedAlibabaSession := oldTrustedAlibabaSessions[i]
		newTrustedAlibabaSessionAccount := *oldTrustedAlibabaSessions[i].Account
		newTrustedAlibabaSession.Account = &newTrustedAlibabaSessionAccount
		newTrustedAlibabaSessions = append(newTrustedAlibabaSessions, newTrustedAlibabaSession)
	}

	for i, session := range newTrustedAlibabaSessions {
		if session.Id == id {
			newTrustedAlibabaSessions[i].Status = Pending
		}
	}

	err = fac.updateState(newTrustedAlibabaSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *trustedAlibabaSessionsFacade) SetTrustedAlibabaSessionStatusToActive(id string) error {
	trustedAlibabaSessionsLock.Lock()
	defer trustedAlibabaSessionsLock.Unlock()

	trustedAlibabaSession, err := fac.GetTrustedAlibabaSessionById(id)
	if err != nil {
		return err
	}

	if !(trustedAlibabaSession.Status == Pending) {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("plain Alibaba session with id " + id + "cannot be started because it's not in pending state"))
	}

	oldTrustedAlibabaSessions := fac.GetTrustedAlibabaSessions()
	newTrustedAlibabaSessions := make([]TrustedAlibabaSession, 0)

	for i := range oldTrustedAlibabaSessions {
		newTrustedAlibabaSession := oldTrustedAlibabaSessions[i]
		newTrustedAlibabaSessionAccount := *oldTrustedAlibabaSessions[i].Account
		newTrustedAlibabaSession.Account = &newTrustedAlibabaSessionAccount
		newTrustedAlibabaSessions = append(newTrustedAlibabaSessions, newTrustedAlibabaSession)
	}

	err = fac.updateState(newTrustedAlibabaSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *trustedAlibabaSessionsFacade) updateState(newState []TrustedAlibabaSession) error {
	oldTrustedAlibabaSessions := fac.GetTrustedAlibabaSessions()
	fac.SetTrustedAlibabaSessions(newState)

	for _, observer := range fac.observers {
		err := observer.UpdateTrustedAlibabaSessions(oldTrustedAlibabaSessions, newState)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
func CreateTrustedAlibabaSession(sessionContainer Container, name string, accountNumber string, roleName string, roleArn string, idpArn string,
	region string, ssoUrl string, profile string) error {

	sessions, err := sessionContainer.GetTrustedAlibabaSessions()
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

	role := TrustedAlibabaRole{
		Name: roleName,
		Arn:  roleArn,
	}

	trustedAlibabaAccount := TrustedAlibabaAccount{
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

	session := TrustedAlibabaSession{
		Id:        uuidString,
		Status:    NotActive,
		StartTime: "",
		Account:   &trustedAlibabaAccount,
		Profile:   namedProfileId,
	}

	err = sessionContainer.SetTrustedAlibabaSessions(append(sessions, &session))
	if err != nil { return err }

	return nil
}

func GetTrustedAlibabaSession(sessionContainer Container, id string) (*TrustedAlibabaSession, error) {
	sessions, err := sessionContainer.GetTrustedAlibabaSessions()
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

func ListTrustedAlibabaSession(sessionContainer Container, query string) ([]*TrustedAlibabaSession, error) {
	sessions, err := sessionContainer.GetTrustedAlibabaSessions()
	if err != nil {
		return nil, err
	}

	filteredList := make([]*TrustedAlibabaSession, 0)

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

func UpdateTrustedAlibabaSession(sessionContainer Container, id string, name string, accountNumber string, roleName string, roleArn string, idpArn string,
	region string, ssoUrl string, profile string) error {

	sessions, err := sessionContainer.GetTrustedAlibabaSessions()
	if err != nil { return err }

	found := false
	for index := range sessions {
		if sessions[index].Id == id {
			namedProfileId, err := named_profile.EditNamedProfile(sessionContainer, sessions[index].Profile, profile)
			if err != nil { return err }

			sessions[index].Profile = namedProfileId
			sessions[index].Account = &TrustedAlibabaAccount{
				AccountNumber: accountNumber,
				Name:          name,
				Region:        region,
				IdpArn: 	   idpArn,
				SsoUrl:        ssoUrl,
			}

			sessions[index].Account.Role = &TrustedAlibabaRole{
				Name: roleName,
				Arn:  roleArn,
			}

			found = true
		}
	}

	if found == false {
		err = http_error2.NewNotFoundError(fmt.Errorf("trusted AWS session with id " + id + " not found"))
		return err
	}

	err = sessionContainer.SetTrustedAlibabaSessions(sessions)
	if err != nil { return err }

	return nil
}

func DeleteTrustedAlibabaSession(sessionContainer Container, id string) error {
	sessions, err := sessionContainer.GetTrustedAlibabaSessions()
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
		err = http_error2.NewNotFoundError(fmt.Errorf("trusted AWS session with id " + id + " not found"))
		return err
	}

	err = sessionContainer.SetTrustedAlibabaSessions(sessions)
	if err != nil {
		return err
	}

	return nil
}

func StartTrustedAlibabaSession(sessionContainer Container, id string) error {
	sess, err := GetTrustedAlibabaSession(sessionContainer, id)
	if err != nil {
		return err
	}

	println("Rotating session with id", sess.Id)
	err = sess.Rotate(nil)
	if err != nil { return err }

	return nil
}

func StopTrustedAlibabaSession(sessionContainer Container, id string) error {
	sess, err := GetTrustedAlibabaSession(sessionContainer, id)
	if err != nil {
		return err
	}

	sess.Status = NotActive
	return nil
}
*/
