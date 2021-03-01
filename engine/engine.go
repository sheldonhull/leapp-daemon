package engine

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"leapp_daemon/controllers"
	"leapp_daemon/controllers/middleware"
	"leapp_daemon/logging"
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
	logging.InitializeLogger()
	engineWrapper.ginEngine.Use(middleware.ErrorHandler.Handle)
	engineWrapper.ginEngine.Use(gin.Recovery())
	initializeRoutes(engineWrapper.ginEngine)
}

func (engineWrapper *engineWrapper) Serve(port int) {
	err := engineWrapper.ginEngine.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		logrus.Fatalln("error:", err.Error())
	}
}

func initializeRoutes(ginEngine *gin.Engine) {
	v1 := ginEngine.Group("/api/v1")
	{
		v1.GET("/echo/:text", controllers.EchoController)
		v1.POST("/configuration/create", controllers.CreateConfigurationController)
		v1.GET("/configuration/read", controllers.ReadConfigurationController)
		v1.POST("/federated_aws_account/create", controllers.CreateFederatedAccountController)
		v1.POST("/g_suite_auth/first_step", controllers.GSuiteAuthFirstStepController)
		v1.POST("/g_suite_auth/second_step", controllers.GSuiteAuthSecondStepController)
		v1.POST("g_suite_auth/third_step", controllers.GSuiteAuthThirdStepController)
		v1.GET("/ws/:roomId", controllers.WsController)
	}
}
