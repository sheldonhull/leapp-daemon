package main

import (
	"leapp_daemon/engine"
	"leapp_daemon/logging"
)

func main() {
	defer logging.CloseLogFile()
	eng := engine.Engine()
	eng.Serve(8080)

	//conf, _ := services.ReadConfiguration()
	//log.Printf("conf: %+v\n", *conf)

	//_ = services.CreateConfiguration()
}
