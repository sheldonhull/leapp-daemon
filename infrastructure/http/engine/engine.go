package engine

import (
	"fmt"
	"leapp_daemon/infrastructure/http/middleware"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/interface/http/controller"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	v1 := ginEngine.Group("/api/v1")
	{
		v1.GET("/session/list", controller.ListSessionController)
		v1.POST("/session/mfa/token/confirm", controller.ConfirmMfaTokenController)

		v1.GET("/session/plain/:id", controller.GetPlainAwsSessionController)
		v1.POST("/session/plain", controller.CreatePlainAwsSessionController)
		v1.PUT("/session/plain/:id", controller.UpdatePlainAwsSessionController)
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

		v1.GET("/region/aws/list", controller.GetAwsRegionListController)
		v1.PUT("/region/aws/", controller.EditAwsRegionController)

		v1.GET("/named_profile/aws/list", controller.ListAwsNamedProfileController)

		v1.GET("/ws", controller.WsController)

		v1.GET("/plain/alibaba/session/:id", controller.GetPlainAlibabaSessionController)
		v1.POST("/plain/alibaba/session/", controller.CreatePlainAlibabaSessionController)
		v1.PUT("/plain/alibaba/session/:id", controller.UpdatePlainAlibabaSessionController)
		v1.DELETE("/plain/alibaba/session/:id", controller.DeletePlainAlibabaSessionController)
		v1.POST("/plain/alibaba/session/:id/start", controller.StartPlainAlibabaSessionController)
		v1.POST("/plain/alibaba/session/:id/stop", controller.StopPlainAlibabaSessionController)
	}
}
