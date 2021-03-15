package main

import (
	"leapp_daemon/api/controller"
	"leapp_daemon/api/engine"
	"leapp_daemon/core/service"
	"leapp_daemon/shared/logging"
)

func main() {
	/*isMfaTokenRequired, err := session.IsMfaRequiredForPlainAwsSession("e2e5541af95c41f495087954d77b2e2d")

	if isMfaTokenRequired {
		err = session.StartPlainAwsSession("e2e5541af95c41f495087954d77b2e2d", "")
		if err != nil {
			logging.Info(err)
		}
	} else {
		err = session.StartPlainAwsSession("e2e5541af95c41f495087954d77b2e2d", "")
		if err != nil {
			logging.Info(err)
		}
	}*/

	defer service.CloseTimer()
	service.InitializeTimer(1, service.CheckAllSessions)

	go controller.Hub.Run()
	defer logging.CloseLogFile()
	eng := engine.Engine()
	eng.Serve(8080)
}