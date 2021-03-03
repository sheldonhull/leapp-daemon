package main

import (
	"leapp_daemon/services"
	"leapp_daemon/services/accounts"
)

func main() {
	// go controllers.Hub.Run()

	// defer logging.CloseLogFile()
	// eng := engine.Engine()
	// eng.Serve(8080)


	// Test Generate config file
	services.CreateConfiguration()

	configuration, err := services.ReadConfiguration()
	if err != nil {
		println(err.Error())
	}

	// Test Delete account
	if len(configuration.PlainAwsAccountSessions) > 0 {
		err = accounts.DeletePlainAwsSession(configuration.PlainAwsAccountSessions[0].Id)
		if err != nil {
			println(err.Error())
		} else {
			println("account deleted")
		}
	}

	// Test accounts crud creation
	err = accounts.CreatePlainAwsSession("Test Session", "12345678911", "eu-west-1", "ciccio", "arn:1234:etzy")
	if err == nil {
		println("account created")
	}

	// We can't create more than one account plain with same number and user
	err = accounts.CreatePlainAwsSession("Test Session 2", "12345678911", "eu-west-2", "ciccio", "arn:4567:etzy")
	if err != nil {
		println(err.Error())
	}

	configuration, err = services.ReadConfiguration()
	println(configuration.PlainAwsAccountSessions[0].Id)
	println(configuration.PlainAwsAccountSessions[0].Account.Name)

	// Edit
	err = accounts.EditPlainAwsSession(configuration.PlainAwsAccountSessions[0].Id,"Test Session 2b", "12345678911", "eu-west-2", "ciccio", "arn:4567:etzy")
	if err != nil {
		println(err.Error())
	}

	configuration, err = services.ReadConfiguration()
	println(configuration.PlainAwsAccountSessions[0].Id)
	println(configuration.PlainAwsAccountSessions[0].Account.Name)

	// Can't Edit if id is wrong
	err = accounts.EditPlainAwsSession("00000000000","Test Session 2b", "12345678911", "eu-west-2", "ciccio", "arn:4567:etzy")
	if err != nil {
		println(err.Error())
	}
}