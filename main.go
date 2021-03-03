package main

import (
	"leapp_daemon/controllers"
	"leapp_daemon/engine"
	"leapp_daemon/logging"
	"leapp_daemon/services"
	"leapp_daemon/services/accounts"
)

func main() {
	go controllers.Hub.Run()

	defer logging.CloseLogFile()
	eng := engine.Engine()
	eng.Serve(8080)


	// Test accounts crud creation
	err := accounts.CreatePlainAwsSession("Test Session", "12345678911", "eu-west-1", "ciccio", "arn:1234:etzy")
	if err != nil {
		println(err)
	}
	// We can't create more than one account plain with same number and user
	err = accounts.CreatePlainAwsSession("Test Session 2", "12345678911", "eu-west-2", "ciccio", "arn:4567:etzy")
	if err != nil {
		println(err)
	}

	configuration, err := services.ReadConfiguration()
	if err != nil {
		println(err)
	}

	println(configuration)
}