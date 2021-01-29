package main

import (
	"leapp_daemon/services"
	"log"
)

func main() {
	// TODO: add JSON parsing and encoding
	// TODO: add error handling
	// TODO: add logging

	// TODO: add DTOs
	// TODO: add DTOs validation
	// TODO: add services layer

	// TODO: add domain layer, which implements the business logic, and should be independent from the REST API engine

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
	log.Printf("configuration: %+v\n", *c)
}
