package engine

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"leapp_daemon/infrastructure/http/middleware"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/interface/http/controller"
	providers2 "leapp_daemon/providers"
)

type engineWrapper struct {
	providers *providers2.Providers
	ginEngine *gin.Engine
}

var engineWrapperInstance *engineWrapper = nil

func newEngineWrapper(providers *providers2.Providers) *engineWrapper {
	ginEngine := gin.New()

	engineWrapper := engineWrapper{
		ginEngine: ginEngine,
		providers: providers,
	}

	engineWrapper.initialize()

	return &engineWrapper
}

func Engine(providers *providers2.Providers) *engineWrapper {
	if engineWrapperInstance != nil {
		return engineWrapperInstance
	} else {
		engineWrapperInstance = newEngineWrapper(providers)
		return engineWrapperInstance
	}
}

func (engineWrapper *engineWrapper) initialize() {
	logging.InitializeLogger()
	engineWrapper.ginEngine.Use(middleware.ErrorHandler.Handle)
	initializeRoutes(engineWrapper.ginEngine, engineWrapper.providers)
}

func (engineWrapper *engineWrapper) Serve(port int) {
	err := engineWrapper.ginEngine.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		logrus.Fatalln("error:", err.Error())
	}
}

func initializeRoutes(ginEngine *gin.Engine, providers *providers2.Providers) {
	controller := controller.EngineController{Providers: providers}

	v1 := ginEngine.Group("/api/v1")
	{
		v1.GET("/session/list", controller.ListSession)
		v1.POST("/session/mfa/token/confirm", controller.ConfirmMfaToken)

		// AWS
		v1.GET("/session/plain/:id", controller.GetAwsPlainSession)
		v1.POST("/session/plain", controller.CreateAwsPlainSession)
		v1.PUT("/session/plain/:id", controller.UpdateAwsPlainSession)
		v1.DELETE("/session/plain/:id", controller.DeleteAwsPlainSession)
		v1.POST("/session/plain/:id/start", controller.StartAwsPlainSession)
		v1.POST("/session/plain/:id/stop", controller.StopAwsPlainSession)

		v1.GET("/session/federated/:id", controller.GetAwsFederatedSession)
		v1.POST("/session/federated", controller.CreateAwsFederatedSession)
		v1.PUT("/session/federated/:id", controller.EditAwsFederatedSession)
		v1.DELETE("/session/federated/:id", controller.DeleteAwsFederatedSession)
		v1.POST("/session/federated/:id/start", controller.StartAwsFederatedSession)
		v1.POST("/session/federated/:id/stop", controller.StopAwsFederatedSession)

		v1.GET("/session/trusted/:id", controller.GetAwsTrustedSession)
		v1.POST("/session/trusted", controller.CreateAwsTrustedSession)
		v1.PUT("/session/trusted/:id", controller.EditAwsTrustedSession)
		v1.DELETE("/session/trusted/:id", controller.DeleteAwsTrustedSession)

		v1.GET("/region/aws/list", controller.GetAwsRegionList)
		v1.PUT("/region/aws/", controller.EditAwsRegion)

		v1.GET("/named_profile/aws/list", controller.ListNamedProfiles)

		// GCP
		v1.GET("/gcp/oauth/url", controller.GetGcpOauthUrl)
		v1.POST("/gcp/session/plain", controller.CreateGcpPlainSession)
		v1.GET("/gcp/session/plain/:id", controller.GetGcpPlainSession)
		v1.PUT("/gcp/session/plain/:id", controller.EditGcpPlainSession)
		v1.DELETE("/gcp/session/plain/:id", controller.DeleteGcpPlainSession)
		v1.POST("/gcp/session/plain/:id/start", controller.StartGcpPlainSession)
		v1.POST("/gcp/session/plain/:id/stop", controller.StopGcpPlainSession)

		v1.GET("/ws", controller.GetWs)
	}
}
