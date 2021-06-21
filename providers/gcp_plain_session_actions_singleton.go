package providers

import (
	"leapp_daemon/use_case"
	"sync"
)

var gcpPlainSessionActionsSingleton *use_case.GcpPlainSessionActions
var gcpPlainSessionActionsMutex sync.Mutex

func (prov *Providers) GetGcpPlainSessionActions() *use_case.GcpPlainSessionActions {
	gcpPlainSessionActionsMutex.Lock()
	defer gcpPlainSessionActionsMutex.Unlock()

	if gcpPlainSessionActionsSingleton == nil {
		gcpPlainSessionActionsSingleton = &use_case.GcpPlainSessionActions{
			GcpApi:                prov.GetGcpApi(),
			Environment:           prov.GetEnvironment(),
			Keychain:              prov.GetKeychain(),
			GcpPlainSessionFacade: prov.GetGcpPlainSessionFacade(),
		}
	}
	return gcpPlainSessionActionsSingleton
}
