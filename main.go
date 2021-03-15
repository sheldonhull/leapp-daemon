package main

import (
	"leapp_daemon/api/controller"
	"leapp_daemon/api/engine"
	"leapp_daemon/core/service"
	"leapp_daemon/shared/logging"
)

func main() {
	// Test MFA
	// testMFA()

	// ======= Deferred functions ========
	defer logging.CloseLogFile()
	defer service.CloseTimer()

	// Check and create config file
	_, err := service.ReadConfiguration()
	// TODO: check the nature of the error: if is no such file is ok, otherwise it must be panicked
	if err != nil {
		err = service.CreateConfiguration()
		if err != nil {
			logging.Entry().Error(err)
			panic(err)
		}
	}

	// ========   Global Timer    ========
	service.InitializeTimer(1, service.CheckAllSessions)

	// ======== WebSocket Channel ========
	go controller.Hub.Run()

	// ========     Api Server    ========
	eng := engine.Engine()
	eng.Serve(8080)
}

func testMFA() {
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
}