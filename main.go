package main

import (
	"leapp_daemon/api/controller"
	"leapp_daemon/api/engine"
	"leapp_daemon/core/service"
	"leapp_daemon/core/service/session"
	"leapp_daemon/shared/logging"
)

func main() {
	isMfaTokenRequired, err := session.IsMfaRequiredForPlainAwsSession("5eda3711cc724380a9c6f39638dbb9db")

	if isMfaTokenRequired {
		err = session.StartPlainAwsSession("5eda3711cc724380a9c6f39638dbb9db", "")
		if err != nil {
			logging.Info(err)
		}
	} else {
		err = session.StartPlainAwsSession("5eda3711cc724380a9c6f39638dbb9db", "")
		if err != nil {
			logging.Info(err)
		}
	}

	defer service.CloseTimer()
	service.InitiliazeTimer(1, service.CheckAllSessions)


	go controller.Hub.Run()
	defer logging.CloseLogFile()
	eng := engine.Engine()
	eng.Serve(8080)
}