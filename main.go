package main

import (
	"leapp_daemon/infrastructure/http/engine"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/providers"
)

func main() {
	defer logging.CloseLogFile()

	prov := providers.NewProviders()

	config := ConfigurationBootstrap(prov)
	NamedProfilesBootstrap(prov, config)
	AwsPlainBootstrap(prov, config)
	GcpPlainBootstrap(prov, config)

	//timer.Initialize(1, use_case.RotateAllSessionsCredentials)
	//defer timer.Close()
	//go websocket.Hub.Run()

	eng := engine.Engine(prov)
	eng.Serve(8080)
}
