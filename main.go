package main

import (
	"leapp_daemon/api/engine"
	"leapp_daemon/core/configuration"
	"leapp_daemon/core/timer"
	"leapp_daemon/core/websocket"
	"leapp_daemon/logging"
	"leapp_daemon/service"
)

func main() {
	// Test MFA
	//testMFA()

	// ======= Deferred functions ========
	defer logging.CloseLogFile()
	defer timer.Close()

	// Check and create config file
	_, err := configuration.ReadConfiguration()
	// TODO: check the nature of the error: if is no such file is ok, otherwise it must be panicked
	if err != nil {
		err = configuration.CreateConfiguration()
		if err != nil {
			logging.Entry().Error(err)
			panic(err)
		}
	}

	// ======== Sessions Timer ========
	timer.Initialize(1, service.RotateAllSessionsCredentials)

	// ======== WebSocket Hub ========
	go websocket.Hub.Run()

	// ======== REST API Server ========
	eng := engine.Engine()
	eng.Serve(8080)
}

func testMFA() {
	isMfaTokenRequired, err := service.IsMfaRequiredForPlainAwsSession("dc6b8f6015084ab885c00b5bc0fceb7b")

	if isMfaTokenRequired {
		var token = "587023"
		err = service.StartPlainAwsSession("dc6b8f6015084ab885c00b5bc0fceb7b", &token)
		if err != nil {
			logging.Info(err)
		}
	} else {
		err = service.StartPlainAwsSession("dc6b8f6015084ab885c00b5bc0fceb7b", nil)
		if err != nil {
			logging.Info(err)
		}
	}
}