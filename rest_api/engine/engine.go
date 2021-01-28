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
	log.Println("engineWrapperInstance:", engineWrapperInstance)
	if engineWrapperInstance != nil {
		log.Println("engineWrapperInstance:", engineWrapperInstance)
		return engineWrapperInstance
	} else {
		engineWrapperInstance = newEngineWrapper()
		log.Println("engineWrapperInstance:", engineWrapperInstance)
		return engineWrapperInstance
	}
}

func (engineWrapper *engineWrapper) initialize() {
	engineWrapper.ginEngine.Use(gin.Logger())
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
