package main

import (
	"leapp_daemon/controllers"
	"leapp_daemon/engine"
	"leapp_daemon/logging"

	"leapp_daemon/services"
)

func main() {
	// Test Secrets Save
	errSave := services.SaveSecret("megasegreto", "test-leapp")
	if errSave != nil {
		println("crash secret service save err:", errSave)
		return
	}
	// Test Secret Retrieve
	secret, errRead := services.RetrieveSecret("test-leapp")
	if errRead != nil {
		println("crash secret service read err:", errRead)
		return
	}
	println("secret:", secret)




	go controllers.Hub.Run()

	defer logging.CloseLogFile()
	eng := engine.Engine()
	eng.Serve(8080)
}