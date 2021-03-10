package main

import (
	"leapp_daemon/core/service/session"
	"leapp_daemon/shared/logging"
)

func main() {
	err := session.StartPlainAwsSession("5eda3711cc724380a9c6f39638dbb9db")
	if err != nil {
		logging.Info(err)
	}

	//go controller.Hub.Run()
	/*defer logging.CloseLogFile()
	eng := engine.Engine()
	eng.Serve(8080)*/
}