package main

import (
	"leapp_daemon/controllers"
	"leapp_daemon/engine"
	"leapp_daemon/logging"
	"leapp_daemon/services"
	"leapp_daemon/services/sessions"
)

func main() {
	go controllers.Hub.Run()

	defer logging.CloseLogFile()
	eng := engine.Engine()
	eng.Serve(8080)


	// Test Generate config file
	services.CreateConfiguration()

	configuration, err := services.ReadConfiguration()
	if err != nil {
		println(err.Error())
	}

	// Test Delete account
	if len(configuration.PlainAwsAccountSessions) > 0 {
		err = sessions.DeletePlainAwsSession(configuration.PlainAwsAccountSessions[0].Id)
		if err != nil {
			println(err.Error())
		} else {
			println("account deleted")
		}
	}

	// Test sessions crud creation
	err = sessions.CreatePlainAwsSession("Test Session", "12345678911", "eu-west-1", "ciccio", "arn:1234:etzy")
	if err == nil {
		println("account created")
	}

	// We can't create more than one account plain with same number and user
	err = sessions.CreatePlainAwsSession("Test Session 2", "12345678911", "eu-west-2", "ciccio", "arn:4567:etzy")
	if err != nil {
		println(err.Error())
	}

	configuration, err = services.ReadConfiguration()
	println(configuration.PlainAwsAccountSessions[0].Id)
	println(configuration.PlainAwsAccountSessions[0].Account.Name)

	// Edit
	err = sessions.EditPlainAwsSession(configuration.PlainAwsAccountSessions[0].Id,"Test Session 2b", "12345678911", "eu-west-2", "ciccio", "arn:4567:etzy")
	if err != nil {
		println(err.Error())
	}

	configuration, err = services.ReadConfiguration()
	println(configuration.PlainAwsAccountSessions[0].Id)
	println(configuration.PlainAwsAccountSessions[0].Account.Name)

	// Can't Edit if id is wrong
	err = sessions.EditPlainAwsSession("00000000000","Test Session 2b", "12345678911", "eu-west-2", "ciccio", "arn:4567:etzy")
	if err != nil {
		println(err.Error())
	}

	// Test list
	// a) add another session
	_ = sessions.CreatePlainAwsSession("Test Session Alpha", "0011001100", "us-east-2", "panzor", "arn:1111:etzy")
	// b) list
	println("------------------------------")
	println("List without filters")
	list, _ := sessions.ListPlainAwsSession("")
	for index, _ := range list { println(list[index].Account.Name) }
	// c) apply some filters
	println("------------------------------")
	println("Find second by name")
	list, _ = sessions.ListPlainAwsSession("Alpha")
	println(list[0].Account.Name)
	println("------------------------------")
	println("Find first by account number")
	list, _ = sessions.ListPlainAwsSession("911")
	println(list[0].Account.Name)
	println("------------------------------")
	println("Find second by region")
	list, _ = sessions.ListPlainAwsSession("east")
	println(list[0].Account.Name)
	println("------------------------------")
	println("Find first by user")
	list, _ = sessions.ListPlainAwsSession("cicc")
	println(list[0].Account.Name)
	println("------------------------------")
	println("Find second by Mfa device")
	list, _ = sessions.ListPlainAwsSession("rn:111")
	println(list[0].Account.Name)
	println("------------------------------")
	println("No found if query doesn't match")
	list, _ = sessions.ListPlainAwsSession("rn:9uiuhu")
	for index, _ := range list { println(list[index].Account.Name) }
}