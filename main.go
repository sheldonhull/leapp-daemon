package main

import (
	"leapp_daemon/engine"
	"leapp_daemon/logging"
)

func main() {
	//go controller.Hub.Run()
	defer logging.CloseLogFile()
	eng := engine.Engine()
	eng.Serve(8080)
}