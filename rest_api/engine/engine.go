package engine

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"leapp_daemon/rest_api/controllers"
	"log"
)

type engineWrapper struct {
	ginEngine *gin.Engine
}

var engineWrapperInstance *engineWrapper = nil

func newEngineWrapper() *engineWrapper {
	ginEngine := gin.New()

	engineWrapper := engineWrapper{
		ginEngine: ginEngine,
	}

	engineWrapper.initialize()

	return &engineWrapper
}

func Engine() *engineWrapper {
	if engineWrapperInstance != nil {
		return engineWrapperInstance
	} else {
		engineWrapperInstance = newEngineWrapper()
		return engineWrapperInstance
	}
}

func (engineWrapper *engineWrapper) initialize() {
	// TODO: configure logging the proper way
	/*
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	    if err != nil {
	        log.Fatal(err)
	    }

	    log.SetOutput(file)
	 */
	//engineWrapper.ginEngine.Use(gin.Logger())
	engineWrapper.ginEngine.Use(gin.Recovery())
	initializeRoutes(engineWrapper.ginEngine)
}

func (engineWrapper *engineWrapper) Serve(port int) {
	err := engineWrapper.ginEngine.Run(fmt.Sprintf(":%d", port))

	if err != nil {
		log.Fatalln("Error:", fmt.Sprintf("%+v", err))
	}
}

func initializeRoutes(ginEngine *gin.Engine) {
	v1 := ginEngine.Group("/api/v1")
	{
		v1.GET("/home/:name", controllers.HomeController)
	}
}
