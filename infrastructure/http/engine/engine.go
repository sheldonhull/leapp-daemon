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
	engineEnvironment := controller.EngineEnvironment{Providers: providers}

	v1 := ginEngine.Group("/api/v1")
	{
		v1.GET("/session/list", engineEnvironment.ListSessionController)
		v1.POST("/session/mfa/token/confirm", engineEnvironment.ConfirmMfaTokenController)

		// AWS
		v1.GET("/session/plain/:id", engineEnvironment.GetPlainAwsSessionController)
		v1.POST("/session/plain", engineEnvironment.CreatePlainAwsSessionController)
		v1.PUT("/session/plain/:id", engineEnvironment.UpdatePlainAwsSessionController)
		v1.DELETE("/session/plain/:id", engineEnvironment.DeletePlainAwsSessionController)
		v1.POST("/session/plain/:id/start", engineEnvironment.StartPlainAwsSessionController)
		v1.POST("/session/plain/:id/stop", engineEnvironment.StopPlainAwsSessionController)

		v1.GET("/session/federated/:id", engineEnvironment.GetFederatedAwsSessionController)
		v1.POST("/session/federated", engineEnvironment.CreateFederatedAwsSessionController)
		v1.PUT("/session/federated/:id", engineEnvironment.EditFederatedAwsSessionController)
		v1.DELETE("/session/federated/:id", engineEnvironment.DeleteFederatedAwsSessionController)
		v1.POST("/session/federated/:id/start", engineEnvironment.StartFederatedAwsSessionController)
		v1.POST("/session/federated/:id/stop", engineEnvironment.StopFederatedAwsSessionController)

		v1.GET("/session/trusted/:id", engineEnvironment.GetTrustedAwsSessionController)
		v1.POST("/session/trusted", engineEnvironment.CreateTrustedAwsSessionController)
		v1.PUT("/session/trusted/:id", engineEnvironment.EditTrustedAwsSessionController)
		v1.DELETE("/session/trusted/:id", engineEnvironment.DeleteTrustedAwsSessionController)

		v1.GET("/region/aws/list", engineEnvironment.GetAwsRegionListController)
		v1.PUT("/region/aws/", engineEnvironment.EditAwsRegionController)

		v1.GET("/named_profile/aws/list", engineEnvironment.ListAwsNamedProfileController)

		// GCP
		v1.GET("/gcp/oauth/url", engineEnvironment.GetGcpOauthUrl)
		v1.POST("/gcp/session/plain", engineEnvironment.CreateGcpPlainSession)
		v1.GET("/gcp/session/plain/:id", engineEnvironment.GetGcpPlainSession)
		v1.PUT("/gcp/session/plain/:id", engineEnvironment.EditGcpPlainSession)
		v1.DELETE("/gcp/session/plain/:id", engineEnvironment.DeleteGcpPlainSession)
		v1.POST("/gcp/session/plain/:id/start", engineEnvironment.StartGcpPlainSession)
		v1.POST("/gcp/session/plain/:id/stop", engineEnvironment.StopGcpPlainSession)

		v1.GET("/ws", engineEnvironment.WsController)
	}
}
