package providers

import (
	"leapp_daemon/domain/session"
	"sync"
)

var gcpSessionsFacadeSingleton *session.GcpPlainSessionsFacade
var gcpSessionsFacadeLock sync.Mutex

func (prov *Providers) GetGcpPlainSessionFacade() *session.GcpPlainSessionsFacade {
	gcpSessionsFacadeLock.Lock()
	defer gcpSessionsFacadeLock.Unlock()

	if gcpSessionsFacadeSingleton == nil {
		gcpSessionsFacadeSingleton = session.NewGcpPlainSessionsFacade()
	}
	return gcpSessionsFacadeSingleton
}
