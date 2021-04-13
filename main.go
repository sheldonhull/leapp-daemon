package main

import (
  "leapp_daemon/infrastructure/http/engine"
  logging2 "leapp_daemon/infrastructure/logging"
  timer2 "leapp_daemon/infrastructure/timer"
  websocket2 "leapp_daemon/infrastructure/websocket"
  "leapp_daemon/use_case"
)

func main() {
	// Test MFA
	//testMFA()

	// ======= Deferred functions ========
	defer logging2.CloseLogFile()
	defer timer2.Close()

	/*
	// Check and create config file
	_, err := configuration.ReadConfiguration()
	// TODO: check the nature of the error: if is no such file is ok, otherwise it must be panicked
	if err != nil {
		fmt.Printf("%+v", err)
		err = configuration.CreateConfiguration()
		if err != nil {
			logging2.Entry().Error(err)
			panic(err)
		}
	}
	 */

	// ======== Sessions Timer ========
	timer2.Initialize(1, use_case.RotateAllSessionsCredentials)

	// ======== WebSocket Hub ========
	go websocket2.Hub.Run()

	// ======== REST API Server ========
	eng := engine.Engine()
	eng.Serve(8080)
}

/*
func testMFA() {
	config, err := configuration.ReadConfiguration()
	if err != nil {
		logging2.Info(err)
	}

	isMfaTokenRequired, err := session2.IsMfaRequiredForPlainAwsSession(config, "dc6b8f6015084ab885c00b5bc0fceb7b")

	if isMfaTokenRequired {
		var token = "014729"
		err = session2.StartPlainAwsSession(config, "dc6b8f6015084ab885c00b5bc0fceb7b", &token)
		if err != nil {
			logging2.Info(err)
		}
	} else {
		err = session2.StartPlainAwsSession(config,"dc6b8f6015084ab885c00b5bc0fceb7b", nil)
		if err != nil {
			logging2.Info(err)
		}
	}
}
 */
