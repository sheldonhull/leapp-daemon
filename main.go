package main

import (
	"leapp_daemon/logging"
	"leapp_daemon/rest_api/engine"
)

func main() {
	// TODO: create specific accounts models
	// TODO: implement keychain service
	// TODO: implement plain account creation service
	// TODO: implement tmp credentials generation flow
	// TODO: implement plain account strategy (with STS)
	// TODO: test suite

	// TODO: how to protect REST API from outside?
	// TODO: integrate with gRPC plugins

	// TODO: documentation

	// TODO: add authentication
	// TODO: add HTTPS
	// TODO: are there different web/application servers to host Gin application?
	defer logging.CloseLogFile()
	eng := engine.Engine()
	eng.Serve(8080)

	/*services.CreateConfiguration()
	c, err := services.ReadConfiguration()
	if err != nil { log.Fatalln(err) }
	c.ProxyConfiguration.Username = `eric`
	services.UpdateConfiguration(c, false)
	c, _ = services.ReadConfiguration()
	log.Printf("configuration: %+v\n", *c)*/
}
