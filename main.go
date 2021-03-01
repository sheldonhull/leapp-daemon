package main

import (
	"leapp_daemon/controllers"
	"leapp_daemon/engine"
	"leapp_daemon/logging"
)

func main() {
	go controllers.Hub.Run()

	defer logging.CloseLogFile()
	eng := engine.Engine()
	eng.Serve(8080)
}