package engine

import (
  "fmt"
  "github.com/gin-gonic/gin"
  "github.com/sirupsen/logrus"
  "leapp_daemon/interface/http/controller"
  "leapp_daemon/interface/http/middleware"
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
	initializeRoutes(engineWrapper.ginEngine)
}

func (engineWrapper *engineWrapper) Serve(port int) {
	err := engineWrapper.ginEngine.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		logrus.Fatalln("error:", err.Error())
	}
}

func initializeRoutes(ginEngine *gin.Engine) {
	v1 := ginEngine.Group("/http/v1")
  {
    v1.POST("/configuration/", controller.CreateConfigurationController)
    v1.GET("/configuration/", controller.ReadConfigurationController)

    v1.GET("/session/list", controller.ListSessionController)
    v1.POST("/session/mfa/token/confirm", controller.ConfirmMfaTokenController)

    v1.GET("/session/plain/:id", controller.GetPlainAwsSessionController)
    v1.POST("/session/plain", controller.CreatePlainAwsSessionController)
    v1.PUT("/session/plain/:id", controller.EditPlainAwsSessionController)
    v1.DELETE("/session/plain/:id", controller.DeletePlainAwsSessionController)
    v1.POST("/session/plain/:id/start", controller.StartPlainAwsSessionController)
    v1.POST("/session/plain/:id/stop", controller.StopPlainAwsSessionController)

    v1.GET("/session/federated/:id", controller.GetFederatedAwsSessionController)
    v1.POST("/session/federated", controller.CreateFederatedAwsSessionController)
    v1.PUT("/session/federated/:id", controller.EditFederatedAwsSessionController)
    v1.DELETE("/session/federated/:id", controller.DeleteFederatedAwsSessionController)
    v1.POST("/session/federated/:id/start", controller.StartFederatedAwsSessionController)
    v1.POST("/session/federated/:id/stop", controller.StopFederatedAwsSessionController)

    v1.GET("/session/trusted/:id", controller.GetTrustedAwsSessionController)
    v1.POST("/session/trusted", controller.CreateTrustedAwsSessionController)
    v1.PUT("/session/trusted/:id", controller.EditTrustedAwsSessionController)
    v1.DELETE("/session/trusted/:id", controller.DeleteTrustedAwsSessionController)

    v1.POST("/g_suite_auth/first_step", controller.GSuiteAuthFirstStepController)
    v1.POST("/g_suite_auth/second_step", controller.GSuiteAuthSecondStepController)
    v1.POST("/g_suite_auth/third_step", controller.GSuiteAuthThirdStepController)

    v1.GET("/region/aws/list", controller.GetAwsRegionListController)
    v1.PUT("/region/aws/", controller.EditAwsRegionController)

    v1.GET("/named_profile/aws/list", controller.ListAwsNamedProfileController)

    v1.GET("/ws", controller.WsController)
  }
}
