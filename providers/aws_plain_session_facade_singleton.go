package providers

import (
	"leapp_daemon/domain/session"
	"sync"
)

var awsSessionsFacadeSingleton *session.AwsPlainSessionsFacade
var awsSessionsFacadeLock sync.Mutex

func (prov *Providers) GetAwsPlainSessionFacade() *session.AwsPlainSessionsFacade {
	awsSessionsFacadeLock.Lock()
	defer awsSessionsFacadeLock.Unlock()

	if awsSessionsFacadeSingleton == nil {
		awsSessionsFacadeSingleton = session.NewAwsPlainSessionsFacade()
	}
	return awsSessionsFacadeSingleton
}
