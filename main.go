package main

import (
	"leapp_daemon/controller"
	"leapp_daemon/engine"
	"leapp_daemon/logging"
	"leapp_daemon/service"
	"leapp_daemon/service/session"
)

func main() {
	go controller.Hub.Run()

	defer logging.CloseLogFile()
	eng := engine.Engine()
	eng.Serve(8080)


	// Test Generate config file
	service.CreateConfiguration()

	configuration, err := service.ReadConfiguration()
	if err != nil {
		println(err.Error())
	}

	// Test Delete account
	if len(configuration.PlainAwsAccountSessions) > 0 {
		err = session.DeletePlainAwsSession(configuration.PlainAwsAccountSessions[0].Id)
		if err != nil {
			println(err.Error())
		} else {
			println("account deleted")
		}
	}

	// Test sessions crud creation
	err = session.CreatePlainAwsSession("Test Session", "12345678911", "eu-west-1", "ciccio", "arn:1234:etzy")
	if err == nil {
		println("account created")
	}

	// We can't create more than one account plain with same number and user
	err = session.CreatePlainAwsSession("Test Session 2", "12345678911", "eu-west-2", "ciccio", "arn:4567:etzy")
	if err != nil {
		println(err.Error())
	}

	configuration, err = service.ReadConfiguration()
	println(configuration.PlainAwsAccountSessions[0].Id)
	println(configuration.PlainAwsAccountSessions[0].Account.Name)

	// Edit
	err = session.EditPlainAwsSession(configuration.PlainAwsAccountSessions[0].Id,"Test Session 2b", "12345678911", "eu-west-2", "ciccio", "arn:4567:etzy")
	if err != nil {
		println(err.Error())
	}

	configuration, err = service.ReadConfiguration()
	println(configuration.PlainAwsAccountSessions[0].Id)
	println(configuration.PlainAwsAccountSessions[0].Account.Name)

	// Can't Edit if id is wrong
	err = session.EditPlainAwsSession("00000000000","Test Session 2b", "12345678911", "eu-west-2", "ciccio", "arn:4567:etzy")
	if err != nil {
		println(err.Error())
	}

	// Test list
	// a) add another session
	_ = session.CreatePlainAwsSession("Test Session Alpha", "0011001100", "us-east-2", "panzor", "arn:1111:etzy")
	// b) list
	println("------------------------------")
	println("List without filters")
	list, _ := session.ListPlainAwsSession("")
	for index, _ := range list { println(list[index].Account.Name) }
	// c) apply some filters
	println("------------------------------")
	println("Find second by name")
	list, _ = session.ListPlainAwsSession("Alpha")
	println(list[0].Account.Name)
	println("------------------------------")
	println("Find first by account number")
	list, _ = session.ListPlainAwsSession("911")
	println(list[0].Account.Name)
	println("------------------------------")
	println("Find second by region")
	list, _ = session.ListPlainAwsSession("east")
	println(list[0].Account.Name)
	println("------------------------------")
	println("Find first by user")
	list, _ = session.ListPlainAwsSession("cicc")
	println(list[0].Account.Name)
	println("------------------------------")
	println("Find second by Mfa device")
	list, _ = session.ListPlainAwsSession("rn:111")
	println(list[0].Account.Name)
	println("------------------------------")
	println("No found if query doesn't match")
	list, _ = session.ListPlainAwsSession("rn:9uiuhu")
	for index, _ := range list { println(list[index].Account.Name) }
}