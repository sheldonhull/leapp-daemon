package providers

import (
	"leapp_daemon/use_case"
	"sync"
)

var awsPlainSessionActionsSingleton *use_case.AwsPlainSessionActions
var awsPlainSessionActionsMutex sync.Mutex

func (prov *Providers) GetAwsPlainSessionActions() *use_case.AwsPlainSessionActions {
	awsPlainSessionActionsMutex.Lock()
	defer awsPlainSessionActionsMutex.Unlock()

	if awsPlainSessionActionsSingleton == nil {
		awsPlainSessionActionsSingleton = &use_case.AwsPlainSessionActions{
			NamedProfilesActions: prov.GetNamedProfilesActions(),
			Environment:          prov.GetEnvironment(),
			Keychain:             prov.GetKeychain(),
		}
	}
	return awsPlainSessionActionsSingleton
}
