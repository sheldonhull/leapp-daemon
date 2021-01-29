package main

import (
	"leapp_daemon/services"
	"log"
)

func main() {
	// TODO: how to protect REST API from outside?
	// TODO: integrate with gRPC plugins
	// TODO: test suite

	// TODO: documentation

	// TODO: add authentication
	// TODO: add HTTPS
	// TODO: are there different web/application servers to host Gin application?
	//eng := engine.Engine()
	//eng.Serve(8080)
	services.CreateConfiguration()
	c, _ := services.ReadConfiguration()
	c.ProxyConfiguration.Username = `eric`
	//log.Printf("configuration: %+v\n", *c)
	services.UpdateConfiguration(c, false)
	c, _ = services.ReadConfiguration()
	log.Printf("configuration: %+v\n", *c)
}
